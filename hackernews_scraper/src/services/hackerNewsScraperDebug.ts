import { Browser } from 'puppeteer';

export interface DebugArticle {
  title: string;
  link: string;
}

export function createHackerNewsScraperDebug(browser: Browser) {
  const debugScrapeTopStories = async (): Promise<DebugArticle[]> => {
    try {
      const page = await browser.newPage();
      await page.goto('https://news.ycombinator.com/', {
        waitUntil: 'networkidle2',
      });

      // Grab all 'tr.athing' rows
      const rowHandles = await page.$$('tr.athing');
      console.log(`Found ${rowHandles.length} 'tr.athing' rows.`);

      const articles: DebugArticle[] = [];

      // Iterate each row
      for (const [idx, rowHandle] of rowHandles.entries()) {
        // Find '.titleline a'
        const titleLink = await rowHandle.$('.titleline a');
        if (!titleLink) {
          console.log(`Row #${idx}: NO '.titleline a' found.`);
          continue;
        }

        // Extract text & href
        const title = await page.evaluate((el) => el.innerText, titleLink);
        const link = await page.evaluate((el) => el.href, titleLink);

        console.log(`Row #${idx} title: ${title}`);
        console.log(`Row #${idx} link: ${link}`);

        articles.push({ title, link });
      }

      await page.close();
      console.log('Debug scrape finished.');
      return articles;
    } catch (error) {
      console.error('Debug scrape error:', error);
      return [];
    }
  };

  return { debugScrapeTopStories };
}
