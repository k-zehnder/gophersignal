// Scrapes Hacker News, summarizes content via Ollama, and saves to MySQL.

import dotenv from 'dotenv';
dotenv.config();

import puppeteer from 'puppeteer-extra';
import StealthPlugin from 'puppeteer-extra-plugin-stealth';

import { connectToDatabase } from './database/connection';
import { createHackerNewsScraper } from './services/articleScraper';
import { createHackerNewsScraperDebug } from './services/hackerNewsScraperDebug'; // <-- import debug
import { createContentFetcher } from './services/articleContentFetcher';
import { createArticleProcessor } from './services/articleProcessor';
import { createArticleSummarizer } from './services/articleSummarizer';

import Instructor from '@instructor-ai/instructor';
import OpenAI from 'openai';

import { Article, SummarySchema } from './types/index';
import config from './config/config';

// Initializes dependencies: database, browser, and services
const initDependencies = async () => {
  puppeteer.use(StealthPlugin());

  const browser = await puppeteer.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox'],
  });

  // Connect to DB
  const db = await connectToDatabase(config.mysql);

  // Initialize OpenAI + Instructor
  const openaiClient = new OpenAI({
    apiKey: config.ollama.apiKey || 'ollama',
    baseURL: config.ollama.baseUrl,
  });
  const instructorClient = Instructor({ client: openaiClient, mode: 'JSON' });

  // Normal HN scraper service
  const hackerNewsScraper = createHackerNewsScraper(browser);

  // OPTIONAL: The debug scraper (only used if needed)
  const hackerNewsScraperDebug = createHackerNewsScraperDebug(browser);

  // Rest of your services
  const contentFetcher = createContentFetcher(browser);
  const articleProcessor = createArticleProcessor(
    hackerNewsScraper,
    contentFetcher
  );
  const articleSummarizer = createArticleSummarizer(
    instructorClient,
    config.ollama,
    SummarySchema
  );

  return {
    db,
    browser,
    articleProcessor,
    articleSummarizer,
    hackerNewsScraperDebug,
  };
};

// Main function orchestrating the workflow
const main = async ({
  db,
  browser,
  articleProcessor,
  articleSummarizer,
  hackerNewsScraperDebug,
}: Awaited<ReturnType<typeof initDependencies>>) => {
  try {
    // (1) -- DEBUG STEP --
    console.log('Running debug scrape first...');
    const debugArticles = await hackerNewsScraperDebug.debugScrapeTopStories();
    console.log('Debug articles found:', debugArticles.length);

    // TODO: Decide if you want to exit or proceed with normal scraping

    // // (2) -- Normal workflow below --
    // const articles = await articleProcessor.processTopStories();
    // const articlesWithContent = articles.filter(
    //   (article): article is Required<Article> => article.content != null
    // );

    // const summarizedArticles = await articleSummarizer.summarizeArticles(
    //   articlesWithContent
    // );

    // await db.saveArticles(summarizedArticles);

    console.log('Script finished successfully.');
  } catch (error) {
    console.error('Error in main:', error);
  } finally {
    // Clean up
    await db.closeDatabaseConnection();
    await browser.close();
  }
};

if (require.main === module) {
  initDependencies()
    .then((dependencies) => main(dependencies))
    .catch((error) => {
      console.error('Failed to initialize dependencies:', error);
      process.exit(1);
    });
}

export { main, initDependencies };
