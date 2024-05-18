// Processes top stories from Hacker News by scraping them, fetching their full content,
// and saving the articles to the database.

const createArticleProcessor = (scraper, contentFetcher, db) => {
  /**
   * Scrapes top stories from Hacker News, fetches their content, and returns the articles.
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

      return articles;
    } catch (error) {
      console.error('Error processing top stories:', error);
      return [];
    }
  };

  return { processTopStories };
};

module.exports = { createArticleProcessor };
