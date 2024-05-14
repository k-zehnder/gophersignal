// Fetches the main content of an article from a given URL using Puppeteer.

const {
  handlePopups,
  acceptCookieConsent,
} = require('../helpers/scraperUtils');

// Fetches the main content of an article from a given URL using Puppeteer.
const fetchArticleContent = async (url, browser) => {
  const navigationTimeout = 30000;

  try {
    const page = await browser.newPage();
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

    // Handle cookie consent if present
    const acceptCookie = await acceptCookieConsent(page);
    if (acceptCookie) {
      console.log('Cookie consent was accepted successfully.');
    } else {
      console.log('No cookie consent found.');
    }

    // Handle popups if present
    const popupClosed = await handlePopups(page);
    if (popupClosed) {
      console.log('A popup was successfully closed.');
    } else {
      console.log('No popup found.');
    }

    // Extract main content
    let content = await page.evaluate(() => {
      const mainContentSelector = 'main, article, .post, .text, [role="main"]';
      let mainContent = document.querySelector(mainContentSelector);

      if (!mainContent) {
        mainContent = document.querySelector('body');
        const nonContentSelectors =
          'nav, aside, footer, header, .sidebar, .menu, .footer';
        document
          .querySelectorAll(nonContentSelectors)
          .forEach((el) => el.remove());
      }

      const paragraphs = mainContent.querySelectorAll('p');
      return Array.from(paragraphs)
        .map((p) => p.innerText.trim())
        .join('\n\n');
    });

    await page.close();
    console.log(`Fetched content for: ${url}`);
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

module.exports = { fetchArticleContent };
