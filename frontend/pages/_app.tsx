import * as React from "react";
import Head from "next/head";
import { AppProps } from "next/app";
import { CssVarsProvider } from "@mui/joy/styles";
import GlobalStyles from "@mui/joy/GlobalStyles";
import theme from "../lib/theme";
import CssBaseline from "@mui/joy/CssBaseline";
import { CacheProvider, EmotionCache } from "@emotion/react";
import createEmotionCache from "../lib/createEmotionCache";

// Client-side cache, shared for the whole session of the user in the browser.
const clientSideEmotionCache = createEmotionCache();

interface MyAppProps extends AppProps {
  emotionCache?: EmotionCache;
}

export default function MyApp(props: MyAppProps) {
  const { Component, emotionCache = clientSideEmotionCache, pageProps } = props;
  return (
    <CacheProvider value={emotionCache}>
      <Head>
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
              overflowY: "scroll",
            },
            a: {
              textDecoration: "none",
              color: "var(--joy-palette-primary-500)",
            },
            "a:hover": {
              color: "var(--joy-palette-primary-600)",
            },
            "a:active": {
              color: "var(--joy-palette-primary-700)",
            },
            li: {
              paddingLeft: "0 !important",
            },
          }}
        />
        <Component {...pageProps} />
      </CssVarsProvider>
    </CacheProvider>
  );
}
