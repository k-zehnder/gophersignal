import { z } from 'zod';

export const MySQLConfigSchema = z.object({
  host: z.string().default('localhost'),
  port: z.coerce
    .number({ invalid_type_error: 'MYSQL_PORT must be a number' })
    .default(3306),
  user: z.string().default('user'),
  password: z.string().default(''),
  database: z.string().default('database_name'),
});
export type MySQLConfig = z.infer<typeof MySQLConfigSchema>;

export const OllamaConfigSchema = z.object({
  baseUrl: z
    .string()
    .url('OLLAMA_BASE_URL must be a valid URL')
    .default('http://localhost:11434/api/generate'),
  model: z.string().default('llama3:instruct'),
  apiKey: z.string().optional(),
  maxContentLength: z.coerce
    .number({ invalid_type_error: 'MAX_CONTENT_LENGTH must be a number' })
    .default(2000),
  maxSummaryLength: z.coerce
    .number({ invalid_type_error: 'MAX_SUMMARY_LENGTH must be a number' })
    .default(500),
});
export type OllamaConfig = z.infer<typeof OllamaConfigSchema>;

export const ConfigSchema = z.object({
  mysql: MySQLConfigSchema,
  ollama: OllamaConfigSchema,
});
export type Config = z.infer<typeof ConfigSchema>;
