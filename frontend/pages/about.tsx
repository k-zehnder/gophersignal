import Layout from "../components/Layout";
import Typography from "@mui/joy/Typography";
import Avatar from "@mui/joy/Avatar";
import Link from "@mui/joy/Link";

import { siteMetaData } from "../lib/siteMetaData";

// Define the About component.
export default function About() {
  return (
    <Layout>
      {/* Display the heading "About". */}
      <Typography level="h2" sx={{ mb: "1rem", fontWeight: 'bold', fontSize: '1.75rem' }}>
        About
      </Typography>
      {/* Display the Gopher Signal logo as an Avatar. */}
      <Avatar
        sx={{
          "--Avatar-size": "100px",
          mb: "1rem",
        }}
        alt="Gopher Signal Logo"
        src={siteMetaData.image}
      />
      {/* Provide information about Gopher Signal. */}
      <Typography sx={{ fontSize: '1rem' }}>
        Gopher Signal uses smart technology to quickly summarize important points from{" "}
        {/* Create a link to Hacker News. */}
        <Link href="https://news.ycombinator.com" target="_blank" rel="noopener noreferrer">
          Hacker News
        </Link>{" "}
        articles, giving you brief and useful updates.
      </Typography>
    </Layout>
  );
}
