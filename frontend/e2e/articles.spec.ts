import { test, expect } from '@playwright/test';

test('latest articles heading', async ({ page }) => {
  // Navigate to the home page
  await page.goto('http://localhost:3000');

  // Check for the presence of an h2 element with the text "Latest Articles"
  const latestArticlesHeading = await page.$('h2:has-text("Latest Articles")');

  // Expect that the "Latest Articles" heading exists in the DOM
  await expect(latestArticlesHeading).toBeTruthy();
});
