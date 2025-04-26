// Scrapes a page and returns articles and the next URL. The isTop flag distinguishes
// between the homepage (fresh stories) and /front (dup, dead, flagged articles).

import { Browser } from 'puppeteer';
import { Article, Scraper } from '../types';

const TOP_BASE_URL = 'https://news.ycombinator.com';
const FRONT_BASE_URL = 'https://news.ycombinator.com/front';

type ScrapeResult = { articles: Article[]; nextUrl: string | null };

const createHackerNewsScraper = (browser: Browser): Scraper => {
  // Scrapes a page and returns articles and the next URL; isTop distinguishes / vs /front
  const scrapePageWithNext = async (
    pageUrl: string,
    isTop = false
  ): Promise<ScrapeResult> => {
    const page = await browser.newPage();
    try {
      await page.goto(pageUrl, { waitUntil: 'networkidle2' });

      const result = await page.evaluate(
        (isTopFlag, frontBase) => {
          const articles: Article[] = [];
          const rows = Array.from(document.querySelectorAll('tr.athing'));

          rows.forEach((row) => {
            // Get HN ID or fallback to comment link
            const rawId = (row as HTMLElement).id;
            let hn_id = parseInt(rawId, 10) || 0;
            if (!hn_id) {
              const href =
                row.nextElementSibling?.querySelector<HTMLAnchorElement>(
                  'a[href*="item?id="]'
                )?.href;
              const match = href?.match(/item\?id=(\d+)/);
              hn_id = match ? parseInt(match[1], 10) : 0;
            }
            if (!hn_id) return;

            // Parse rank
            const rankEl = row.querySelector(
              'td.title > span.rank'
            ) as HTMLElement | null;
            const rankText = rankEl?.innerText.replace(/\D/g, '').trim() || '0';
            const article_rank = parseInt(rankText, 10) || 0;

            // Parse title and link
            const titleEl = row.querySelector(
              'td.title > span.titleline a'
            ) as HTMLAnchorElement | null;
            const title = titleEl?.innerText || 'No title';
            const link = titleEl?.href || 'No link';

            // Determine flags
            const containerText = titleEl?.parentElement?.innerText || '';
            const flagged = containerText.includes('[flagged]');
            const dead = containerText.includes('[dead]');
            const dupe = containerText.includes('[dupe]');
            if (!isTopFlag && !(flagged || dead || dupe)) return;

            // Parse upvotes and comments
            let upvotes = 0;
            let comment_count = 0;
            let comment_link = 'No comments link';
            const subtextRow = row.nextElementSibling;
            if (subtextRow) {
              const scoreEl = subtextRow.querySelector('.score');
              const scoreMatch = scoreEl?.textContent?.match(/\d+/);
              upvotes = scoreMatch ? parseInt(scoreMatch[0], 10) : 0;

              const commentAnchor = Array.from(
                subtextRow.querySelectorAll('a')
              ).find((a) => a.textContent?.includes('comment')) as
                | HTMLAnchorElement
                | undefined;
              if (commentAnchor) {
                const countMatch = commentAnchor.textContent?.match(/\d+/);
                comment_count = countMatch ? parseInt(countMatch[0], 10) : 0;
                comment_link = commentAnchor.href;
              }
            }

            articles.push({
              title,
              link,
              hn_id,
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
          const moreLink = document.querySelector(
            'a.morelink'
          ) as HTMLAnchorElement | null;
          let nextUrl: string | null = null;
          if (moreLink) {
            if (isTopFlag) {
              nextUrl = moreLink.href;
            } else {
              const match = moreLink.href.match(
                /day=(\d{4}-\d{2}-\d{2})&p=(\d+)/
              );
              if (match) {
                const day = match[1];
                const pageNum = parseInt(match[2], 10);
                nextUrl = `${frontBase}?day=${day}&p=${pageNum}`;
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

  // Scrapes top stories from /
  const scrapeTopStories = async (maxPages?: number): Promise<Article[]> => {
    const articles: Article[] = [];
    let nextUrl: string | null = TOP_BASE_URL;
    let pageCount = 0;

    while (nextUrl) {
      console.log(`Scraping top stories: ${nextUrl}`);
      const { articles: pageArticles, nextUrl: newNext } =
        await scrapePageWithNext(nextUrl, true);
      articles.push(...pageArticles);
      nextUrl = newNext;
      pageCount++;
      if (maxPages && pageCount >= maxPages) break;
    }

    return articles;
  };

  // Scrapes flagged, dead, and duplicate stories from /front
  const scrapeFront = async (maxPages = 10): Promise<Article[]> => {
    const articles: Article[] = [];
    let nextUrl: string | null = FRONT_BASE_URL;
    let pageCount = 0;

    while (nextUrl && pageCount < maxPages) {
      console.log(`Scraping front page: ${nextUrl}`);
      const { articles: pageArticles, nextUrl: newNext } =
        await scrapePageWithNext(nextUrl, false);
      articles.push(...pageArticles);
      nextUrl = newNext;
      pageCount++;
    }

    return articles;
  };

  // Scrapes flagged stories for a specific day from /front
  const scrapeFrontForDay = async (day: string): Promise<Article[]> => {
    const articles: Article[] = [];
    let nextUrl: string | null = `${FRONT_BASE_URL}?day=${day}&p=1`;

    while (nextUrl) {
      console.log(`Scraping front for ${day}: ${nextUrl}`);
      const { articles: pageArticles, nextUrl: newNext } =
        await scrapePageWithNext(nextUrl, false);
      articles.push(...pageArticles);
      nextUrl = newNext;
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
