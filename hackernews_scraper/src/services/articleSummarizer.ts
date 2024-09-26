// Provides functions to summarize content using Instructor and Ollama.

import Instructor from '@instructor-ai/instructor';
import OpenAI from 'openai';
import { z } from 'zod';
import { Article, OllamaConfig } from '../types/index';

function createArticleSummarizer(config: OllamaConfig) {
  // Initialize the OpenAI client pointing to Ollama's API
  const openaiClient = new OpenAI({
    apiKey: config.apiKey || 'ollama', // Placeholder API key
    baseURL: config.baseUrl,
  });

  // Initialize the Instructor client with the OpenAI client
  const client = Instructor({
    client: openaiClient,
    mode: 'JSON',
  });

  // Define the Zod schema for response validation
  const SummarySchema = z.object({
    summary: z.string(),
  });

  const summarizeContent = async (
    title: string,
    content: string
  ): Promise<string> => {
    try {
      const MAX_CONTENT_LENGTH = 2000;
      const MAX_OUTPUT_TOKENS = 150;

      // Truncate content if it exceeds the maximum length
      let truncatedContent = content;
      if (content.length > MAX_CONTENT_LENGTH) {
        truncatedContent = content.slice(0, MAX_CONTENT_LENGTH);
        console.warn(
          `Content for "${title}" was truncated from ${content.length} to ${MAX_CONTENT_LENGTH} characters.`
        );
      }

      const prompt = `Provide a concise summary (max 150 words) of the article below.

Title: ${title}

Content:
${truncatedContent}

Summary:
`;

      const response = await client.chat.completions.create({
        model: config.model,
        messages: [{ role: 'user', content: prompt }],
        max_tokens: MAX_OUTPUT_TOKENS,
        temperature: 0.5,
        top_p: 0.9,
        response_model: { schema: SummarySchema, name: 'SummarySchema' },
      });

      console.log('Response received:', response);

      const { summary } = response;
      return summary;
    } catch (error) {
      console.error('Error summarizing content:', error);
      return '';
    }
  };

  const summarizeArticles = async (
    articles: Required<Article>[]
  ): Promise<Article[]> => {
    for (const article of articles) {
      console.log(`Summarizing article: ${article.title}`);
      const summary = await summarizeContent(article.title, article.content);

      if (summary) {
        article.summary = summary;
      } else {
        console.warn(
          `Failed to generate summary for article: ${article.title}`
        );
      }
    }

    return articles;
  };

  return {
    summarizeArticles,
  };
}

export { createArticleSummarizer };
