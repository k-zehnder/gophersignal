// Loads environment variables and provides configuration settings for the application.

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
  huggingFace: {
    apiKey: process.env.HUGGING_FACE_API_KEY || '',
    model: process.env.HUGGING_FACE_MODEL || 'facebook/bart-large-cnn',
  },
};

export default config;
