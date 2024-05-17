// Scrapes top stories from Hacker News, extracting their titles and links.

const createHackerNewsScraper = (browser) => {
  /**
   * Scrapes top stories from Hacker News and returns an array of articles with title and link.
   */
  const scrapeTopStories = async () => {
    try {
      const page = await browser.newPage();
      await page.setUserAgent(
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36'
      );
      await page.goto('https://news.ycombinator.com/', {
        waitUntil: 'networkidle2',
      });

      // Extract article titles and links from the page
      const articles = await page.evaluate(() => {
        const rows = Array.from(document.querySelectorAll('tr.athing'));
        return rows.map((row) => {
          const titleElement = row.querySelector('.titleline a');
          const title = titleElement
            ? titleElement.innerText
            : 'No title found';
          const link =
            titleElement && titleElement.href
              ? titleElement.href
              : 'No link found';
          return { title, link };
        });
      });

      await page.close();
      console.log('Scraped top stories successfully');
      return articles;
    } catch (error) {
      console.error('Scraping top stories failed:', error);
      return [];
    }
  };

  return { scrapeTopStories };
};

module.exports = { createHackerNewsScraper };
