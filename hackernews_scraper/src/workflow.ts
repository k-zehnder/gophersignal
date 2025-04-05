// Scrapes (https://news.ycombinator.com/front, https://news.ycombinator.com/news)
// processes, summarizes, and saves articles

import { Article } from './types/article';
import { Services } from './services/createServices';

export const createWorkflow = (services: Services) => {
  const MAX_TOP_STORY_PAGES = 2; // Limit for / scraping
  const MAX_SUMMARIZED_ARTICLES = 30; // Cap for regular articles summarization

  // Merge top stories with categorized front-page articles
  const mergeArticles = (
    topArticles: Article[],
    categorized: { flagged: Article[]; dead: Article[]; dupe: Article[] }
  ): Article[] => [
    ...topArticles,
    ...categorized.flagged,
    ...categorized.dead,
    ...categorized.dupe,
  ];

  // Filter and summarize articles, preserving original order
  const filterAndSummarizeArticles = async (
    processed: Article[],
    topArticles: Article[],
    flaggedArticles: Article[]
  ): Promise<Article[]> => {
    const topWithContent =
      services.articleProcessor.helpers.getTopArticlesWithContent(
        processed,
        topArticles
      );
    const flaggedWithContent =
      services.articleProcessor.helpers.getTopArticlesWithContent(
        processed,
        flaggedArticles
      );

    const summarizedTop = await services.articleSummarizer.summarizeArticles(
      topWithContent.slice(0, MAX_SUMMARIZED_ARTICLES)
    );
    const summarizedFlagged =
      await services.articleSummarizer.summarizeArticles(flaggedWithContent);

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

  // Main workflow to scrape, process, summarize, and save articles
  const run = async (): Promise<void> => {
    try {
      const frontArticles = await services.scraper.scrapeFront();
      console.info(`Front articles count: ${frontArticles.length}`);
      const fetchedFrontTitles = frontArticles.map((a) => a.title);
      console.log(frontArticles.map((a) => a.title));

      const categorizedArticles =
        services.articleProcessor.helpers.categorizeArticles(frontArticles);

      const topArticles = await services.scraper.scrapeTopStories(
        MAX_TOP_STORY_PAGES
      );
      console.info(`Top stories scraped: ${topArticles.length}`);

      const allArticles = mergeArticles(topArticles, categorizedArticles);
      console.info(`Total articles to process: ${allArticles.length}`);

      const processedArticles = await services.articleProcessor.processArticles(
        allArticles
      );
      const finalArticles = await filterAndSummarizeArticles(
        processedArticles,
        topArticles,
        categorizedArticles.flagged
      );

      await services.db.saveArticles(finalArticles);
      console.info(
        `Workflow completed. Saved ${finalArticles.length} articles`
      );
    } catch (error) {
      console.error('Workflow execution error:', error);
      throw error;
    }
  };

  // Clean up resources
  const shutdown = async (): Promise<void> => {
    try {
      await Promise.all([
        services.db.closeDatabaseConnection(),
        services.browser.close(),
      ]);
      console.info('Resources released successfully');
    } catch (error) {
      console.error('Error during shutdown:', error);
    }
  };

  return { run, shutdown };
};
