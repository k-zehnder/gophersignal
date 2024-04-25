import * as React from "react";
import Head from "next/head";
import { AppProps } from "next/app";
import { CssVarsProvider } from "@mui/joy/styles";
import GlobalStyles from "@mui/joy/GlobalStyles";
import theme from "../lib/theme";
import CssBaseline from "@mui/joy/CssBaseline";
import { CacheProvider, EmotionCache } from "@emotion/react";
import createEmotionCache from "../lib/createEmotionCache";

// Create a client-side cache for Emotion (CSS-in-JS) styles to be shared throughout the user's session.
const clientSideEmotionCache = createEmotionCache();

// Define the MyAppProps interface that extends AppProps and includes an optional emotionCache property.
interface MyAppProps extends AppProps {
  emotionCache?: EmotionCache;
}

// Define the main MyApp component.
export default function MyApp(props: MyAppProps) {
  // Destructure the Component, emotionCache, and pageProps from the props.
  const { Component, emotionCache = clientSideEmotionCache, pageProps } = props;
  return (
    <CacheProvider value={emotionCache}>
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
              overflowY: "scroll",
            },
            a: {
              // Define link styles for text decoration and color using CSS variables.
              textDecoration: "none",
              color: "var(--joy-palette-primary-500)",
            },
            "a:hover": {
              // Define link hover state color.
              color: "var(--joy-palette-primary-600)",
            },
            "a:active": {
              // Define link active state color.
              color: "var(--joy-palette-primary-700)",
            },
            li: {
              // Remove left padding for list items.
              paddingLeft: "0 !important",
            },
          }}
        />
        {/* Render the component passed through props with its pageProps. */}
        <Component {...pageProps} />
      </CssVarsProvider>
    </CacheProvider>
  );
}
