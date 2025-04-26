// Initializes infrastructure clients.

import { createBrowserClient } from '../clients/puppeteer';
import { createMySqlClient } from '../database/mysql';
import { createOpenAIClient } from '../clients/openAI';
import { createInstructorClient } from '../clients/instructor';
import { createGitHubClient } from './github';
import createArticleHelpers from '../utils/article';
import { createTimeUtil } from '../utils/time';
import { Config } from '../types/config';

export const createClients = async (config: Config) => {
  const browser = await createBrowserClient();
  const db = await createMySqlClient(config);
  const openaiClient = createOpenAIClient();
  const instructorClient = createInstructorClient(openaiClient);
  const githubClient = createGitHubClient();
  const articleHelpers = createArticleHelpers();
  const timeUtil = createTimeUtil();

  return {
    browser,
    db,
    openaiClient,
    instructorClient,
    githubClient,
    articleHelpers,
    timeUtil,
  };
};

export type Clients = Awaited<ReturnType<typeof createClients>>;
