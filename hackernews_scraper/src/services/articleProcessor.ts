import { z } from 'zod';
import { ArticleSchema, type Article, Scraper, ContentFetcher } from '../types';
import { ArticleHelpers } from '../utils/article';

const createArticleProcessor = (
  scraper: Scraper,
  fetcher: ContentFetcher,
  helpers: ArticleHelpers
) => {
  // Fetch top stories and validate against ArticleSchema
  const scrapeTopStories = async (numPages?: number): Promise<Article[]> => {
    const raw = await scraper.scrapeTopStories(numPages);
    const parsed = z.array(ArticleSchema).safeParse(raw);
    if (!parsed.success) {
      console.error('Top stories validation failed:', parsed.error.format());
      return [];
    }
    return parsed.data;
  };

  // Fetch full content for each article and re-validate
  const processArticles = async (articles: Article[]): Promise<Article[]> => {
    const output: Article[] = [];
    for (const art of articles) {
      let enriched: Article = art;
      try {
        const content = await fetcher.fetchArticleContent(art.link);
        enriched = ArticleSchema.parse({ ...art, content });
      } catch (e) {
        console.error(`Error fetching content for ${art.link}:`, e);
      }
      output.push(enriched);
      await new Promise((r) => setTimeout(r, 1000));
    }
    return output;
  };

  return { scrapeTopStories, processArticles, helpers };
};

export { createArticleProcessor };
