// Provides functions to summarize content using Instructor and Ollama.

import { z } from 'zod';
import { SingleBar, Presets } from 'cli-progress';
import { Article, OllamaConfig, SummaryResponseSchema } from '../types/index';
import Instructor from '@instructor-ai/instructor';

// Creates the article summarizer
const createArticleSummarizer = (
  client: ReturnType<typeof Instructor>,
  config: OllamaConfig,
  schema: z.AnyZodObject = SummaryResponseSchema
) => {
  const MAX_CONTENT_LENGTH = config.maxContentLength || 2000;
  const MAX_OUTPUT_TOKENS = config.maxSummaryLength || 150;
  const MIN_CONTENT_LENGTH = 300;

  // Replace special HTML characters
  const sanitizeInput = (text: string) =>
    text.replace(
      /[<>&]/g,
      (char) => ({ '<': '&lt;', '>': '&gt;', '&': '&amp;' }[char] || char)
    );

  // Check for captcha and redact IPv4 addresses
  const sanitizeSummary = (summary: string): string => {
    if (/captcha/i.test(summary)) {
      console.error('Captcha detected in summary, flagging as error.');
      return 'No summary available';
    }
    const ipRegex = /(?:(?:\d{1,3}\.){3}\d{1,3})/g;
    return summary.replace(ipRegex, 'REDACTED');
  };

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
      - NEVER hallucinate or fabricate content; only summarize what’s provided.
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
      ${sanitizeInput(truncatedContent)} ${truncationNotice}
    `;

    try {
      const response = await client.chat.completions.create({
        model: config.model,
        messages: [
          {
            role: 'system',
            content: prompt,
          },
        ],
        max_tokens: MAX_OUTPUT_TOKENS,
        temperature: 0.2,
        top_p: 0.9,
        response_model: { schema, name: 'SummarySchema' },
      });

      const parsed = SummaryResponseSchema.safeParse(response);
      const rawSummary = parsed.data?.summary ?? 'No summary available';
      return sanitizeSummary(rawSummary);
    } catch (error) {
      console.error(
        `Error processing "${title.slice(0, 50)}...":`,
        error instanceof Error ? error.message : 'Unknown error'
      );
      return 'No summary available';
    }
  };

  const summarizeArticles = async (articles: Article[]): Promise<Article[]> => {
    const progressBar = new SingleBar(
      {
        format:
          'Summarizing Articles |{bar}| {percentage}% | {value}/{total} Articles',
      },
      Presets.shades_classic
    );

    progressBar.start(articles.length, 0);

    for (const [index, article] of articles.entries()) {
      console.log(`\nProcessing article: ${article.title.slice(0, 60)}...`);
      article.summary = await summarizeContent(
        article.title,
        article.content ?? ''
      );
      progressBar.update(index + 1);
    }

    progressBar.stop();
    return articles;
  };

  return { summarizeArticles };
};

export { createArticleSummarizer };
