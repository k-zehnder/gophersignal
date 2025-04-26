import { DBClient } from './db';
import { TimeUtil } from '../utils/time';
import { InstructorClient } from '../clients/instructor';
import { BrowserClient } from '../clients/puppeteer';
import {
  Scraper,
  ArticleProcessor,
  ArticleSummarizer,
  GitHubService,
} from './services';

export interface Dependencies {
  db: DBClient;
  browser: BrowserClient;
  timeUtil: TimeUtil;
  instructorClient: InstructorClient;
  scraper: Scraper;
  articleProcessor: ArticleProcessor;
  articleSummarizer: ArticleSummarizer;
  githubService: GitHubService;
}
