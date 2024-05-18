// Summarizes the content of articles fetched from top Hacker News stories using the Hugging Face API.

const createArticleSummarizer = (axios, config) => {
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
   * Summarizes a list of articles.
   */
  const summarizeArticles = async (articles) => {
    for (const article of articles) {
      if (!article.content) {
        console.warn(`Skipping article with missing content: ${article.title}`);
        continue;
      }

      console.log(`Summarizing article: ${article.title}`);
      article.summary = await summarizeContentWithRetry(article.content);
      await new Promise((resolve) => setTimeout(resolve, 1000)); // Delay between summaries
    }

    return articles;
  };

  return {
    summarizeContentWithRetry,
    summarizeArticles,
  };
};

module.exports = { createArticleSummarizer };
