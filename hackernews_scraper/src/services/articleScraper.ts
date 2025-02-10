// Scrapes top stories from Hacker News, extracting their titles, links, article_rank, upvotes, and comment counts.

import { Browser } from 'puppeteer';
import { Article, Scraper } from '../types';

// Creates a Hacker News scraper
const createHackerNewsScraper = (browser: Browser): Scraper => {
  // Navigates to a URL, extracts articles, and returns a nextUrl from the "more" button
  const scrapePageWithNext = async (
    pageUrl: string
  ): Promise<{ articles: Article[]; nextUrl: string | null }> => {
    const page = await browser.newPage();
    try {
      await page.goto(pageUrl, { waitUntil: 'networkidle2' });
      const result = await page.evaluate(() => {
        const articles: Article[] = [];
        // Get all rows that represent an article
        const rows = Array.from(document.querySelectorAll('tr.athing'));
        rows.forEach((row: Element) => {
          // Extract the article_rank from the span with class "rank"
          const rankElement = row.querySelector(
            'td.title > span.rank'
          ) as HTMLElement | null;
          // Remove any non-digit characters (such as the trailing dot)
          const rankStr = rankElement
            ? rankElement.innerText.replace(/\D/g, '').trim()
            : '0';
          const article_rank = parseInt(rankStr, 10) || 0;
          // Extract title and link from the title element
          const titleElement = row.querySelector(
            'td.title > span.titleline a'
          ) as HTMLAnchorElement | null;
          const title = titleElement?.innerText || 'No title';
          const link = titleElement?.href || 'No link';
          const titleContainer =
            titleElement?.parentElement as HTMLElement | null;
          const titleText = titleContainer?.innerText || '';
          const flagged = titleText.includes('[flagged]');
          const dead = titleText.includes('[dead]');
          const dupe = titleText.includes('[dupe]');
          let upvotes = 0;
          let comment_count = 0;
          let comment_link = 'No comments link';
          // Get the next row which holds subtext details (like upvotes and comments)
          const subtextRow = row.nextElementSibling;
          if (subtextRow) {
            const subtext = subtextRow.querySelector('.subtext');
            if (subtext) {
              const scoreEl = subtext.querySelector('.score');
              if (scoreEl?.textContent) {
                const match = scoreEl.textContent.match(/\d+/);
                upvotes = match ? parseInt(match[0], 10) : 0;
              }
              // Cast NodeList to NodeListOf<HTMLAnchorElement> so that 'a' is properly typed
              const commentLinkElement = Array.from(
                subtext.querySelectorAll('a') as NodeListOf<HTMLAnchorElement>
              ).find((a: HTMLAnchorElement) =>
                a.textContent?.includes('comment')
              );
              if (commentLinkElement) {
                const match = commentLinkElement.textContent?.match(/\d+/);
                comment_count = match ? parseInt(match[0], 10) : 0;
                comment_link = commentLinkElement.href;
              }
            }
          }
          // Push the article object including the new "article_rank" property
          articles.push({
            title,
            link,
            article_rank,
            flagged,
            dead,
            dupe,
            upvotes,
            comment_count,
            comment_link,
          });
        });
        // Get the "more" button link, ensuring correct pagination
        const moreLinkElem = document.querySelector(
          'a.morelink'
        ) as HTMLAnchorElement | null;
        let nextUrl: string | null = null;
        if (moreLinkElem) {
          const match = moreLinkElem.href.match(
            /day=(\d{4}-\d{2}-\d{2})&p=(\d+)/
          );
          if (match) {
            const day = match[1];
            const nextPage = parseInt(match[2], 10);
            nextUrl = `https://news.ycombinator.com/front?day=${day}&p=${nextPage}`;
          }
        }
        return { articles, nextUrl };
      });
      return result;
    } catch (error) {
      console.error(`Error scraping ${pageUrl}:`, error);
      return { articles: [], nextUrl: null };
    } finally {
      await page.close();
    }
  };

  // Scrapes top stories with an optional page limit
  const scrapeTopStories = async (maxPages?: number): Promise<Article[]> => {
    let articles: Article[] = [];
    let nextUrl: string | null = 'https://news.ycombinator.com/';
    let pageCount = 0;
    while (nextUrl) {
      console.log(`Scraping top stories page: ${nextUrl}`);
      const { articles: pageArticles, nextUrl: newNextUrl } =
        await scrapePageWithNext(nextUrl);
      articles = articles.concat(pageArticles);
      nextUrl = newNextUrl;
      pageCount++;
      if (maxPages && pageCount >= maxPages) break;
    }
    return articles;
  };

  // Scrapes front pages with a page limit
  const scrapeFront = async (maxPages: number = 10): Promise<Article[]> => {
    let articles: Article[] = [];
    let nextUrl: string | null = 'https://news.ycombinator.com/front';
    let pageCount = 0;
    while (nextUrl) {
      console.log(`Scraping front page: ${nextUrl}`);
      const { articles: pageArticles, nextUrl: newNextUrl } =
        await scrapePageWithNext(nextUrl);
      articles = articles.concat(pageArticles);
      nextUrl = newNextUrl;
      pageCount++;
      if (pageCount >= maxPages) break;
    }
    return articles;
  };

  // Retrieves all front pages for a given day by starting at the day-specific URL
  const scrapeFrontForDay = async (day: string): Promise<Article[]> => {
    let articles: Article[] = [];
    let nextUrl:
      | string
      | null = `https://news.ycombinator.com/front?day=${day}&p=1`;
    while (nextUrl) {
      console.log(`Scraping front page for ${day}: ${nextUrl}`);
      const { articles: pageArticles, nextUrl: newNextUrl } =
        await scrapePageWithNext(nextUrl);
      articles = articles.concat(pageArticles);
      nextUrl = newNextUrl;
    }
    return articles;
  };

  return {
    scrapeTopStories,
    scrapeFront,
    scrapeFrontForDay,
  };
};

export { createHackerNewsScraper };
