// Assembles high-level services and includes missing dependencies

import { createHackerNewsScraper } from '../services/articleScraper';
import { createContentFetcher } from '../services/articleContent';
import { createArticleProcessor } from '../services/articleProcessor';
import { createArticleSummarizer } from '../services/articleSummarizer';
import { createGitHubService } from '../services/githubService';
import { SummaryResponseSchema } from '../types';
import config from '../config/config';
import { Clients } from '../clients/createClients';

export const createServices = (clients: Clients) => {
  const { browser, instructorClient, articleHelpers, githubClient } = clients;

  const githubService = createGitHubService(githubClient, config.github);
  const scraper = createHackerNewsScraper(browser);
  const contentFetcher = createContentFetcher(browser);
  const articleProcessor = createArticleProcessor(
    scraper,
    contentFetcher,
    articleHelpers
  );
  const articleSummarizer = createArticleSummarizer(
    instructorClient,
    config.ollama,
    SummaryResponseSchema
  );

  return {
    ...clients,
    scraper,
    articleProcessor,
    articleSummarizer,
    githubService,
  };
};

export type Services = ReturnType<typeof createServices>;
