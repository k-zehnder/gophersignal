// Provides functions to summarize content using Instructor and Ollama with structured output and proper message roles.

import { z } from 'zod';
import { SingleBar, Presets } from 'cli-progress';
import Instructor from '@instructor-ai/instructor';
// Keep existing SummaryResponseSchema import
import { Article, OllamaConfig, SummaryResponseSchema } from '../types/index';

// Define the new StructuredSummarySchema
const StructuredSummarySchema = z.object({
  thinking: z.string().optional().describe("Internal analysis plan. NOT for user summary."),
  context: z.string().optional().describe("Article background/setting."),
  core_idea: z.string().optional().describe("Article central message/thesis."),
  insight_1: z.string().optional().describe("First key insight/argument."),
  insight_2: z.string().optional().describe("Second key insight/argument."),
  insight_3: z.string().optional().describe("Third insight. Omit if weak/absent."),
  insight_4: z.string().optional().describe("Fourth insight. Omit if weak/absent."),
  insight_5: z.string().optional().describe("Fifth insight. Omit if weak/absent."),
  author_conclusion: z.string().optional().describe("Author's conclusion/call to action."),
  warning: z.string().optional().describe("Note non-critical issues (e.g., ambiguity). Attempt summary."),
  error: z.string().optional().describe("CRITICAL: Use if summary IMPOSSIBLE (unreadable, CAPTCHA). Other fields: 'No summary available'."),
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
You are an assistant. Summarize Hacker News articles into structured JSON using the 'StructuredSummary' schema.

Instructions:
1.  **JSON ONLY**: Respond with a single, valid JSON object matching 'StructuredSummary'. No extra text.
2.  **Schema**:
    *   'thinking': Internal analysis. Not for user.
    *   Required: 'context', 'core_idea', 'insight_1', 'insight_2', 'insight_3', 'author_conclusion'.
    *   Optional: 'insight_4', 'insight_5'. Omit if weak/absent or article too short.
    *   NO HALLUCINATION. Use only provided text. Empty string or omit if info missing.
3.  **Issues**:
    *   'error': For CRITICAL issues (unreadable, CAPTCHA, too short <${MIN_CONTENT_LENGTH} chars, irrelevant). Sets other fields to "No summary available".
    *   'warning': For NON-CRITICAL issues (e.g., ambiguity). Still attempt summary.
4.  **Tone**: Neutral, factual, technical.
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

          const essentialFields = [context, core_idea, insight_1, insight_2, insight_3, author_conclusion];
          if (essentialFields.some(field => !field?.trim())) {
            console.warn(`[Structured Summary] Essential fields (context, core_idea, insight_1, insight_2, insight_3, author_conclusion) missing or empty for "${title}". Proceeding to fallback.`);
          } else {
            // Log warnings for missing truly optional fields (insights 4-5)
            if (!insight_4?.trim()) console.warn(`[Structured Summary] Optional field 'insight_4' missing or empty for "${title}".`);
            if (!insight_5?.trim()) console.warn(`[Structured Summary] Optional field 'insight_5' missing or empty for "${title}".`);

            const summaryLines = [
              context, core_idea, insight_1, insight_2, insight_3,
              insight_4, insight_5, // These can be undefined/empty
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
    
    const fallbackSystemPrompt = 'Assistant: Respond with JSON: {"summary": "your_summary_here"}.';
    const fallbackUserPrompt = `
      SUMMARY REQUEST:
      - Content missing, unreadable, or < ${MIN_CONTENT_LENGTH} chars? "summary": "No summary available".
      - NO HALLUCINATION. Summarize provided text only.
      - Summary: 5 lines, each with a unique role:
        1. Context (no "Context:" prefix).
        2. Core idea (no "Core idea:" prefix).
        3 & 4. Main insights (no labels).
        5. Author's conclusion (no label).
      - JSON ONLY: {"summary": "formatted_summary"}.
      - Tone: Neutral, factual, technical.

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
