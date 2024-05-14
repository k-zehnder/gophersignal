// Loads environment variables and provides configuration settings for the application.

require('dotenv').config();

// Configuration settings for the application, including debug mode and Hugging Face API details.
const config = {
  debugMode: process.env.DEBUG_MODE,
  huggingFace: {
    apiUrl:
      'https://api-inference.huggingface.co/models/facebook/bart-large-cnn',
    apiKey: process.env.HUGGING_FACE_API_KEY,
  },
};

module.exports = config;
