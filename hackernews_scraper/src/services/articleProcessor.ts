// Scrapes top stories from Hacker News, fetches their full content,
// and returns the articles.

import { Article } from './articleScraper';

type Scraper = {
  scrapeTopStories: () => Promise<Article[]>;
};

type ContentFetcher = {
  fetchArticleContent: (url: string) => Promise<string>;
};

const createArticleProcessor = (
  scraper: Scraper,
  contentFetcher: ContentFetcher
) => {
  const processTopStories = async (): Promise<Article[]> => {
    try {
      const articles = await scraper.scrapeTopStories();
      articles.reverse();

      for (const article of articles) {
        article.content = await contentFetcher.fetchArticleContent(
          article.link
        );
        await new Promise((resolve) => setTimeout(resolve, 1000));
      }

      return articles;
    } catch (error) {
      console.error('Error processing top stories:', error);
      return [];
    }
  };

  return { processTopStories };
};

export { createArticleProcessor };
