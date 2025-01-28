import Instructor from '@instructor-ai/instructor';
import { createOpenAIClient } from './createOpenAIClient';

// Wraps the OpenAI client with the Instructor library
export const createInstructorClient = (
  openaiClient: ReturnType<typeof createOpenAIClient>
) => {
  const instructorClient = Instructor({
    client: openaiClient,
    mode: 'JSON',
  });
  return instructorClient;
};
