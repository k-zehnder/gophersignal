// Scrapes top stories from Hacker News, extracting their titles and links.

import { Browser } from 'puppeteer';
import { Article } from '../types/index';

const createHackerNewsScraper = (browser: Browser) => {
  const scrapeTopStories = async (): Promise<Article[]> => {
    try {
      const page = await browser.newPage();
      await page.goto('https://news.ycombinator.com/', {
        waitUntil: 'networkidle2',
      });

      const stories = await page.evaluate(() => {
        const rows = Array.from(document.querySelectorAll('tr.athing'));
        return rows.map((row) => {
          // Extract title and link
          const titleElement = row.querySelector(
            '.titleline a'
          ) as HTMLAnchorElement;
          const title = titleElement?.innerText ?? 'No title';
          const link = titleElement?.href ?? 'No link';

          // Find subtext row
          const subtextRow = row.nextElementSibling;
          if (!subtextRow) {
            return {
              title,
              link,
              upvotes: 0,
              comment_count: 0,
              comment_link: 'No comments link',
            };
          }

          const subtext = subtextRow.querySelector('.subtext');
          if (!subtext) {
            return {
              title,
              link,
              upvotes: 0,
              comment_count: 0,
              comment_link: 'No comments link',
            };
          }

          // Extract upvotes
          let upvotes = 0;
          const scoreEl = subtext.querySelector('.score');
          if (scoreEl?.textContent) {
            const match = scoreEl.textContent.match(/\d+/);
            upvotes = match ? parseInt(match[0], 10) : 0;
          }

          // Extract comment count and link
          let comment_count = 0;
          let comment_link = 'No comments link';
          const commentLinkElement = Array.from(
            subtext.querySelectorAll('a')
          ).find((a) => a.textContent?.includes('comment'));
          if (commentLinkElement) {
            const match = commentLinkElement.textContent?.match(/\d+/);
            comment_count = match ? parseInt(match[0], 10) : 0;
            comment_link = commentLinkElement.href;
          }

          return { title, link, upvotes, comment_count, comment_link };
        });
      });

      await page.close();
      console.log('Scraped stories:', stories);
      return stories;
    } catch (error) {
      console.error('Error during scraping:', error);
      return [];
    }
  };

  return { scrapeTopStories };
};

export { createHackerNewsScraper, Article };
