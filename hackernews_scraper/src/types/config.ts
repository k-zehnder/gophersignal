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
