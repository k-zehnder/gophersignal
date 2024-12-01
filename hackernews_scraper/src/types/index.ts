import { z } from 'zod';

export interface Article {
  title: string;
  link: string;
  content?: string;
  summary?: string;
}

export interface MySQLConfig {
  host: string;
  port: number;
  user: string;
  password: string;
  database: string;
}

export interface OllamaConfig {
  baseUrl: string;
  model: string;
  apiKey?: string;
  maxContentLength: number;
  maxSummaryLength: number;
}

export interface Config {
  mysql: MySQLConfig;
  ollama: OllamaConfig;
}

export const SummarySchema = z.object({
  summary: z.string().optional(),
  response: z.string().optional(),
});
