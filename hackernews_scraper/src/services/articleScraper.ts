// Scrapes a page and returns the articles and next URL. The isTop flag distinguishes
// between the homepage (fresh stories) and /front (dup, dead, flagged articles).

import { Browser } from 'puppeteer';
import { ArticleSchema, type Article } from '../types';
import { Scraper } from '../types';

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
          const articles: unknown[] = [];
          const rows = Array.from(document.querySelectorAll('tr.athing'));

          rows.forEach((row: Element) => {
            // Rank parsing
            const rankElement = row.querySelector('td.title > span.rank');
            const rankStr =
              rankElement?.textContent?.replace(/\D/g, '').trim() || '0';
            const article_rank = parseInt(rankStr, 10) || 0;

            // Title and link parsing
            const titleElement = row.querySelector<HTMLAnchorElement>(
              'td.title > span.titleline a'
            );
            const title = titleElement?.textContent?.trim() || 'No title';
            const link = titleElement?.href || '';
            const titleText = titleElement?.parentElement?.textContent || '';

            // Flags detection
            const flagged = titleText.includes('[flagged]');
            const dead = titleText.includes('[dead]');
            const dupe = titleText.includes('[dupe]');

            if (!isTop && !(flagged || dead || dupe)) return;

            // Engagement metrics
            let upvotes = 0,
              comment_count = 0;
            let comment_link = '';
            const subtextRow = row.nextElementSibling;

            if (subtextRow) {
              const subtext = subtextRow.querySelector('.subtext');
              if (subtext) {
                // Score parsing
                const scoreEl = subtext.querySelector('.score');
                const scoreText = scoreEl?.textContent?.replace(/,/g, '') || '';
                upvotes = parseInt(scoreText) || 0;

                // Comment parsing (last link in subtext)
                const links = Array.from(
                  subtext.querySelectorAll<HTMLAnchorElement>('a')
                );
                const commentLink = links[links.length - 1];
                if (commentLink) {
                  comment_link = commentLink.href;
                  const commentText =
                    commentLink.textContent?.replace(/,/g, '') || '';
                  comment_count = parseInt(
                    commentText.match(/\d+/)?.[0] || '0',
                    10
                  );
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

          // Next URL resolution
          const moreLink =
            document.querySelector<HTMLAnchorElement>('a.morelink');
          let nextUrl: string | null = null;

          if (moreLink) {
            if (isTop) {
              nextUrl = moreLink.href;
            } else {
              const match = moreLink.href.match(
                /day=(\d{4}-\d{2}-\d{2})&p=(\d+)/
              );
              nextUrl = match
                ? `${frontBase}?day=${match[1]}&p=${match[2]}`
                : null;
            }
          }

          return { articles, nextUrl };
        },
        isTop,
        FRONT_BASE_URL
      );

      if (!result) return { articles: [], nextUrl: null };

      // Parse articles with Zod
      const parsedArticles = result.articles.map((rawArticle) => {
        const article = rawArticle as {
          title: string;
          link: string;
          article_rank: number;
          flagged: boolean;
          dead: boolean;
          dupe: boolean;
          upvotes: number;
          comment_count: number;
          comment_link: string;
        };

        // Create proper URL objects
        const makeAbsoluteUrl = (path: string) =>
          new URL(path, TOP_BASE_URL).toString();

        return ArticleSchema.parse({
          ...article,
          link: article.link.startsWith('http')
            ? article.link
            : makeAbsoluteUrl(article.link),
          comment_link: article.comment_link.startsWith('http')
            ? article.comment_link
            : makeAbsoluteUrl(article.comment_link),
        });
      });

      return {
        articles: parsedArticles,
        nextUrl: result.nextUrl
          ? new URL(result.nextUrl, TOP_BASE_URL).toString()
          : null,
      };
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
      const { articles: pageArticles, nextUrl: newNextUrl } =
        await scrapePageWithNext(nextUrl, false);
      articles = articles.concat(pageArticles);
      nextUrl = newNextUrl;
      if (++pageCount >= maxPages) break;
    }

    return articles;
  };

  // Scrapes front pages for a specific day from (/front)
  const scrapeFrontForDay = async (day: string): Promise<Article[]> => {
    let articles: Article[] = [];
    let nextUrl: string | null = `${FRONT_BASE_URL}?day=${day}&p=1`;

    while (nextUrl) {
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
