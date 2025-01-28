// Orchestrates article scraping, summarizing, and saving articles.

import { initClients } from '../src/clients/initClients';
import { initServices } from '../src/services/initServices';
import { Dependencies, Article } from './types';

// Status codes for success and failure
const STATUS_SUCCESS = 0;
const STATUS_FAILURE = 1;

// Core workflow
export const orchestrateWorkflow = async ({
  db,
  browser,
  articleProcessor,
  articleSummarizer,
}: Dependencies): Promise<number> => {
  try {
    // Fetch raw articles
    const articles = await articleProcessor.processTopStories();

    // Filter out articles without content
    const articlesWithContent = articles.filter(
      (article): article is Required<Article> => Boolean(article.content)
    );

    // Summarize the filtered articles
    const summarizedArticles = await articleSummarizer.summarizeArticles(
      articlesWithContent
    );

    // Save the summarized articles to the database
    await db.saveArticles(summarizedArticles);

    console.log('Workflow completed successfully.');
    return STATUS_SUCCESS;
  } catch (error) {
    console.error('Error in orchestrateWorkflow:', error);
    return STATUS_FAILURE;
  } finally {
    await db
      .closeDatabaseConnection()
      .catch((err) => console.error('Error closing database:', err));
    await browser
      .close()
      .catch((err) => console.error('Error closing browser:', err));
  }
};

// Entry point
export const main = async (): Promise<void> => {
  try {
    // Initialize clients
    const { db, browser, instructorClient } = await initClients();

    // Initialize services
    const { articleProcessor, articleSummarizer } = initServices({
      browser,
      instructorClient,
    });

    // Run the main workflow
    const statusCode = await orchestrateWorkflow({
      db,
      browser,
      articleProcessor,
      articleSummarizer,
    });

    process.exit(statusCode);
  } catch (err) {
    console.error('Failed to execute main:', err);
    process.exit(STATUS_FAILURE);
  }
};

if (require.main === module) {
  main();
}
