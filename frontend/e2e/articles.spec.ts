// Import Playwright's test and expect functions for testing.
import { test, expect } from '@playwright/test';

// Define a test case with a description.
test('has title', async ({ page }) => {
  // Navigate to the specified URL.
  await page.goto('http://localhost:3000/');

  // Use the expect function to check if the page has the expected title using a regular expression.
  await expect(page).toHaveTitle(/Gopher Signal/);
});
