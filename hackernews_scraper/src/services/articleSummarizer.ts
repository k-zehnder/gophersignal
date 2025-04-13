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
    const truncatedContent = content.slice(0, MAX_CONTENT_LENGTH);
    const truncationNotice =
      content.length > MAX_CONTENT_LENGTH
        ? '\n[Truncated for length constraints]'
        : '';

    const prompt = `
      SUMMARY REQUEST
      ---------------
      INSTRUCTIONS:
      - Provide a clear, concise summary of the Hacker News article.
      - The summary should be **2 to 3 sentences** long (approximately 50–70 words) and capture the core idea of the article.
      - Write in a clear, factual style suitable for a tech-savvy audience; assume the reader wants a quick, informative gist.
      - Highlight the main point and any important outcome or insight, while omitting trivial details or general background.
      - The tone should be neutral and informative.
      - Ensure the summary can stand on its own and remains within the optimal length for easy reading on both desktop and mobile.

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
            content: `You are a precise summarization AI specialized in Hacker News content. Follow these rules strictly:
            1. Provide factual, technical summaries in a single paragraph.
            2. Do NOT preface with generic lead-in phrases.
            3. Return ONLY a JSON object with a "summary" key containing the summary. Format must be: { "summary": "..." }.`,
          },
          {
            role: 'user',
            content: prompt,
          },
        ],
        max_tokens: MAX_OUTPUT_TOKENS,
        temperature: 0.5,
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
