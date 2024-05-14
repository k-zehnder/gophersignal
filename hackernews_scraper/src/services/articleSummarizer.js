// Summarizes content using the Hugging Face API with retry logic.

const axios = require('axios');

// Summarizes content using the Hugging Face API with retry logic.
// Retries the request up to a specified number of times if a server error occurs.
const summarizeContentWithRetry = async (
  apiUrl,
  apiKey,
  content,
  maxRetries = 3
) => {
  let attempt = 0;
  let delay = 2000; // Start with a 2-second delay

  while (attempt < maxRetries) {
    try {
      const response = await axios.post(
        apiUrl,
        { inputs: content, parameters: { truncation: 'only_first' } },
        {
          headers: {
            Authorization: `Bearer ${apiKey}`,
            'Content-Type': 'application/json',
          },
        }
      );

      if (response.status === 200) {
        // Return the summary if the request is successful
        return response.data[0]?.summary_text ?? '';
      }

      console.error(
        `Attempt ${attempt + 1}: Unexpected response status: ${response.status}`
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
        // Exit on non-retryable errors
        console.error('Non-retryable error:', err.message);
        break;
      }
    }
  }

  // Return an empty summary if all retries fail
  return '';
};

module.exports = { summarizeContentWithRetry };
