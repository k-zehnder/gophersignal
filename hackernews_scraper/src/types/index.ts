import { z } from 'zod';

export interface Article {
  title: string;
  link: string;
  content?: string;
  summary?: string;
  source?: string;
  upvotes?: number;
  comment_count?: number;
  comment_link?: string;
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

export interface Dependencies {
  db: {
    saveArticles: (articles: Article[]) => Promise<void>;
    updateArticleSummary?: (id: number, summary: string) => Promise<void>;
    closeDatabaseConnection: () => Promise<void>;
    connection?: any;
  };
  browser: {
    close: () => Promise<void>;
  };
  articleProcessor: {
    processTopStories: () => Promise<Article[]>;
  };
  articleSummarizer: {
    summarizeArticles: (articles: Required<Article>[]) => Promise<Article[]>;
  };
}

export const SummaryResponseSchema = z.object({
  summary: z.string().optional(),
  response: z.string().optional(),
  _meta: z.any().optional(),
});
