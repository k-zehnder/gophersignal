// Provides functions to summarize content using Instructor and Ollama.

import { z } from 'zod';
import Instructor from '@instructor-ai/instructor';
import { SingleBar, Presets } from 'cli-progress';
import { type Article, OllamaConfig } from '../types';

export const SummaryResponseSchema = z.object({
  summary: z.string().optional(),
  _meta: z.any().optional(),
});

// Creates the article summarizer
const createArticleSummarizer = (
  client: ReturnType<typeof Instructor>,
  config: OllamaConfig,
  schema: z.AnyZodObject = SummaryResponseSchema
) => {
  const MAX_CONTENT = config.maxContentLength ?? 2000;
  const MAX_TOKENS = config.maxSummaryLength ?? 150;
  const MIN_LENGTH = 300;

  // Replace special HTML characters
  const sanitizeInput = (text: string) =>
    text.replace(
      /[<>&]/g,
      (c) => ({ '<': '&lt;', '>': '&gt;', '&': '&amp;' }[c] || c)
    );

  // Redact captchas or IPs
  const sanitizeSummary = (s: string) =>
    /captcha/i.test(s)
      ? 'No summary available'
      : s.replace(/(?:(?:\d{1,3}\.){3}\d{1,3})/g, 'REDACTED');

  // Generate or fallback
  const summarizeContent = async (
    title: string,
    content: string
  ): Promise<string> => {
    if (!content || content.trim().length < MIN_LENGTH) {
      return 'No summary available';
    }

    const snippet = content.slice(0, MAX_CONTENT);
    const notice =
      content.length > MAX_CONTENT
        ? '\n[Truncated for length constraints]'
        : '';
    const prompt = `
      SUMMARY REQUEST
      ---------------
      INSTRUCTIONS:
      - If the article content is missing, unreadable, or under ${MIN_LENGTH} characters, return "No summary available".
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
      ${sanitizeInput(snippet)} ${notice}
    `;

    try {
      const response = await client.chat.completions.create({
        model: config.model,
        messages: [{ role: 'system', content: prompt }],
        max_tokens: MAX_TOKENS,
        temperature: 0.2,
        top_p: 0.9,
        response_model: { schema, name: 'SummarySchema' },
      });

      // Zod parse enforces schema and provides default summary
      const { summary } = schema.parse(response);
      return sanitizeSummary(summary);
    } catch {
      return 'No summary available';
    }
  };

  // Process list with progress
  const summarizeArticles = async (articles: Article[]): Promise<Article[]> => {
    const bar = new SingleBar(
      { format: '|{bar}| {value}/{total}' },
      Presets.shades_classic
    );
    bar.start(articles.length, 0);

    for (const [idx, article] of articles.entries()) {
      console.log(`Processing: ${article.title}`);
      article.summary = await summarizeContent(
        article.title,
        article.content ?? ''
      );
      bar.update(idx + 1);
    }

    bar.stop();
    return articles;
  };

  return { summarizeArticles };
};

export { createArticleSummarizer };
