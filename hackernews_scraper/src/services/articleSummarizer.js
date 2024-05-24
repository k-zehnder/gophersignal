// Provides functions to summarize content using the Hugging Face API.
// Includes content truncation and retry logic with error handling.

const createArticleSummarizer = (axios, config) => {
  /**
   * Truncates content to a specified maximum length.
   */
  const truncateContent = (content, maxLength = 2000) => {
    if (content.length > maxLength) {
      return content.substring(0, maxLength);
    }
    return content;
  };

  /**
   * Handles errors and retries for the API request.
   */
  const handleError = async (err, attempt, delay) => {
    if (err.response) {
      if ([500, 503].includes(err.response.status)) {
        console.error(
          `Attempt ${attempt + 1}: Server error, retrying after ${
            delay / 1000
          } seconds.`
        );
      } else if (
        err.response.data?.error.includes(
          'Model facebook/bart-large-cnn is currently loading'
        )
      ) {
        const estimatedTime = err.response.data?.estimated_time || 60;
        console.error(
          `Attempt ${
            attempt + 1
          }: Model loading, retrying after ${estimatedTime} seconds.`
        );
        await new Promise((resolve) =>
          setTimeout(resolve, estimatedTime * 1000)
        );
        return true; // Indicate to retry without incrementing attempt
      } else {
        console.error('Non-retryable error:', err.message);
        return false; // Indicate not to retry
      }
    } else {
      console.error('Non-retryable error:', err.message);
      return false; // Indicate not to retry
    }

    await new Promise((resolve) => setTimeout(resolve, delay));
    return true; // Indicate to retry
  };

  /**
   * Summarizes content using the Hugging Face API with retry logic.
   */
  const summarizeContentWithRetry = async (content, maxRetries = 3) => {
    let attempt = 0;
    let delay = 2000;

    // Truncate content before sending to API
    const truncatedContent = truncateContent(content, 2048);

    while (attempt < maxRetries) {
      try {
        const requestData = JSON.stringify({
          inputs: truncatedContent,
          parameters: { max_length: 150, truncation: true },
        });

        const response = await axios.post(
          config.huggingFace.apiUrl,
          requestData,
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
          }.`
        );
        attempt++;
      } catch (err) {
        const shouldRetry = await handleError(err, attempt, delay);
        if (!shouldRetry) {
          break;
        }
        delay *= 2; // Double the delay for the next retry
        attempt++;
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
    }

    return articles;
  };

  return {
    summarizeContentWithRetry,
    summarizeArticles,
  };
};

module.exports = { createArticleSummarizer };
