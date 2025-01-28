// Assembles low-level infrastructure clients into a single object.

import { createBrowserClient } from './createBrowserClient';
import { createDBClient } from './createDBClient';
import { createOpenAIClient } from './createOpenAIClient';
import { createInstructorClient } from './createInstructorClient';
import config from '../config/config';

export const initClients = async () => {
  const browser = await createBrowserClient();
  const db = await createDBClient(config);
  const openaiClient = createOpenAIClient();
  const instructorClient = createInstructorClient(openaiClient);

  return { browser, db, openaiClient, instructorClient };
};
