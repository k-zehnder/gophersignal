// Scrapes (https://news.ycombinator.com/front, https://news.ycombinator.com/news),
// processes, summarizes, and saves articles.

import { Article } from './types/article';
import { Services } from './services/createServices';

export class Workflow {
  private readonly NUMBER_OF_TOP_STORY_PAGES = 2;
  private readonly MAX_SUMMARIZED_TOP_ARTICLES = 30;
  private readonly MAX_TOTAL_SUMMARIES = 40;

  constructor(private services: Services) {}

  public async run(): Promise<void> {
    try {
      // Fetch front-page articles from /front
      const combinedFrontArticles = await this.services.scraper.scrapeFront();
      console.info(`Front articles count: ${combinedFrontArticles.length}`);

      // Categorize front-page articles (includes flagged, dead, and duplicate articles)
      const categorizedArticles =
        this.services.articleProcessor.helpers.categorizeArticles(
          combinedFrontArticles
        );

      // Scrape top stories from /news
      const topArticles = await this.services.scraper.scrapeTopStories(
        this.NUMBER_OF_TOP_STORY_PAGES
      );
      console.info(`Top stories scraped: ${topArticles.length}`);

      // Merge categorized articles with top stories for full processing
      const allArticles = this.mergeArticles(topArticles, categorizedArticles);
      console.info(`Total articles to process: ${allArticles.length}`);

      // Process articles to fetch full content.
      const processedArticles =
        await this.services.articleProcessor.processArticles(allArticles);

      // Summarize articles:
      // Summarize top stories up to 30.
      // Then, if we haven't reached 40 total summaries, summarize flagged articles
      // until the overall total is 40.
      const finalArticles = await this.filterAndSummarizeArticles(
        processedArticles,
        topArticles,
        categorizedArticles
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

  // Merge top stories with categorized articles, ensuring order consistency
  private mergeArticles(
    topArticles: Article[],
    categorized: { flagged: Article[]; dead: Article[]; dupe: Article[] }
  ): Article[] {
    return [
      // Reverse top stories to ensure the newest appear first
      ...topArticles.reverse(),
      // Reverse flagged articles for correct ordering
      ...categorized.flagged.reverse(),
      // Reverse dead articles
      ...categorized.dead.reverse(),
      // Reverse duplicate articles
      ...categorized.dupe.reverse(),
    ];
  }

  // Summarize articles:
  // First, summarize top stories (up to 30 articles).
  // Then, summarize flagged articles only up to the remaining amount (so that the total
  // number of summarized articles does not exceed 40).
  // Finally, append any processed articles that weren't summarized.
  private async filterAndSummarizeArticles(
    processed: Article[],
    topArticles: Article[],
    categorized: { flagged: Article[]; dead: Article[]; dupe: Article[] }
  ): Promise<Article[]> {
    // Get top stories with content from processed articles
    const topWithContent =
      this.services.articleProcessor.helpers.getTopArticlesWithContent(
        processed,
        topArticles
      );
    // Summarize up to MAX_SUMMARIZED_TOP_ARTICLES top stories
    const summarizedTop =
      await this.services.articleSummarizer.summarizeArticles(
        topWithContent.slice(0, this.MAX_SUMMARIZED_TOP_ARTICLES)
      );

    // Calculate remaining summaries allowed so that total does not exceed MAX_TOTAL_SUMMARIES 
    const remainingSummaries = this.MAX_TOTAL_SUMMARIES - summarizedTop.length;

    // Get flagged articles with content from processed articles
    const flaggedWithContent =
      this.services.articleProcessor.helpers.getTopArticlesWithContent(
        processed,
        categorized.flagged
      );
    // Summarize flagged articles up to the remaining number
    const summarizedFlagged =
      await this.services.articleSummarizer.summarizeArticles(
        flaggedWithContent.slice(0, remainingSummaries)
      );

    // Build a list of article links that were summarized
    const summarizedLinks = [
      ...summarizedTop.map((a) => a.link),
      ...summarizedFlagged.map((a) => a.link),
    ];

    return [
      ...summarizedTop,
      ...summarizedFlagged,
      // Append any processed articles that weren't summarized
      ...processed.filter((article) => !summarizedLinks.includes(article.link)),
    ];
  }

  public async shutdown(): Promise<void> {
    try {
      // Close database and browser resources
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
