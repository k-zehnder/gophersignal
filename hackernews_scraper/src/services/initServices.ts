// Creates and returns high-level services using the initialized clients

import { Browser } from 'puppeteer';
import { createHackerNewsScraper } from './articleScraper';
import { createContentFetcher } from './articleContentFetcher';
import { createArticleProcessor } from './articleProcessor';
import { createArticleSummarizer } from './articleSummarizer';
import { createInstructorClient } from '../clients/createInstructorClient';
import config from '../config/config';
import { SummaryResponseSchema } from '../types';
import { Article } from '../types';

export const initServices = ({
  browser,
  instructorClient,
}: {
  browser: Browser;
  instructorClient: ReturnType<typeof createInstructorClient>;
}) => {
  const hackerNewsScraper = createHackerNewsScraper(browser);
  const contentFetcher = createContentFetcher(browser);
  const articleProcessor = createArticleProcessor(
    hackerNewsScraper,
    contentFetcher
  );
  const articleSummarizer = createArticleSummarizer(
    instructorClient,
    config.ollama,
    SummaryResponseSchema
  );

  return {
    articleProcessor,
    articleSummarizer: {
      summarizeArticles: (articles: Required<Article>[]) =>
        articleSummarizer.summarizeArticles(articles),
    },
  };
};
