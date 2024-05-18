// Summarizes the content of articles fetched from top Hacker News stories using the Hugging Face API,
// and updates the database with the summaries.

const createArticleSummarizer = (axios, config, db) => {
  /**
   * Summarizes content using the Hugging Face API with retry logic.
   */
  const summarizeContentWithRetry = async (content, maxRetries = 3) => {
    let attempt = 0;
    let delay = 2000;

    while (attempt < maxRetries) {
      try {
        const response = await axios.post(
          config.huggingFace.apiUrl,
          { inputs: content, parameters: { truncation: 'only_first' } },
          {
            headers: {
              Authorization: `Bearer ${config.huggingFace.apiKey}`,
              'Content-Type': 'application/json',
            },
          }
        );

        if (response.status === 200) {
          return response.data[0]?.summary_text ?? '';
        }

        console.error(
          `Attempt ${attempt + 1}: Unexpected response status: ${
            response.status
          }`
        );
        attempt++;
      } catch (err) {
        if (err.response && [500, 503].includes(err.response.status)) {
          console.error(
            `Attempt ${attempt + 1}: Server error, retrying after ${
              delay / 1000
            } seconds...`
          );
          await new Promise((resolve) => setTimeout(resolve, delay));
          delay *= 2; // Double the delay for the next retry
          attempt++;
        } else {
          console.error('Non-retryable error:', err.message);
          break;
        }
      }
    }

    return ''; // Return an empty summary if all retries fail
  };

  /**
   * Fetches unsummarized articles from the database, summarizes them, and updates the database.
   */
  const summarizeFetchedArticles = async () => {
    const rows = await db.fetchUnsummarizedArticles();
    for (const { id, content } of rows) {
      if (!content) {
        console.warn(`Skipping Article ID ${id}: content is empty`);
        continue;
      }

      const summary = await summarizeContentWithRetry(content);
      await db.updateArticleSummary(id, summary);
      await new Promise((resolve) => setTimeout(resolve, 1000)); // Delay between summaries
    }
  };

  return {
    summarizeContentWithRetry,
    summarizeFetchedArticles,
  };
};

module.exports = { createArticleSummarizer };
