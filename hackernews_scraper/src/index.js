// Main script to scrape articles from Hacker News, fetch their content,
// summarize the content using the Hugging Face API, and save the articles
// and their summaries to a MySQL database.

require('dotenv').config();

const axios = require('axios');
const puppeteer = require('puppeteer');
const { connectToDatabase } = require('./database/connection');
const { createHackerNewsScraper } = require('./services/articleScraper');
const { createContentFetcher } = require('./services/articleContentFetcher');
const { createArticleProcessor } = require('./services/articleProcessor');
const { createSummarizer } = require('./services/articleSummarizer');
const config = require('./config/config');

/**
 * Main function to:
 * 1. Scrape top stories from Hacker News.
 * 2. Fetch the full content of the top stories.
 * 3. Summarize the fetched content using the Hugging Face API.
 * 4. Save the articles and their summaries to a MySQL database.
 */
const main = async () => {
  let db;
  let browser;

  try {
    // Initialize the database connection
    db = await connectToDatabase(config);

    // Launch a new Puppeteer browser instance
    browser = await puppeteer.launch({
      headless: 'new',
      args: ['--no-sandbox', '--disable-setuid-sandbox'],
    });

    // Initialize the components with dependencies
    const hackerNewsScraper = createHackerNewsScraper(browser);
    const contentFetcher = createContentFetcher(browser);
    const articleProcessor = createArticleProcessor(
      hackerNewsScraper,
      contentFetcher,
      db
    );
    const summarizer = createSummarizer(axios, config, db);

    // Begin a transaction to keep data consistent during updates.
    // This prevents indeterminate database states and ensures Cloudflare
    // doesn't cache incomplete data, maintaining synchronization between
    // the database and cached content.
    await db.connection.beginTransaction();

    // Scrape top stories from Hacker News and fetch their content
    await articleProcessor.processTopStories();

    // Summarize the content of the fetched articles
    await summarizer.summarizeFetchedArticles();

    // Commit the transaction to finalize the changes
    await db.connection.commit();

    // Close the browser
    await browser.close();

    console.log('Script finished successfully.');
  } catch (error) {
    console.error('Error in main function:', error);

    // Roll back the transaction and close the browser in case of an error
    if (db) await db.connection.rollback();
    if (browser) await browser.close();
  } finally {
    // Ensure the database connection is closed
    if (db) await db.closeDatabaseConnection(db.connection);
  }
};

// Execute the main function
main();
