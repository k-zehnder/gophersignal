// Loads environment variables and provides configuration settings.

import dotenv from 'dotenv';
import { ConfigSchema, type Config } from '../types';

dotenv.config();

export const config: Config = ConfigSchema.parse({
  mysql: {
    host: process.env.MYSQL_HOST,
    port: process.env.MYSQL_PORT,
    user: process.env.MYSQL_USER,
    password: process.env.MYSQL_PASSWORD,
    database: process.env.MYSQL_DATABASE,
  },
  ollama: {
    baseUrl: process.env.OLLAMA_BASE_URL,
    model: process.env.OLLAMA_MODEL,
    apiKey: process.env.OLLAMA_API_KEY,
    maxContentLength: process.env.MAX_CONTENT_LENGTH,
    maxSummaryLength: process.env.MAX_SUMMARY_LENGTH,
  },
});

export default config;
