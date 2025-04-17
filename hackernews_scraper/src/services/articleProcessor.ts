import { Article, Scraper, ContentFetcher } from '../types';
import { ArticleHelpers } from '../utils/article';

const createArticleProcessor = (
  scraper: Scraper,
  contentFetcher: ContentFetcher,
  helpers: ArticleHelpers
) => {
  // Scrapes top stories using page-number based pagination or next-button logic
  const scrapeTopStories = async (numPages?: number): Promise<Article[]> => {
    return await scraper.scrapeTopStories(numPages);
  };

  // Processes articles by fetching full content
  const processArticles = async (articles: Article[]): Promise<Article[]> => {
    for (const article of articles) {
      try {
        article.content = await contentFetcher.fetchArticleContent(
          article.link
        );
      } catch (error) {
        console.error(`Error processing article at ${article.link}:`, error);
      }
      await new Promise((resolve) => setTimeout(resolve, 1000));
    }
    return articles;
  };

  return {
    scrapeTopStories,
    processArticles,
    helpers,
  };
};

export { createArticleProcessor };
