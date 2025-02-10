import { z } from 'zod';
import puppeteer from 'puppeteer-extra';
import StealthPlugin from 'puppeteer-extra-plugin-stealth';
import { Browser } from 'puppeteer';

// Creates and configures a Puppeteer browser client
export const createBrowserClient = async () => {
  puppeteer.use(StealthPlugin());
  const browser = await puppeteer.launch({
    headless: true,
    args: [
      '--no-sandbox',
      '--disable-setuid-sandbox',
      '--disable-dev-shm-usage',
      '--disable-cache',
      '--disk-cache-size=0',
      '--incognito',
      '--disable-gpu',
    ],
    protocolTimeout: 30000,
  });

  process.on('exit', async () => {
    await browser.close();
  });

  return browser;
};

export const BrowserClientSchema = z.instanceof(Browser);
export type BrowserClient = z.infer<typeof BrowserClientSchema>;
