import React from 'react';
import Document, {
  Html,
  Head,
  Main,
  NextScript,
  DocumentContext,
} from 'next/document';
import { ReactElement } from 'react';

// Extends Next.js's default Document to customize the HTML document structure, allowing inclusion of global elements like favicon for branding and Google Analytics for traffic analysis.
export default class MyDocument extends Document {
  static async getInitialProps(ctx: DocumentContext): Promise<any> {
    const initialProps = await Document.getInitialProps(ctx);
    return { ...initialProps };
  }

  render(): ReactElement {
    return (
      <Html lang="en">
        <Head>
          {/* Add a shortcut icon for the website. */}
          <link rel="shortcut icon" href="/favicon.ico" />

          {/* Google Analytics Script */}
          <script
            async
            src="https://www.googletagmanager.com/gtag/js?id=G-H03QDKFRJ0"
          ></script>
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
        </Head>
        <body>
          {/* Render the main content of the application. */}
          <Main />

          {/* Render Next.js scripts. */}
          <NextScript />
        </body>
      </Html>
    );
  }
}
