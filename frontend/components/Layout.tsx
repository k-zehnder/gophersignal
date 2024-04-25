import Box from "@mui/joy/Box";
import Grid from "@mui/joy/Grid";
import NavBar from "./NavBar";
import { PropsWithChildren } from "react";
import { siteMetaData } from "../lib/siteMetaData";
import Head from "next/head";
import { useRouter } from "next/router";

// Layout component defines the overall structure of the application.
export default function Layout(props: PropsWithChildren) {
  const router = useRouter();
  
  return (
    <>
      {/* Define HTML head content, including metadata and title */}
      <Head>
        <title>{siteMetaData.title}</title>
        {/* Define metadata for web crawlers */}
        <meta name="robots" content="follow, index" />
        <meta name="description" content={siteMetaData.description} />
        {/* Open Graph (OG) metadata for sharing on social media */}
        <meta property="og:url" content={`${siteMetaData.siteUrl}${router.asPath}`} />
        <meta property="og:type" content="blog" />
        <meta property="og:site_name" content={siteMetaData.title} />
        <meta property="og:description" content={siteMetaData.description} />
        <meta property="og:title" content={siteMetaData.title} />
        <meta property="og:image" content={`${siteMetaData.siteUrl}${siteMetaData.ogImage}`} key={siteMetaData.ogImage} />
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:title" content={siteMetaData.title} />
        <meta name="twitter:description" content={siteMetaData.description} />
        <meta name="twitter:image" content={`${siteMetaData.siteUrl}${siteMetaData.ogImage}`} />
        {/* Define canonical link */}
        <link rel="canonical" href={`${siteMetaData.siteUrl}${router.asPath}`} />
      </Head>
      {/* Create a responsive grid layout */}
      <Grid
        container
        spacing={2}
        sx={{
          maxWidth: "1000px",
          display: "flex",
          flexGrow: 1,
          flexDirection: { xs: "column", md: "row" },
          mx: "auto",
          mt: { xs: "2rem", md: "6rem" },
          p: "2rem",
        }}
      >
        {/* Left column: Navigation bar */}
        <Grid xs={12} md={4}>
          <NavBar />
        </Grid>
        {/* Right column: Content */}
        <Grid xs={12} md={8}>
          <Box>{props.children}</Box>
        </Grid>
      </Grid>
    </>
  );
}
