// Provides functions to summarize content using Instructor and Ollama.

import { z } from 'zod';
import { SingleBar, Presets } from 'cli-progress';
import { Article, OllamaConfig, SummaryResponseSchema } from '../types/index';
import Instructor from '@instructor-ai/instructor';

// Emotional stimuli phrases inspired by EmotionPrompt research
// @see {@link https://python.useinstructor.com/prompting/zero_shot/emotion_prompting/|EmotionPrompt Documentation}
const EMOTIONAL_STIMULI = [
  'This summary is crucial for understanding cutting-edge tech trends.',
  'Your precise summary will guide critical decision-making.',
  'Deliver clarity to help readers quickly grasp key insights.',
  'Accuracy is key—provide a focused, detail-rich summary.',
];

const createArticleSummarizer = (
  client: ReturnType<typeof Instructor>,
  config: OllamaConfig,
  schema: z.AnyZodObject = SummaryResponseSchema
) => {
  const MAX_CONTENT_LENGTH = config.maxContentLength || 2000;
  const MAX_OUTPUT_TOKENS = config.maxSummaryLength || 150;

  const sanitizeInput = (text: string) =>
    text.replace(
      /[<>&]/g,
      (char) => ({ '<': '&lt;', '>': '&gt;', '&': '&amp;' }[char] || char)
    );

  const summarizeContent = async (
    title: string,
    content: string
  ): Promise<string> => {
    const truncatedContent = content.slice(0, MAX_CONTENT_LENGTH);
    const truncationNotice =
      content.length > MAX_CONTENT_LENGTH
        ? '\n[Truncated for length constraints]'
        : '';

    // Select a random emotional directive
    const emotionalDirective =
      EMOTIONAL_STIMULI[Math.floor(Math.random() * EMOTIONAL_STIMULI.length)];

    const prompt = `
SUMMARY REQUEST
---------------
INSTRUCTIONS:
- ${emotionalDirective}
- Provide a clear, concise summary of the Hacker News article.
- Emphasize key technical details, context, and innovative ideas.
- Do NOT preface with generic phrases such as "The article discusses..."
- Format the output as a single, well-structured paragraph.

ARTICLE:
--- TITLE ---
${sanitizeInput(title)}

--- TRUNCATED CONTENT ---
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
2. Avoid generic lead-ins.`,
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

      const parsedResponse = SummaryResponseSchema.parse(response);

      return parsedResponse.summary ?? 'No summary available';
    } catch (error) {
      console.error(
        `Error processing "${title.slice(0, 50)}...":`,
        error instanceof z.ZodError ? error.errors : error
      );
      return 'SUMMARY_ERROR';
    }
  };

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
      console.log(`\nProcessing article: ${article.title.slice(0, 60)}...`);
      const summary = await summarizeContent(article.title, article.content);

      article.summary = summary || 'Summary unavailable';
      progressBar.update(index + 1);
    }

    progressBar.stop();
    return articles;
  };

  return { summarizeArticles };
};

export { createArticleSummarizer };
