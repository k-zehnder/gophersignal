// Provides functions to summarize content using HuggingFace.

import fetch from 'node-fetch';
import { z } from 'zod';
import { Article, HuggingFaceConfig } from '../types/index';

// Factory to create the summarizer using the Hugging Face API
function createArticleSummarizer(config: HuggingFaceConfig) {
  const HUGGING_FACE_API_URL = 'https://api-inference.huggingface.co/models';

  // Define the Zod schema for response validation
  const SummarySchema = z.array(
    z.object({
      summary_text: z.string(),
    })
  );

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

      // Hugging Face API request
      const response = await fetch(`${HUGGING_FACE_API_URL}/${config.model}`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${config.apiKey}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          inputs: prompt,
          parameters: {
            max_new_tokens: MAX_OUTPUT_TOKENS,
            temperature: 0.5,
            top_p: 0.9,
          },
        }),
      });

      // Check if the response is ok
      if (!response.ok) {
        throw new Error(
          `Hugging Face API error: ${response.status} - ${response.statusText}`
        );
      }

      // Parse the response
      const data = await response.json();

      // Validate response using the Zod schema
      const parsedData = SummarySchema.safeParse(data);
      if (!parsedData.success) {
        throw new Error(
          'Invalid response structure received from Hugging Face API.'
        );
      }

      // Extract the summary text from the array
      const summary = parsedData.data[0]?.summary_text;
      if (!summary) {
        throw new Error('No summary text found in response.');
      }

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
