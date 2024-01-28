// Import necessary dependencies.
import * as React from "react";
import Document, { Html, Head, Main, NextScript } from "next/document";
import createEmotionServer from "@emotion/server/create-instance";
import createEmotionCache from "../lib/createEmotionCache";
import { getInitColorSchemeScript } from "@mui/joy/styles";

// Define the MyDocument class which extends the Document class.
export default class MyDocument extends Document {
  render() {
    return (
      <Html lang="en">
        <Head>
          {/* Add a shortcut icon for the website. */}
          <link rel="shortcut icon" href="/favicon.ico" />

          {/* Google Analytics Script */}
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

          {/* Inject Emotion styles */}
          {(this.props as any).emotionStyleTags}
        </Head>
        <body>
          {/* Initialize the color scheme script with defaultMode as "system". */}
          {getInitColorSchemeScript({ defaultMode: "system" })}

          {/* Render the main content of the application. */}
          <Main />

          {/* Render Next.js scripts. */}
          <NextScript />
        </body>
      </Html>
    );
  }
}

// Define getInitialProps function for server-side rendering.
MyDocument.getInitialProps = async (ctx) => {
  const originalRenderPage = ctx.renderPage;
  const cache = createEmotionCache();
  const { extractCriticalToChunks } = createEmotionServer(cache);

  ctx.renderPage = () =>
    originalRenderPage({
      enhanceApp: (App: any) =>
        function EnhanceApp(props) {
          return <App emotionCache={cache} {...props} />;
        },
    });

  const initialProps = await Document.getInitialProps(ctx);
  const emotionStyles = extractCriticalToChunks(initialProps.html);
  const emotionStyleTags = emotionStyles.styles.map((style) => (
    <style
      data-emotion={`${style.key} ${style.ids.join(" ")}`}
      key={style.key}
      // eslint-disable-next-line react/no-danger
      dangerouslySetInnerHTML={{ __html: style.css }}
    />
  ));

  return {
    ...initialProps,
    emotionStyleTags,
  };
};
