import OpenAI from 'openai';
import { type Config } from '../types';

// Creates an OpenAI client configured for Ollama
export const createOpenAIClient = (config: Config) => {
  const openaiClient = new OpenAI({
    apiKey: config.ollama.apiKey || 'ollama',
    baseURL: config.ollama.baseUrl,
  });
  return openaiClient;
};
