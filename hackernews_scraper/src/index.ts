// Scrapes Hacker News, summarizes content via Ollama, and saves to MySQL.

import dotenv from 'dotenv';
dotenv.config();

import puppeteer from 'puppeteer-extra';
import StealthPlugin from 'puppeteer-extra-plugin-stealth';
import { connectToDatabase } from './database/connection';
import { createHackerNewsScraper } from './services/articleScraper';
import { createContentFetcher } from './services/articleContentFetcher';
import { createArticleProcessor } from './services/articleProcessor';
import { createArticleSummarizer } from './services/articleSummarizer';
import Instructor from '@instructor-ai/instructor';
import OpenAI from 'openai';
import { Article, SummaryResponseSchema } from './types/index';
import config from './config/config';

// Initializes dependencies: database, browser, and services
const initDependencies = async () => {
  // Initialize Puppeteer with Stealth Plugin
  puppeteer.use(StealthPlugin());
  const browser = await puppeteer.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox'],
  });

  // Connect to the database
  const db = await connectToDatabase(config.mysql);

  // Initialize OpenAI and Instructor clients
  const openaiClient = new OpenAI({
    apiKey: config.ollama.apiKey || 'ollama',
    baseURL: config.ollama.baseUrl,
  });

  const instructorClient = Instructor({
    client: openaiClient,
    mode: 'JSON',
  });

  // Initialize services
  const hackerNewsScraper = createHackerNewsScraper(browser);
  const contentFetcher = createContentFetcher(browser);
  const articleProcessor = createArticleProcessor(
    hackerNewsScraper,
    contentFetcher
  );
  const articleSummarizer = createArticleSummarizer(
    instructorClient,
    config.ollama,
    SummaryResponseSchema
  );

  return { db, browser, articleProcessor, articleSummarizer };
};

// Main function orchestrating the workflow
const main = async ({
  db,
  browser,
  articleProcessor,
  articleSummarizer,
}: Awaited<ReturnType<typeof initDependencies>>) => {
  try {
    // Process articles
    const articles = await articleProcessor.processTopStories();
    const articlesWithContent = articles.filter(
      (article): article is Required<Article> => article.content != null
    );

    // Summarize articles
    const summarizedArticles = await articleSummarizer.summarizeArticles(
      articlesWithContent
    );

    // Save summaries to the database
    await db.saveArticles(summarizedArticles);

    console.log('Script finished successfully.');
  } catch (error) {
    console.error('Error in main:', error);
  } finally {
    // Clean up resources
    await db.closeDatabaseConnection();
    await browser.close();
  }
};

// Run the script
if (require.main === module) {
  initDependencies()
    .then((dependencies) => main(dependencies))
    .catch((error) => {
      console.error('Failed to initialize dependencies:', error);
      process.exit(1);
    });
}

export { main, initDependencies };
