import OpenAI from 'openai';
import config from '../config/config';

// Creates an OpenAI client configured for Ollama
export const createOpenAIClient = () => {
  const openaiClient = new OpenAI({
    apiKey: config.ollama.apiKey || 'ollama',
    baseURL: config.ollama.baseUrl,
  });
  return openaiClient;
};
