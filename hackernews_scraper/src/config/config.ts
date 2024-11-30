// Loads environment variables and provides configuration settings.

import dotenv from 'dotenv';
dotenv.config();

import { Config } from '../types/index';

const config: Config = {
  mysql: {
    host: process.env.MYSQL_HOST || 'localhost',
    port: parseInt(process.env.MYSQL_PORT || '3306'),
    user: process.env.MYSQL_USER || 'user',
    password: process.env.MYSQL_PASSWORD || '',
    database: process.env.MYSQL_DATABASE || 'database_name',
  },
  ollama: {
    baseUrl:
      process.env.OLLAMA_BASE_URL || 'http://localhost:11434/api/generate',
    model: process.env.OLLAMA_MODEL || 'llama3:instruct',
    apiKey: process.env.OLLAMA_API_KEY || 'ollama',
    maxContentLength: parseInt(process.env.MAX_CONTENT_LENGTH || '2000'),
    maxSummaryLength: parseInt(process.env.MAX_SUMMARY_LENGTH || '500'),
  },
};

export default config;
