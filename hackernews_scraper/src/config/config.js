// Loads environment variables and provides configuration settings for the application.

require('dotenv').config();

// Configuration settings for the application.
const config = {
  huggingFace: {
    apiUrl:
      'https://api-inference.huggingface.co/models/facebook/bart-large-cnn',
    apiKey: process.env.HUGGING_FACE_API_KEY,
  },
  mysql: {
    host: process.env.MYSQL_HOST,
    port: process.env.MYSQL_PORT,
    user: process.env.MYSQL_USER,
    password: process.env.MYSQL_PASSWORD,
    database: process.env.MYSQL_DATABASE,
  },
};

module.exports = config;
