import puppeteer from 'puppeteer-extra';
import StealthPlugin from 'puppeteer-extra-plugin-stealth';

// Creates and configures a Puppeteer browser client
export const createBrowserClient = async () => {
  puppeteer.use(StealthPlugin());
  const browser = await puppeteer.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox'],
    protocolTimeout: 30000,
  });
  return browser;
};

export type BrowserClient = Awaited<ReturnType<typeof createBrowserClient>>;
