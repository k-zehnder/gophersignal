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

    // Scrape top stories from Hacker News and fetch their content
    const articles = await articleProcessor.processTopStories();

    // Summarize the content of the fetched articles
    const summarizedArticles = await summarizer.summarizeArticles(articles);

    // Save the summarized articles to the database concurrently
    await Promise.all(
      summarizedArticles.map((article) => db.saveArticle(article))
    );

    // Close the browser
    await browser.close();

    console.log('Script finished successfully.');
  } catch (error) {
    console.error('Error in main function:', error);

    // Close the browser in case of an error
    if (browser) await browser.close();
  } finally {
    // Ensure the database connection is closed
    if (db) await db.closeDatabaseConnection(db.connection);
  }
};

// Execute the main function
main();
