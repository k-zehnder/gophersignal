import * as React from 'react';
import Head from 'next/head';
import { CssVarsProvider } from '@mui/joy/styles';
import GlobalStyles from '@mui/joy/GlobalStyles';
import theme from '../lib/theme';
import CssBaseline from '@mui/joy/CssBaseline';

const MyApp: React.FC<{ Component: React.ElementType; pageProps: any }> = ({
  Component,
  pageProps,
}) => {
  return (
    <>
      <Head>
        {/* Set viewport meta tag for responsive design */}
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </Head>
      <CssVarsProvider
        defaultMode="system"
        theme={theme}
        modeStorageKey="mode-key"
        disableNestedContext
      >
        <CssBaseline />
        <GlobalStyles
          styles={{
            html: {
              // Ensure vertical scrolling is always available for content overflow.
              overflowY: 'scroll',
            },
            a: {
              textDecoration: 'none',
              color: 'var(--joy-palette-primary-500)',
            },
            'a:hover': {
              color: 'var(--joy-palette-primary-600)',
            },
            'a:active': {
              color: 'var(--joy-palette-primary-700)',
            },
            li: {
              // Remove left padding for list items.
              paddingLeft: '0 !important',
            },
          }}
        />
        {/* Render the component passed through props with its pageProps. */}
        <Component {...pageProps} />
      </CssVarsProvider>
    </>
  );
};

export default MyApp;
