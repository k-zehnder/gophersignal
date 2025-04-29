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
  numCtx: number; // Context window size for Ollama model, read from OLLAMA_CONTEXT_LENGTH env var
}

export interface GitHubConfig {
  token?: string;
  owner: string;
  repo: string;
  branch: string;
}

export interface Config {
  mysql: MySQLConfig;
  ollama: OllamaConfig;
  github: GitHubConfig;
}
