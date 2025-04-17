// Scrapes a page and returns the articles and next URL. The isTop flag distinguishes
// between the homepage (fresh stories) and /front (dup, dead, flagged articles).

import { Browser } from 'puppeteer';
import { Article, Scraper } from '../types';

const TOP_BASE_URL = 'https://news.ycombinator.com';
const FRONT_BASE_URL = 'https://news.ycombinator.com/front';

// Creates a Hacker News scraper
const createHackerNewsScraper = (browser: Browser): Scraper => {
  // Scrapes a page and returns articles and the next URL; isTop distinguishes / vs /front
  const scrapePageWithNext = async (
    pageUrl: string,
    isTop: boolean = false
  ): Promise<{ articles: Article[]; nextUrl: string | null }> => {
    const page = await browser.newPage();
    try {
      await page.goto(pageUrl, { waitUntil: 'networkidle2' });
      const result = await page.evaluate(
        (isTop: boolean, frontBase: string) => {
          const articles: Article[] = [];
          const rows = Array.from(document.querySelectorAll('tr.athing'));
          rows.forEach((row: Element) => {
            const rankElement = row.querySelector(
              'td.title > span.rank'
            ) as HTMLElement | null;
            const rankStr = rankElement
              ? rankElement.innerText.replace(/\D/g, '').trim()
              : '0';
            const article_rank = parseInt(rankStr, 10) || 0;
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
            // For /front, include only flagged, dead, or duplicate articles
            if (!isTop && !(flagged || dead || dupe)) return;
            let upvotes = 0,
              comment_count = 0;
            let comment_link = 'No comments link';
            const subtextRow = row.nextElementSibling;
            if (subtextRow) {
              const subtext = subtextRow.querySelector('.subtext');
              if (subtext) {
                const scoreEl = subtext.querySelector('.score');
                if (scoreEl?.textContent) {
                  const match = scoreEl.textContent.match(/\d+/);
                  upvotes = match ? parseInt(match[0], 10) : 0;
                }
                const commentLinkElement = Array.from(
                  subtext.querySelectorAll('a') as NodeListOf<HTMLAnchorElement>
                ).find((a) => a.textContent?.includes('comment'));
                if (commentLinkElement) {
                  const match = commentLinkElement.textContent?.match(/\d+/);
                  comment_count = match ? parseInt(match[0], 10) : 0;
                  comment_link = commentLinkElement.href;
                }
              }
            }
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
          // Find the "more" button to get the next URL
          const moreLinkElem = document.querySelector(
            'a.morelink'
          ) as HTMLAnchorElement | null;
          let nextUrl: string | null = null;
          if (moreLinkElem) {
            if (isTop) {
              nextUrl = moreLinkElem.href;
            } else {
              const match = moreLinkElem.href.match(
                /day=(\d{4}-\d{2}-\d{2})&p=(\d+)/
              );
              if (match) {
                const day = match[1];
                const nextPage = parseInt(match[2], 10);
                nextUrl = `${frontBase}?day=${day}&p=${nextPage}`;
              }
            }
          }
          return { articles, nextUrl };
        },
        isTop,
        FRONT_BASE_URL
      );
      return result;
    } catch (error) {
      console.error(`Error scraping ${pageUrl}:`, error);
      return { articles: [], nextUrl: null };
    } finally {
      await page.close();
    }
  };

  // Scrapes top stories from the root (/)
  const scrapeTopStories = async (maxPages?: number): Promise<Article[]> => {
    let articles: Article[] = [];
    let nextUrl: string | null = TOP_BASE_URL;
    let pageCount = 0;
    while (nextUrl) {
      console.log(`Scraping top stories page: ${nextUrl}`);
      const { articles: pageArticles, nextUrl: newNextUrl } =
        await scrapePageWithNext(nextUrl, true);
      articles = articles.concat(pageArticles);
      nextUrl = newNextUrl;
      pageCount++;
      if (maxPages && pageCount >= maxPages) break;
    }
    return articles;
  };

  // Scrapes flagged articles from the front pages (/front)
  const scrapeFront = async (maxPages: number = 10): Promise<Article[]> => {
    let articles: Article[] = [];
    let nextUrl: string | null = FRONT_BASE_URL;
    let pageCount = 0;
    while (nextUrl) {
      console.log(`Scraping front page: ${nextUrl}`);
      const { articles: pageArticles, nextUrl: newNextUrl } =
        await scrapePageWithNext(nextUrl, false);
      articles = articles.concat(pageArticles);
      nextUrl = newNextUrl;
      pageCount++;
      if (pageCount >= maxPages) break;
    }
    return articles;
  };

  // Scrapes front pages for a specific day from (/front)
  const scrapeFrontForDay = async (day: string): Promise<Article[]> => {
    let articles: Article[] = [];
    let nextUrl: string | null = `${FRONT_BASE_URL}?day=${day}&p=1`;
    while (nextUrl) {
      console.log(`Scraping front page for ${day}: ${nextUrl}`);
      const { articles: pageArticles, nextUrl: newNextUrl } =
        await scrapePageWithNext(nextUrl, false);
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
