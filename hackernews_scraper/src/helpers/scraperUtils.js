// Provides utility functions for handling cookie consent, popups, delays, and logging.

// Handles cookie consent on web pages by clicking on common cookie consent buttons.
const acceptCookieConsent = async (page) => {
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
      return true; // Indicate that the cookie consent button was clicked
    }
  }
  return false; // Indicate that no cookie consent button was found
};

// Handles popups and modal dialogs on web pages.
const handlePopups = async (page) => {
  const popupSelectors = ['.popup', '.overlay', '.modal', '.modal-dialog'];
  for (const selector of popupSelectors) {
    const elements = await page.$$(selector);
    if (elements.length > 0) {
      await elements[0].click(); // Click on the first element found
      console.log('Closed popup or overlay');
      return true; // Indicate that a popup was closed
    }
  }

  // Handles modal dialogs (e.g., alert, prompt, confirmation).
  page.on('dialog', async (dialog) => {
    console.log('Dialog message:', dialog.message());
    await dialog.accept();
  });

  // Simple delay to allow for potential page changes.
  await delay(1000);

  return false; // Indicate that no popup was closed
};

// Provides a delay for a specified number of milliseconds.
const delay = async (ms) => new Promise((resolve) => setTimeout(resolve, ms));

// Provides logging functionality with different levels of severity.
const logger = {
  // Logs informational messages.
  info: (message, ...args) => {
    console.log(`INFO: ${message}`, ...args);
  },
  // Logs warning messages.
  warn: (message, ...args) => {
    console.warn(`WARN: ${message}`, ...args);
  },
  // Logs error messages.
  error: (message, ...args) => {
    console.error(`ERROR: ${message}`, ...args);
  },
};

module.exports = {
  acceptCookieConsent,
  handlePopups,
  delay,
  logger,
};
