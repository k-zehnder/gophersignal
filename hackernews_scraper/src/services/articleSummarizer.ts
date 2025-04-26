// Provides functions to summarize content using Instructor and Ollama with robust JSON parsing and proper message roles.

import { z } from 'zod';
import { SingleBar, Presets } from 'cli-progress';
import Instructor from '@instructor-ai/instructor';
import { Article, OllamaConfig, SummaryResponseSchema } from '../types/index';

export const createArticleSummarizer = (
  client: ReturnType<typeof Instructor>,
  config: OllamaConfig,
  schema: z.AnyZodObject = SummaryResponseSchema
) => {
  const MAX_CONTENT_LENGTH = config.maxContentLength || 2000;
  const MAX_OUTPUT_TOKENS = config.maxSummaryLength || 150;
  const MIN_CONTENT_LENGTH = 300;

  // Escape HTML chars to avoid prompt injection.
  const sanitizeInput = (text: string) =>
    text.replace(
      /[<>&]/g,
      (c) => ({ '<': '&lt;', '>': '&gt;', '&': '&amp;' }[c] || c)
    );

  // Redact IPs and handle captcha flags.
  const sanitizeSummary = (summary: string): string =>
    /captcha/i.test(summary)
      ? 'No summary available'
      : summary.replace(/\b(?:\d{1,3}\.){3}\d{1,3}\b/g, 'REDACTED');

  // Capture the model name once for metadata
  const modelName = config.model;

  // Summarize a single article's content.
  const summarizeContent = async (
    title: string,
    content: string
  ): Promise<string> => {
    if (!content || content.trim().length < MIN_CONTENT_LENGTH) {
      return 'No summary available';
    }

    const truncatedContent = content.slice(0, MAX_CONTENT_LENGTH);
    const truncationNotice =
      content.length > MAX_CONTENT_LENGTH
        ? '\n[Truncated for length constraints]'
        : '';

    const prompt = `
      SUMMARY REQUEST
      ---------------
      INSTRUCTIONS:
      - If the article content is missing, unreadable, or under ${MIN_CONTENT_LENGTH} characters, return "No summary available".
      - NEVER hallucinate or fabricate content; only summarize what's provided.
      - Provide a clear, concise summary of the Hacker News article.
      - The summary must be exactly 5 lines long, with each line serving a unique role:
        * Line 1: Provide concise context.
        * Line 2: State the core idea.
        * Lines 3 & 4: Present the main insights supporting the core idea.
        * Line 5: Summarize the author's ultimate conclusion.
      - Return ONLY a JSON object with a single key "summary" containing the formatted summary.
      - Write in a neutral, factual tone suitable for a tech-savvy audience.

      ARTICLE:
      --- TITLE ---
      ${sanitizeInput(title)}

      --- CONTENT (truncated) ---
      ${sanitizeInput(truncatedContent)}${truncationNotice}
    `.trim();

    try {
      const response = await client.chat.completions.create({
        model: config.model,
        messages: [
          {
            role: 'system',
            content:
              'You are a helpful assistant. Respond with only the JSON object containing a "summary" field.',
          },
          { role: 'user', content: prompt },
        ],
        max_tokens: MAX_OUTPUT_TOKENS,
        temperature: 0.2,
        top_p: 0.9,
        response_model: { schema, name: 'SummaryResponse' },
      });

      // Extract parsed data only
      const container = (response as any).data ?? (response as any);
      const summaryStr: string | undefined =
        container && typeof container.summary === 'string'
          ? container.summary.trim()
          : undefined;

      // Finalize or fallback immediately
      if (!summaryStr) {
        return 'No summary available';
      }

      const finalSummary = sanitizeSummary(summaryStr);
      console.log('Result:', finalSummary);
      return finalSummary;
    } catch (err) {
      console.error(
        `Error processing "${title.slice(0, 50)}...":`,
        err instanceof Error ? err.message : err
      );
      return 'No summary available';
    }
  };

  // Summarize an array of articles with a progress bar.
  const summarizeArticles = async (articles: Article[]): Promise<Article[]> => {
    const bar = new SingleBar(
      {
        format: 'Summarizing Articles |{bar}| {percentage}% | {value}/{total}',
      },
      Presets.shades_classic
    );
    bar.start(articles.length, 0);

    for (let i = 0; i < articles.length; i++) {
      console.log(`\nProcessing: ${articles[i].title.slice(0, 60)}...`);
      articles[i].summary = await summarizeContent(
        articles[i].title,
        articles[i].content || ''
      );
      articles[i].modelName = modelName;
      bar.update(i + 1);
    }

    bar.stop();
    return articles;
  };

  return { summarizeArticles };
};
