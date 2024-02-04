/**
 * createEmotionCache: A utility to customize the insertion point for Emotion's style tag.
 * This configuration ensures Material UI styles are loaded first in the <head> section,
 * enabling developers to override these styles using other CSS solutions.
 * The cache is configured for client-side rendering, where the browser environment is available.
*/
import createCache from "@emotion/cache";

const isBrowser = typeof document !== "undefined";

// On the client side, Create a meta tag at the top of the <head> and set it as insertionPoint.
// This assures that MUI styles are loaded first.
// It allows developers to easily override MUI styles with other styling solutions, like CSS modules.
export default function createEmotionCache() {
  let insertionPoint;

  if (isBrowser) {
    const emotionInsertionPoint = document.querySelector<HTMLMetaElement>(
      'meta[name="emotion-insertion-point"]'
    );
    insertionPoint = emotionInsertionPoint ?? undefined;
  }

  return createCache({ key: "mui-style", insertionPoint });
}
