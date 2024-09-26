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
}

export interface Config {
  mysql: MySQLConfig;
  ollama: OllamaConfig;
}
