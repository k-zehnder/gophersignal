// Scrapes (https://news.ycombinator.com/front, https://news.ycombinator.com/news),
// processes, summarizes, and saves articles.

import { Article } from './types/article';
import { Services } from './services/createServices';

export const createWorkflow = (services: Services) => {
  const MAX_TOP_STORY_PAGES = 2; // Maximum /news pages to scrape
  const MAX_SUMMARIZED_ARTICLES = 30; // Maximum number of articles to summarize

  // Merge top stories with categorized articles, ensuring order consistency
  const mergeArticles = (
    topArticles: Article[],
    categorized: { flagged: Article[]; dead: Article[]; dupe: Article[] }
  ): Article[] => {
    return [
      ...topArticles,
      ...categorized.flagged,
      ...categorized.dead,
      ...categorized.dupe,
    ];
  };

  // Filter and summarize articles, summarized top stories will be saved last
  const filterAndSummarizeArticles = async (
    processed: Article[],
    topArticles: Article[],
    flaggedArticles: Article[]
  ): Promise<Article[]> => {
    // Extract top articles that have content
    const topWithContent =
      services.articleProcessor.helpers.getTopArticlesWithContent(
        processed,
        topArticles
      );

    // Extract flagged articles that have content
    const flaggedWithContent =
      services.articleProcessor.helpers.getTopArticlesWithContent(
        processed,
        flaggedArticles
      );

    // Summarize top articles (limit to MAX_SUMMARIZED_ARTICLES)
    const summarizedTop = await services.articleSummarizer.summarizeArticles(
      topWithContent.slice(0, MAX_SUMMARIZED_ARTICLES)
    );

    // Summarize all flagged articles
    const summarizedFlagged =
      await services.articleSummarizer.summarizeArticles(flaggedWithContent);

    // Build list of summarized article links
    const summarizedLinks = [
      ...summarizedTop.map((a: Article) => a.link),
      ...summarizedFlagged.map((a: Article) => a.link),
    ];

    return [
      ...processed
        .filter((article) => !summarizedLinks.includes(article.link))
        .reverse(),
      ...summarizedFlagged.reverse(),
      ...summarizedTop.reverse(),
    ];
  };

  // Runs the entire workflow: scrape, process, summarize, and save articles
  const run = async (): Promise<void> => {
    try {
      // Scrape front-page articles from /front
      const combinedFrontArticles = await services.scraper.scrapeFront();
      console.info(`Front articles count: ${combinedFrontArticles.length}`);

      // Categorize front-page articles
      const categorizedArticles =
        services.articleProcessor.helpers.categorizeArticles(
          combinedFrontArticles
        );

      // Scrape top articles from /news
      const topArticles = await services.scraper.scrapeTopStories(
        MAX_TOP_STORY_PAGES
      );
      console.info(`Top stories scraped: ${topArticles.length}`);

      // Merge categorized articles with top articles
      const allArticles = mergeArticles(topArticles, categorizedArticles);
      console.info(`Total articles to process: ${allArticles.length}`);

      // Process articles to fetch full content
      const processedArticles = await services.articleProcessor.processArticles(
        allArticles
      );

      // Filter and summarize articles and order them so that unprocessed come first,
      // then flagged, and top stories last (ensuring the top article is saved last)
      const finalArticles = await filterAndSummarizeArticles(
        processedArticles,
        topArticles,
        categorizedArticles.flagged
      );

      // Save final articles to the database
      await services.db.saveArticles(finalArticles);
      console.info(
        `Workflow completed. Saved ${finalArticles.length} articles.`
      );
    } catch (error) {
      console.error('Workflow execution error:', error);
      throw error;
    }
  };

  // Releases resources such as database connections and the browser
  const shutdown = async (): Promise<void> => {
    try {
      await Promise.all([
        services.db.closeDatabaseConnection(),
        services.browser.close(),
      ]);
      console.info('Resources released successfully.');
    } catch (error) {
      console.error('Error during shutdown:', error);
    }
  };

  return { run, shutdown };
};
