// Contains main functions for processing and summarizing articles.

const { scrapeHackerNews } = require('./articleScraper');
const { fetchArticleContent } = require('./articleContentFetcher');
const { summarizeContentWithRetry } = require('./articleSummarizer');
const {
  saveArticle,
  updateArticleSummary,
  fetchUnsummarizedArticles,
} = require('../database/queries');
const { logger, delay } = require('../helpers/scraperUtils');

// Scrapes articles from Hacker News, fetches their content, and saves them to the database.
const processArticles = async (browser, connection) => {
  // Scrape articles from Hacker News
  let articles = await scrapeHackerNews(browser);
  articles.reverse();

  // Fetch content for each article
  for (const article of articles) {
    article.content = await fetchArticleContent(article.link, browser);
    await delay(1000);
  }

  // Save articles to the database
  for (const article of articles) {
    await saveArticle(connection, article);
  }
};

// Fetches unsummarized articles from the database, summarizes them, and updates the database.
const summarizeArticles = async (connection, config) => {
  // Fetch unsummarized articles from the database
  const rows = await fetchUnsummarizedArticles(connection);
  for (const { id, content } of rows) {
    if (!content) {
      logger.warn(`Skipping Article ID ${id}: content is empty`);
      continue;
    }

    // Summarize content and update the database
    const summary = await summarizeContentWithRetry(
      config.huggingFace.apiUrl,
      config.huggingFace.apiKey,
      content
    );
    await updateArticleSummary(connection, id, summary);
    await delay(1000);
  }
};

module.exports = { processArticles, summarizeArticles };
