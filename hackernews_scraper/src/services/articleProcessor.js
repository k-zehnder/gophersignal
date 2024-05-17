// Processes top stories from Hacker News by scraping them, fetching their full content,
// and saving the articles to the database.

const createArticleProcessor = (scraper, contentFetcher, db) => {
  /**
   * Scrapes top stories from Hacker News, fetches their content, and saves them to the database.
   */
  const processTopStories = async () => {
    try {
      const articles = await scraper.scrapeTopStories();
      articles.reverse();

      for (const article of articles) {
        article.content = await contentFetcher.fetchArticleContent(
          article.link
        );
        await new Promise((resolve) => setTimeout(resolve, 1000));
      }

      for (const article of articles) {
        await db.saveArticle(article);
      }
    } catch (error) {
      console.error('Error processing top stories:', error);
    }
  };

  return { processTopStories };
};

module.exports = { createArticleProcessor };
