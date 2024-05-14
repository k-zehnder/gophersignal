// Scrapes articles from Hacker News and returns an array of articles with title and link.

// Scrapes articles from Hacker News and returns an array of articles with title and link.
const scrapeHackerNews = async (browser) => {
  try {
    const page = await browser.newPage();
    await page.setUserAgent(
      'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36'
    );
    await page.goto('https://news.ycombinator.com/', {
      waitUntil: 'networkidle2',
    });

    // Extracts article titles and links from the page
    const articles = await page.evaluate(() => {
      const rows = Array.from(document.querySelectorAll('tr.athing'));
      return rows.map((row) => {
        const titleElement = row.querySelector('.titleline a');
        const title = titleElement ? titleElement.innerText : 'No title found';
        const link =
          titleElement && titleElement.href
            ? titleElement.href
            : 'No link found';
        return { title, link };
      });
    });

    await page.close();
    console.log('Scraped articles successfully');
    return articles;
  } catch (error) {
    console.error('Scraping failed:', error);
    return [];
  }
};

module.exports = { scrapeHackerNews };
