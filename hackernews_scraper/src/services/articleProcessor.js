// Scrapes top stories from Hacker News, fetches their full content,
// and returns the articles.

const createArticleProcessor = (scraper, contentFetcher) => {
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
