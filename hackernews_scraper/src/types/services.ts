import { Article } from './article';
import { ArticleHelpers } from '../utils/article';

export interface Scraper {
  scrapeTopStories: (numPages?: number) => Promise<Article[]>;
  scrapeFront: (numPages?: number) => Promise<Article[]>;
  scrapeFrontForDay: (day: string) => Promise<Article[]>;
}

export interface ContentFetcher {
  fetchArticleContent: (url: string) => Promise<string>;
}

export interface ArticleProcessor {
  scrapeTopStories: (numPages?: number) => Promise<Article[]>;
  scrapeFrontForDay: (day: string) => Promise<Article[]>;
  processArticles: (articles: Article[]) => Promise<Article[]>;
  helpers: ArticleHelpers;
}

export interface ArticleSummarizer {
  summarizeArticles: (articles: Required<Article>[]) => Promise<Article[]>;
}
