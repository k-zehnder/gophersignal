// Provides functions to summarize content using Instructor and Ollama.

import { z } from 'zod';
import { SingleBar, Presets } from 'cli-progress';
import { Article, OllamaConfig, SummaryResponseSchema } from '../types/index';
import Instructor from '@instructor-ai/instructor';

// Creates an article summarizer service using Instructor and Ollama
const createArticleSummarizer = (
  client: ReturnType<typeof Instructor>,
  config: OllamaConfig,
  schema: z.AnyZodObject = SummaryResponseSchema
) => {
  const MAX_CONTENT_LENGTH = config.maxContentLength || 2000;
  const MAX_OUTPUT_TOKENS = config.maxSummaryLength || 150;

  // Summarizes the content of a single article
  const summarizeContent = async (
    title: string,
    content: string
  ): Promise<string> => {
    const truncatedContent =
      content.length > MAX_CONTENT_LENGTH
        ? content.slice(0, MAX_CONTENT_LENGTH)
        : content;

    const prompt = `Provide a concise summary (max 150 words) of the article below.
    Title: ${title}

    Content:
    ${truncatedContent}

    Summary:
  `;

    try {
      const response = await client.chat.completions.create({
        model: config.model,
        messages: [{ role: 'user', content: prompt }],
        max_tokens: MAX_OUTPUT_TOKENS,
        temperature: 0.5,
        top_p: 0.9,
        response_model: { schema, name: 'SummarySchema' },
      });

      // Validate the response using Zod
      const parsedResponse = SummaryResponseSchema.parse(response);

      return parsedResponse.summary ?? 'No summary available';
    } catch (error) {
      if (error instanceof z.ZodError) {
        console.error('Zod validation error:', error.errors);
      } else {
        console.error(`Error summarizing content for "${title}":`, error);
      }
      return 'No summary available';
    }
  };

  // Summarizes a list of articles
  const summarizeArticles = async (
    articles: Required<Article>[]
  ): Promise<Article[]> => {
    const progressBar = new SingleBar(
      {
        format:
          'Summarizing Articles |{bar}| {percentage}% | {value}/{total} Articles',
      },
      Presets.shades_classic
    );

    progressBar.start(articles.length, 0);

    for (const [index, article] of articles.entries()) {
      console.log(`\nSummarizing article: ${article.title}`);
      const summary = await summarizeContent(article.title, article.content);

      if (summary) {
        article.summary = summary;
      } else {
        console.warn(
          `\nFailed to generate summary for article: ${article.title}`
        );
      }

      progressBar.update(index + 1);
    }

    progressBar.stop();

    return articles;
  };

  return {
    summarizeArticles,
  };
};

export { createArticleSummarizer };
