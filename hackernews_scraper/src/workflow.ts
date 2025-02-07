// Scrapes (https://news.ycombinator.com/front, https://news.ycombinator.com/news), processes, summarizes, and saves articles.

import { Article } from './types/article';
import { Services } from './services/createServices';

export class Workflow {
  private readonly NUMBER_OF_TOP_STORY_PAGES = 2;
  private readonly MAX_SUMMARIZED_ARTICLES = 30;

  constructor(private services: Services) {}

  public async run(): Promise<void> {
    try {
      // Fetch front-page articles from /front (today & yesterday)
      const combinedFrontArticles = await this.fetchFrontArticles();
      console.info(`Front articles count: ${combinedFrontArticles.length}`);

      // Categorize front-page articles
      const categorizedArticles =
        this.services.articleProcessor.helpers.categorizeArticles(
          combinedFrontArticles
        );

      // Scrape top stories from /news
      const topArticles = await this.services.scraper.scrapeTopStories(
        this.NUMBER_OF_TOP_STORY_PAGES
      );
      console.info(`Top stories scraped: ${topArticles.length}`);

      // Merge top stories with categorized front-page articles
      const allArticles = this.mergeArticles(topArticles, categorizedArticles);
      console.info(`Total articles to process: ${allArticles.length}`);

      // Process articles to fetch full content
      const processedArticles =
        await this.services.articleProcessor.processArticles(allArticles);

      // Filter and summarize top articles
      const finalArticles = await this.filterAndSummarizeArticles(
        processedArticles,
        topArticles
      );

      // Save final articles to the database
      await this.services.db.saveArticles(finalArticles);
      console.info(
        `Workflow completed. Saved ${finalArticles.length} articles.`
      );
    } catch (error) {
      console.error('Workflow execution error:', error);
      throw error;
    }
  }

  private async fetchFrontArticles(): Promise<Article[]> {
    const { scraper, timeUtil } = this.services;
    const [todayArticles, yesterdayArticles] = await Promise.all([
      scraper.scrapeFrontForDay(timeUtil.today),
      scraper.scrapeFrontForDay(timeUtil.yesterday),
    ]);
    return [...todayArticles, ...yesterdayArticles];
  }

  private mergeArticles(
    topArticles: Article[],
    categorized: { flagged: Article[]; dead: Article[]; dupe: Article[] }
  ): Article[] {
    return [
      ...topArticles,
      ...categorized.flagged,
      ...categorized.dead,
      ...categorized.dupe,
    ];
  }

  private async filterAndSummarizeArticles(
    processed: Article[],
    topArticles: Article[]
  ): Promise<Article[]> {
    const topWithContent =
      this.services.articleProcessor.helpers.getTopArticlesWithContent(
        processed,
        topArticles
      );
    const summarized = await this.services.articleSummarizer.summarizeArticles(
      topWithContent.slice(0, this.MAX_SUMMARIZED_ARTICLES).reverse()
    );
    return [
      ...summarized,
      ...processed.filter(
        (article) => !topArticles.some((top) => top.link === article.link)
      ),
    ];
  }

  public async shutdown(): Promise<void> {
    try {
      await Promise.all([
        this.services.db.closeDatabaseConnection(),
        this.services.browser.close(),
      ]);
      console.info('Resources released successfully.');
    } catch (error) {
      console.error('Error during shutdown:', error);
    }
  }
}
