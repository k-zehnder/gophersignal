// Initializes infrastructure clients.

import { createBrowserClient } from '../clients/puppeteer';
import { createMySqlClient } from '../database/mysql';
import { createOpenAIClient } from '../clients/openAI';
import { createInstructorClient } from '../clients/instructor';
import createArticleHelpers from '../utils/article';
import { createTimeUtil } from '../utils/time';
import { type Config } from '../types';

export const createClients = async (config: Config) => {
  const browser = await createBrowserClient();
  const db = await createMySqlClient(config);
  const openaiClient = createOpenAIClient(config);
  const instructorClient = createInstructorClient(openaiClient);
  const articleHelpers = createArticleHelpers();
  const timeUtil = createTimeUtil();

  return {
    browser,
    db,
    openaiClient,
    instructorClient,
    articleHelpers,
    timeUtil,
  };
};

export type Clients = Awaited<ReturnType<typeof createClients>>;
