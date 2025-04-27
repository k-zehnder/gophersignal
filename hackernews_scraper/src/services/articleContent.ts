// Fetches full article content from Hacker News links, handling cookie consents, popups, and arXiv abstracts.

import { Browser, Page } from 'puppeteer';

const createContentFetcher = (browser: Browser) => {
  // Detect arXiv abstract pages by URL
  const isArxivUrl = (url: string): boolean =>
    /^https?:\/\/(?:www\.)?arxiv\.org\/abs\//.test(url);

  // Handles cookie consent on web pages by clicking on common cookie consent buttons
  const acceptCookieConsent = async (page: Page): Promise<boolean> => {
    const commonSelectors = [
      'cookie-banner',
      '.cookie-consent',
      '#cookie-accept',
    ];
    for (const selector of commonSelectors) {
      const elements = await page.$$(selector);
      if (elements.length > 0) {
        await elements[0].click();
        console.log('Clicked on cookie consent button');
        return true;
      }
    }
    return false;
  };

  // Handles popups and modal dialogs on web pages
  const handlePopups = async (page: Page): Promise<boolean> => {
    const popupSelectors = ['.popup', '.overlay', '.modal', '.modal-dialog'];
    for (const selector of popupSelectors) {
      const elements = await page.$$(selector);
      if (elements.length > 0) {
        await elements[0].click();
        console.log('Closed popup or overlay');
        return true;
      }
    }

    page.on('dialog', async (dialog: any) => {
      console.log('Dialog message:', dialog.message());
      await dialog.accept();
    });

    // Give dialogs a moment to appear
    await new Promise((resolve) => setTimeout(resolve, 1000));
    return false;
  };

  // Try parsing an arXiv abstract page
  const tryParseArxiv = async (
    page: Page,
    url: string
  ): Promise<string | null> => {
    if (!isArxivUrl(url)) {
      return null;
    }

    const { title, authors, abstract } = await page.evaluate(() => {
      const t =
        document
          .querySelector('h1.title')
          ?.textContent?.replace('Title:', '')
          .trim() || '';
      const a = Array.from(document.querySelectorAll('.authors a'))
        .map((el) => el.textContent?.trim() || '')
        .filter(Boolean);
      const abs =
        document
          .querySelector('blockquote.abstract')
          ?.textContent?.replace('Abstract:', '')
          .trim() || '';
      return { title: t, authors: a, abstract: abs };
    });

    console.info(`Parsed arXiv page: ${url}`);
    return [
      `Title: ${title}`,
      `Authors: ${authors.join(', ')}`,
      '',
      'Abstract:',
      abstract,
    ].join('\n');
  };

  // Fetches the main content of an article from a given URL using Puppeteer
  const fetchArticleContent = async (url: string): Promise<string> => {
    const navigationTimeout = 30_000;
    const page = await browser.newPage();

    try {
      // Block images, stylesheets, and fonts for faster loading
      await page.setRequestInterception(true);
      page.on('request', (req) => {
        if (['image', 'stylesheet', 'font'].includes(req.resourceType())) {
          req.abort();
        } else {
          req.continue();
        }
      });
      page.setDefaultNavigationTimeout(navigationTimeout);

      await page.goto(url, {
        waitUntil: 'networkidle2',
        timeout: navigationTimeout,
      });

      // arXiv-specific parsing
      const arxivResult = await tryParseArxiv(page, url);
      if (arxivResult !== null) {
        await page.close();
        return arxivResult;
      }

      // Generic flow for other domains
      const accepted = await acceptCookieConsent(page);
      console.info(
        accepted ? 'Cookie consent accepted' : 'No cookie consent found'
      );

      const closedPopup = await handlePopups(page);
      console.info(closedPopup ? 'Popup closed' : 'No popup found');

      const content = await page.evaluate(() => {
        const mainSelectors = 'main, article, .post, .text, [role="main"]';
        let root =
          document.querySelector(mainSelectors) ||
          (document.body as HTMLElement);

        // Remove navigation, sidebars, footers, headers
        document
          .querySelectorAll('nav, aside, footer, header, .sidebar, .menu')
          .forEach((el) => el.remove());

        const paragraphs = root.querySelectorAll('p');
        return Array.from(paragraphs)
          .map((p) => p.innerText.trim())
          .join('\n\n');
      });

      await page.close();
      console.info(`Fetched content for: ${url}`);
      return content;
    } catch (error) {
      console.error(
        `Failed to fetch content for ${url} within ${
          navigationTimeout / 1000
        } seconds:`,
        error
      );
      await page.close();
      return '';
    }
  };

  return { fetchArticleContent };
};

export { createContentFetcher };
