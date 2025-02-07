import Instructor from '@instructor-ai/instructor';
import { createOpenAIClient } from './openAI';

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
