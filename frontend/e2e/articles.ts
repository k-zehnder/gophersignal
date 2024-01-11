import { chromium, Browser, Page } from 'playwright';

(async () => {
  const browser: Browser = await chromium.launch();
  const page: Page = await browser.newPage();

  // Navigate to the home page
  await page.goto('http://localhost:3000');

  // Check for the presence of an h2 element with the text "Latest Articles"
  const latestArticlesHeading = await page.$('h2:has-text("Latest Articles")');

  if (latestArticlesHeading) {
    console.log('Test passed: "Latest Articles" heading found.');
  } else {
    console.error('Test failed: "Latest Articles" heading not found.');
  }

  await browser.close();
})();
