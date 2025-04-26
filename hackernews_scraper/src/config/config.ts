// Loads environment variables and provides configuration settings.

import dotenv from 'dotenv';
import { z } from 'zod';
import { Config } from '../types/index';

dotenv.config();

const envSchema = z.object({
  MYSQL_HOST: z.string().default('localhost'),
  MYSQL_PORT: z.string().default('3306'),
  MYSQL_USER: z.string().default('user'),
  MYSQL_PASSWORD: z.string().default(''),
  MYSQL_DATABASE: z.string().default('database_name'),

  OLLAMA_BASE_URL: z.string().default('http://localhost:11434/api/generate'),
  OLLAMA_MODEL: z.string().default('llama3:instruct'),
  OLLAMA_API_KEY: z.string().optional(),
  MAX_CONTENT_LENGTH: z.string().default('2000'),
  MAX_SUMMARY_LENGTH: z.string().default('500'),

  GH_TOKEN: z.string().optional(),
  GITHUB_OWNER: z.string().default('k-zehnder'),
  GITHUB_REPO: z.string().default('gophersignal'),
  GITHUB_BRANCH: z.string().default('main'),
});

const env = envSchema.parse(process.env);

const config: Config = {
  mysql: {
    host: env.MYSQL_HOST,
    port: parseInt(env.MYSQL_PORT, 10),
    user: env.MYSQL_USER,
    password: env.MYSQL_PASSWORD,
    database: env.MYSQL_DATABASE,
  },
  ollama: {
    baseUrl: env.OLLAMA_BASE_URL,
    model: env.OLLAMA_MODEL,
    apiKey: env.OLLAMA_API_KEY,
    maxContentLength: parseInt(env.MAX_CONTENT_LENGTH, 10),
    maxSummaryLength: parseInt(env.MAX_SUMMARY_LENGTH, 10),
  },
  github: {
    token: env.GH_TOKEN,
    owner: env.GITHUB_OWNER,
    repo: env.GITHUB_REPO,
    branch: env.GITHUB_BRANCH,
  },
};

export default config;
