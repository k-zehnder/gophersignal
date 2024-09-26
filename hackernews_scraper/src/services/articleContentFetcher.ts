// Fetches the full content of articles linked from top Hacker News stories,
// including handling cookie consents and popups on web pages.

import { Browser, Page } from 'puppeteer';

const createContentFetcher = (browser: Browser) => {
  /**
   * Handles cookie consent on web pages by clicking on common cookie consent buttons.
   */
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

  /**
   * Handles popups and modal dialogs on web pages.
   */
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

    await new Promise((resolve) => setTimeout(resolve, 1000));
    return false;
  };

  /**
   * Fetches the main content of an article from a given URL using Puppeteer.
   */
  const fetchArticleContent = async (url: string): Promise<string> => {
    const navigationTimeout = 30000;

    try {
      const page = await browser.newPage();
      await page.setRequestInterception(true);
      page.on('request', (req: any) => {
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

      const acceptCookie = await acceptCookieConsent(page);
      if (acceptCookie) {
        console.info('Cookie consent was accepted successfully.');
      } else {
        console.info('No cookie consent found.');
      }

      const popupClosed = await handlePopups(page);
      if (popupClosed) {
        console.info('A popup was successfully closed.');
      } else {
        console.info('No popup found.');
      }

      const content = await page.evaluate(() => {
        const mainContentSelector =
          'main, article, .post, .text, [role="main"]';
        let mainContent = document.querySelector(mainContentSelector);

        if (!mainContent) {
          mainContent = document.querySelector('body');
          const nonContentSelectors =
            'nav, aside, footer, header, .sidebar, .menu, .footer';
          document
            .querySelectorAll(nonContentSelectors)
            .forEach((el) => el.remove());
        }

        const paragraphs = mainContent?.querySelectorAll('p');
        return Array.from(paragraphs!)
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
      return '';
    }
  };

  return { fetchArticleContent };
};

export { createContentFetcher };
