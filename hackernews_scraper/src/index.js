// Main script to scrape articles from Hacker News, fetch their content,
// summarize the content using the Hugging Face API, and save the articles
// and their summaries to a MySQL database.

require('dotenv').config();

const puppeteer = require('puppeteer');
const {
  initializeDatabase,
  closeDatabaseConnection,
} = require('./database/connection');
const {
  processArticles,
  summarizeArticles,
} = require('./services/mainFunctions');
const { logger } = require('./helpers/scraperUtils');
const config = require('./config/config');

// Main function to execute the scraping, summarizing, and saving of articles.
const main = async () => {
  let connection;

  try {
    // Initialize the database connection
    connection = await initializeDatabase(config);

    // Launch a new Puppeteer browser instance
    const browser = await puppeteer.launch({ headless: 'new' });

    // Begin a database transaction
    await connection.beginTransaction();

    // Scrape and save articles
    await processArticles(browser, connection);

    // Summarize articles and update the database
    await summarizeArticles(connection, config);

    // Commit the transaction
    await connection.commit();

    // Close the browser
    await browser.close();

    logger.info('Script finished successfully.');
  } catch (error) {
    logger.error('Error in main function:', error);

    // Roll back the transaction in case of an error
    if (connection) await connection.rollback();
  } finally {
    // Ensure the database connection is closed
    await closeDatabaseConnection(connection);
  }
};

// Execute the main function
main();
