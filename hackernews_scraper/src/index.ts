// Main script to scrape articles from Hacker News, fetch their content,
// summarize the content using Instructor and Ollama, and save the articles
// and their summaries to a MySQL database.

import dotenv from 'dotenv';
dotenv.config();

import puppeteer from 'puppeteer-extra';
import { Browser } from 'puppeteer';
import StealthPlugin from 'puppeteer-extra-plugin-stealth';
import { connectToDatabase } from './database/connection';
import { createHackerNewsScraper } from './services/articleScraper';
import { createContentFetcher } from './services/articleContentFetcher';
import { createArticleProcessor } from './services/articleProcessor';
import { createArticleSummarizer } from './services/articleSummarizer';
import { Article } from './types/index';
import config from './config/config';

// Apply the stealth plugin to puppeteer
puppeteer.use(StealthPlugin());

/**
 * Main function to:
 * 1. Scrape top stories from Hacker News.
 * 2. Fetch the full content of the top stories.
 * 3. Summarize the fetched content using Instructor and Ollama.
 * 4. Save the articles and their summaries to a MySQL database.
 */
const main = async () => {
  let db: any;
  let browser: Browser | null = null;

  try {
    // Initialize the database connection
    db = await connectToDatabase(config.mysql);

    // Launch a new Puppeteer browser instance with Stealth plugin
    browser = await puppeteer.launch({
      headless: true,
      args: ['--no-sandbox', '--disable-setuid-sandbox'],
    });

    // Initialize the components with dependencies
    const hackerNewsScraper = createHackerNewsScraper(browser);
    const contentFetcher = createContentFetcher(browser);
    const articleProcessor = createArticleProcessor(
      hackerNewsScraper,
      contentFetcher
    );
    const articleSummarizer = createArticleSummarizer(config.ollama);

    // Scrape top stories from Hacker News and fetch their content
    const articles = await articleProcessor.processTopStories();

    // Filter articles with defined content
    const articlesWithContent = articles.filter(
      (article): article is Required<Article> => article.content !== undefined
    );

    // Summarize the content of the fetched articles
    const summarizedArticles = await articleSummarizer.summarizeArticles(
      articlesWithContent
    );

    // Save the summarized articles to the database
    await db.saveArticles(summarizedArticles);

    // Close the browser
    await browser.close();

    console.log('Script finished successfully.');
  } catch (error) {
    console.error('Error in main function:', error);

    // Close the browser in case of an error
    if (browser) {
      await browser.close();
    }
  } finally {
    // Ensure the database connection is closed
    if (db) {
      await db.closeDatabaseConnection();
    }
  }
};

// Execute the main function
main();
