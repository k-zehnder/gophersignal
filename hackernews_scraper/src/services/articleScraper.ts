// Scrapes a page and returns articles and the next URL. The IsTop flag distinguishes
// between the homepage (fresh stories) and /front (dup, dead, flagged articles).

import { Browser } from 'puppeteer';
import { Article, Scraper } from '../types';

const TOP_BASE_URL = 'https://news.ycombinator.com';
const FRONT_BASE_URL = 'https://news.ycombinator.com/front';

type ScrapeResult = { articles: Article[]; nextUrl: string | null };

const createHackerNewsScraper = (browser: Browser): Scraper => {
  // Scrapes a page and returns articles + next URL; IsTop flags / vs /front
  const scrapePageWithNext = async (
    pageUrl: string,
    isTop = false
  ): Promise<ScrapeResult> => {
    const page = await browser.newPage();
    try {
      await page.goto(pageUrl, { waitUntil: 'networkidle2' });

      // Run extraction in the page context
      const result: ScrapeResult = await page.evaluate(
        (isTopFlag: boolean, frontBase: string): ScrapeResult => {
          const articles: Article[] = [];

          // Select only rows with class "athing submission"
          const rows = Array.from(
            document.querySelectorAll<HTMLTableRowElement>(
              'tr.athing.submission'
            )
          );

          rows.forEach((row) => {
            // Get HN ID or fallback to comment link
            let hn_id = parseInt(row.id, 10) || 0;
            if (!hn_id) {
              const href =
                row.nextElementSibling?.querySelector<HTMLAnchorElement>(
                  'a[href*="item?id="]'
                )?.href;
              const match = href?.match(/item\?id=(\d+)/);
              hn_id = match ? parseInt(match[1], 10) : 0;
            }
            if (!hn_id) return;

            // Parse rank text
            const rankEl = row.querySelector<HTMLSpanElement>(
              'td.title > span.rank'
            );
            const rankText = rankEl?.innerText.replace(/\D/g, '').trim() ?? '0';
            const article_rank = parseInt(rankText, 10) || 0;

            // Parse title and link
            const titleEl = row.querySelector<HTMLAnchorElement>(
              'td.title > span.titleline a'
            );
            const title = titleEl?.innerText ?? 'No title';
            const link = titleEl?.href ?? 'No link';

            // Determine flags
            const containerText = titleEl?.parentElement?.innerText ?? '';
            const flagged = containerText.includes('[flagged]');
            const dead = containerText.includes('[dead]');
            const dupe = containerText.includes('[dupe]');
            if (!isTopFlag && !(flagged || dead || dupe)) return;

            // Parse upvotes and comments
            let upvotes = 0;
            let comment_count = 0;
            let comment_link = '';
            const subtextRow = row.nextElementSibling as Element | null;
            if (subtextRow) {
              const scoreEl =
                subtextRow.querySelector<HTMLSpanElement>('.score');
              const scoreMatch = scoreEl?.textContent?.match(/\d+/);
              upvotes = scoreMatch ? parseInt(scoreMatch[0], 10) : 0;

              const commentEl = Array.from(
                subtextRow.querySelectorAll<HTMLAnchorElement>('a')
              ).find((a) => a.textContent?.includes('comment'));
              if (commentEl) {
                const countMatch = commentEl.textContent?.match(/\d+/);
                comment_count = countMatch ? parseInt(countMatch[0], 10) : 0;
                comment_link = commentEl.href;
              }
            }

            articles.push({
              hn_id,
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

          // Find the "more" link to get the next URL
          const moreLink =
            document.querySelector<HTMLAnchorElement>('a.morelink');
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
    const all: Article[] = [];
    let nextUrl: string | null = TOP_BASE_URL;
    let pageCount = 0;

    while (nextUrl) {
      console.log(`Scraping top stories: ${nextUrl}`);
      const { articles, nextUrl: next } = await scrapePageWithNext(
        nextUrl,
        true
      );
      all.push(...articles);
      nextUrl = next;
      pageCount++;
      if (maxPages && pageCount >= maxPages) break;
    }

    return all;
  };

  // Scrapes flagged, dead, and duplicate stories from /front
  const scrapeFront = async (maxPages = 10): Promise<Article[]> => {
    const all: Article[] = [];
    let nextUrl: string | null = FRONT_BASE_URL;
    let pageCount = 0;

    while (nextUrl && pageCount < maxPages) {
      console.log(`Scraping front page: ${nextUrl}`);
      const { articles, nextUrl: next } = await scrapePageWithNext(
        nextUrl,
        false
      );
      all.push(...articles);
      nextUrl = next;
      pageCount++;
    }

    return all;
  };

  // Scrapes front feed for a specific day
  const scrapeFrontForDay = async (day: string): Promise<Article[]> => {
    const all: Article[] = [];
    let nextUrl: string | null = `${FRONT_BASE_URL}?day=${day}&p=1`;

    while (nextUrl) {
      console.log(`Scraping front for ${day}: ${nextUrl}`);
      const { articles, nextUrl: next } = await scrapePageWithNext(
        nextUrl,
        false
      );
      all.push(...articles);
      nextUrl = next;
    }

    return all;
  };

  return { scrapeTopStories, scrapeFront, scrapeFrontForDay };
};

export { createHackerNewsScraper };
