// Provides functions to summarize content using Instructor and Ollama with structured output and proper message roles.

import { z } from 'zod';
import { SingleBar, Presets } from 'cli-progress';
import Instructor from '@instructor-ai/instructor';
// Keep existing SummaryResponseSchema import
import { Article, OllamaConfig, SummaryResponseSchema } from '../types/index';

// Define the new StructuredSummarySchema
const StructuredSummarySchema = z.object({
  thinking: z.string().optional().describe(
    "Briefly outline your analysis and summary plan here. This is for internal reasoning and NOT part of the user-facing summary."
  ),
  context: z.string().optional().describe("Article's background/setting."),
  core_idea: z.string().optional().describe("Article's central message/thesis."),
  insight_1: z.string().optional().describe("First key insight or argument."),
  insight_2: z.string().optional().describe("Second key insight or argument."),
  insight_3: z.string().optional().describe("Third key insight. Omit if not distinct or strong."),
  insight_4: z.string().optional().describe("Fourth key insight. Omit if not distinct or strong."),
  insight_5: z.string().optional().describe("Fifth key insight. Omit if not distinct or strong."),
  author_conclusion: z.string().optional().describe("Author's main conclusion or call to action."),
  warning: z.string().optional().describe(
    "Note non-critical issues (e.g., minor ambiguity, slight off-topic). Still attempt summary."
  ),
  error: z.string().optional().describe(
    "CRITICAL: Use ONLY if summary is impossible (e.g., unreadable, CAPTCHA, irrelevant). If used, other fields may be 'No summary available'."
  ),
});

export const createArticleSummarizer = (
  client: ReturnType<typeof Instructor>,
  config: OllamaConfig,
  // Default schema for the `response_model` parameter is not directly used here anymore,
  // as summarizeContent will manage its schemas explicitly.
  // We can leave the original default or remove `schema` if not used elsewhere.
  // For this change, we ensure `summarizeContent` uses the correct schemas internally.
  // The diff showed `schema: z.AnyZodObject = StructuredSummarySchema`, let's update this for consistency,
  // though summarizeContent will override it.
  _schema_param_not_directly_used_in_summarize_content: z.AnyZodObject = StructuredSummarySchema // Renamed to clarify
) => {
  const MAX_CONTENT_LENGTH = config.maxContentLength || 2000;
  const MAX_OUTPUT_TOKENS = config.maxSummaryLength || 150; // Consider if structured needs more
  const MIN_CONTENT_LENGTH = 300;
  const NUM_CTX = config.numCtx; // From OllamaConfig

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

  // Strip any “Label:” prefixes (e.g. “Context:”, “Core idea:”) from each line.
  const stripLabels = (text: string): string =>
    text
      .split('\n')
      .map((line) => line.replace(/^\s*\w+:\s*/, ''))
      .join('\n');

  // Collapse multiple blank lines into a single newline.
  const collapseBlankLines = (text: string): string =>
    text.replace(/\n\s*\n+/g, '\n');

  // Helper for post-processing text
  const postProcessText = (text: string): string => {
    const sanitized = sanitizeSummary(text.trim());
    const labeledStripped = stripLabels(sanitized);
    return collapseBlankLines(labeledStripped);
  };

  // Capture the model name for metadata.
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

    // --- Attempt 1: Structured Summary ---
    const structuredSystemPrompt = `
You are an assistant specializing in summarizing Hacker News articles into a structured JSON format.
Your task is to populate a JSON object according to the 'StructuredSummary' schema.
Refer to the schema field descriptions for specific instructions on each field.

Key Instructions:
1.  **JSON Output ONLY**: Your entire response MUST be a single, valid JSON object matching the 'StructuredSummary' schema. No extra text, markdown, or explanations before or after the JSON.
2.  **Schema Adherence**:
    *   Use the 'thinking' field for your internal pre-summary analysis (not for the user).
    *   Prioritize 'context', 'core_idea', 'insight_1', 'insight_2', 'author_conclusion'.
    *   'insight_3' to 'insight_5' are optional; omit if not clearly present or if article is too short.
    *   NEVER hallucinate. Summarize based ONLY on the provided text. If info for a field is missing, use an empty string or omit (all fields are optional strings).
3.  **Issue Handling**:
    *   **Critical Issues**: If content is unusable (e.g., unreadable, CAPTCHA, too short <${MIN_CONTENT_LENGTH} chars, irrelevant to title), populate the 'error' field and set summary fields (context, core_idea, etc.) to "No summary available".
    *   **Non-Critical Issues**: Note minor issues (e.g., ambiguity) in the 'warning' field. Still attempt to provide the main summary.
4.  **Tone**: Neutral, factual, for a technical audience.
    `.trim();

    const structuredUserPrompt = `
<title>${sanitizeInput(title)}</title>
<content>${sanitizeInput(truncatedContent)}${truncationNotice}</content>
    `.trim();

    try {
      const structuredResponseData = await client.chat.completions.create({
        model: config.model,
        messages: [
          { role: 'system', content: structuredSystemPrompt },
          { role: 'user', content: structuredUserPrompt },
        ],
        max_tokens: MAX_OUTPUT_TOKENS + 50, // Allow a bit more for structured output
        temperature: 0.2,
        top_p: 0.9,
        response_model: { schema: StructuredSummarySchema, name: 'StructuredSummary' },
        options: { num_ctx: NUM_CTX },
      });

      const parseResult = StructuredSummarySchema.safeParse(structuredResponseData);

      if (parseResult.success) {
        const data = parseResult.data;
        const { context, core_idea, insight_1, insight_2, insight_3, insight_4, insight_5, author_conclusion, warning, error } = data;

        // Check for LLM-reported critical errors first
        if (error && error.trim().length >= 10) {
          console.warn(`[Structured Summary] LLM reported critical error for "${title}": ${error}. Proceeding to fallback.`);
        } else {
          // Log LLM-reported warnings
          if (warning && warning.trim().length > 0) {
            console.warn(`[Structured Summary LLM Warning] For article "${title}": ${warning}`);
          }

          const essentialFields = [context, core_idea, insight_1, insight_2, author_conclusion];
          if (essentialFields.some(field => !field?.trim())) {
            console.warn(`[Structured Summary] Essential fields (context, core_idea, insight_1, insight_2, author_conclusion) missing or empty for "${title}". Proceeding to fallback.`);
          } else {
            // Log warnings for missing truly optional fields (insights 3-5)
            if (!insight_3?.trim()) console.warn(`[Structured Summary] Optional field 'insight_3' missing or empty for "${title}".`);
            if (!insight_4?.trim()) console.warn(`[Structured Summary] Optional field 'insight_4' missing or empty for "${title}".`);
            if (!insight_5?.trim()) console.warn(`[Structured Summary] Optional field 'insight_5' missing or empty for "${title}".`);

            const summaryLines = [
              context, core_idea, insight_1, insight_2,
              insight_3, insight_4, insight_5, // These can be undefined/empty
              author_conclusion
            ]
            .map(line => line?.trim() || '')
            .filter(line => line.length > 0);

            if (summaryLines.length < 2) {
              console.warn(`[Structured Summary] Resulting summary for "${title}" has less than 2 lines (${summaryLines.length}). Proceeding to fallback.`);
            } else {
              return postProcessText(summaryLines.join('\n'));
            }
          }
        }
      } else {
        console.warn(`[Structured Summary] Parsing failed for "${title}": ${parseResult.error.message}. Proceeding to fallback.`);
      }
    } catch (error: any) {
      console.warn(`[Structured Summary] LLM call or unexpected error for "${title}": ${error.message}. Proceeding to fallback.`);
    }

    // --- Fallback: Simple "summary" field (original method) ---
    console.log(`[Fallback Summary] Attempting simple summary for "${title}".`);
    
    const fallbackSystemPrompt = 'You are a helpful assistant. Respond with only the JSON object containing a "summary" field.';
    const fallbackUserPrompt = `
      SUMMARY REQUEST
      ---------------
      INSTRUCTIONS:
      - If the article content is missing, unreadable, or under ${MIN_CONTENT_LENGTH} characters, return "No summary available" inside the "summary" field.
      - NEVER hallucinate or fabricate content; only summarize what's provided.
      - Provide a clear, concise summary of the Hacker News article.
      - The summary must be exactly 5 lines long, with each line serving a unique role:
        * Line 1: Provide concise context (no “Context:” prefix).
        * Line 2: State the core idea (no “Core idea:” prefix).
        * Lines 3 & 4: Present the main insights supporting the core idea (no literal labels).
        * Line 5: Summarize the author's ultimate conclusion (no label).
      - Return ONLY a JSON object with a single key "summary" containing the formatted summary.
      - Write in a neutral, factual tone suitable for a tech-savvy audience.

      ARTICLE:
      --- TITLE ---
      ${sanitizeInput(title)}

      --- CONTENT (truncated) ---
      ${sanitizeInput(truncatedContent)}${truncationNotice}
    `.trim();

    try {
      const fallbackResponse = await client.chat.completions.create({
        model: config.model,
        messages: [
          {role: 'system', content: fallbackSystemPrompt },
          {role: 'user', content: fallbackUserPrompt },
        ],
        max_tokens: MAX_OUTPUT_TOKENS,
        temperature: 0.2,
        top_p: 0.9,
        response_model: { schema: SummaryResponseSchema, name: 'SummaryResponseFallback' },
        options: { num_ctx: NUM_CTX },
      });
      
      // Assuming response_model ensures fallbackResponse is { summary: string } or throws
      const rawSummary = (fallbackResponse as { summary: string }).summary;

      if (!rawSummary?.trim() || rawSummary.trim().toLowerCase() === 'no summary available') {
        console.warn(`[Fallback Summary] LLM returned empty or "No summary available" for "${title}".`);
        return 'No summary available';
      }
      return postProcessText(rawSummary);
    } catch (fallbackError: any) {
      const errorMessage = `CRITICAL: Fallback summarization for article "${title}" failed. Original error: ${fallbackError.message}`;
      console.error(errorMessage, fallbackError);
      // As per requirement "if the fallback fails, we crash with a helpful message."
      throw new Error(errorMessage);
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
      try {
        articles[i].summary = await summarizeContent(
          articles[i].title,
          articles[i].content || ''
        );
      } catch (error) {
        // If summarizeContent throws (e.g. critical fallback failure), this catch block will handle it.
        // We will log the error and set a default "Error in summarization" message for this article.
        // The entire process for other articles will continue.
        console.error(`Error processing article "${articles[i].title}" in summarizeArticles loop: `, error);
        articles[i].summary = 'Error during summarization process.';
        // Optionally, re-throw if a single failure should stop everything:
        // throw error; 
      }
      articles[i].modelName = modelName;
      bar.update(i + 1);
    }

    bar.stop();
    return articles;
  };

  return { summarizeArticles };
};
