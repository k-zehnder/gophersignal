import { DBClient } from './db';
import { TimeUtil } from '../utils/time';
import { BrowserClient } from '../clients/puppeteer';
import { Scraper, ArticleProcessor, ArticleSummarizer } from './services';

export interface Dependencies {
  db: DBClient;
  browser: BrowserClient;
  timeUtil: TimeUtil;
  instructorClient: any;
  scraper: Scraper;
  articleProcessor: ArticleProcessor;
  articleSummarizer: ArticleSummarizer;
}
