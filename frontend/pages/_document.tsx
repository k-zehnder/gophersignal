// Import necessary dependencies for customizing the default document structure in Next.js
import * as React from "react";
import Document, { Html, Head, Main, NextScript } from "next/document";
import createEmotionServer from "@emotion/server/create-instance";
import createEmotionCache from "../lib/createEmotionCache";
import { getInitColorSchemeScript } from "@mui/joy/styles";

// Define the MyDocument class which extends the Document class.
// This custom document configuration is used to augment our application's <html> and <body> tags.
export default class MyDocument extends Document {
  render() {
    return (
      <Html lang="en"> {/* Setting the language helps with accessibility */}
        <Head>
          {/* Favicon hosted on S3, optimized and cached via Cloudflare for better performance */}
          <link rel="shortcut icon" href="https://gophersignal-cloudflare-assets.s3.us-west-1.amazonaws.com/favicon.ico" />

          {/* Google Analytics Script for tracking website traffic and user interactions */}
          <script async src="https://www.googletagmanager.com/gtag/js?id=G-H03QDKFRJ0"></script>
          <script
            dangerouslySetInnerHTML={{
              __html: `
                window.dataLayer = window.dataLayer || [];
                function gtag(){dataLayer.push(arguments);}
                gtag('js', new Date());
                gtag('config', 'G-H03QDKFRJ0');
              `,
            }}
          />

          {/* Inject Emotion styles - CSS-in-JS solution for styling components */}
          {(this.props as any).emotionStyleTags}
        </Head>
        <body>
          {/* Initializes color scheme based on the system preference, enhancing user experience */}
          {getInitColorSchemeScript({ defaultMode: "system" })}

          {/* Main application content will be rendered here */}
          <Main />

          {/* Next.js scripts required for the application's functionality */}
          <NextScript />
        </body>
      </Html>
    );
  }
}

// getInitialProps allows for server-side rendering with initial data population.
// It augments the render process with additional steps like creating a cache for Emotion.
MyDocument.getInitialProps = async (ctx) => {
  const originalRenderPage = ctx.renderPage;
  const cache = createEmotionCache();
  const { extractCriticalToChunks } = createEmotionServer(cache);

  // Enhance the app with server-side rendered styles
  ctx.renderPage = () =>
    originalRenderPage({
      enhanceApp: (App: any) =>
        function EnhanceApp(props) {
          return <App emotionCache={cache} {...props} />;
        },
    });

  // Extract the critical CSS from the server-side rendering process
  const initialProps = await Document.getInitialProps(ctx);
  const emotionStyles = extractCriticalToChunks(initialProps.html);
  const emotionStyleTags = emotionStyles.styles.map((style) => (
    <style
      data-emotion={`${style.key} ${style.ids.join(" ")}`}
      key={style.key}
      dangerouslySetInnerHTML={{ __html: style.css }}
    />
  ));

  // Return the initial props along with the extracted styles
  return {
    ...initialProps,
    emotionStyleTags,
  };
};
