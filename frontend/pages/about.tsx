import Layout from "../components/Layout";
import Typography from "@mui/joy/Typography";
import Avatar from "@mui/joy/Avatar";
import { siteMetaData } from "../lib/siteMetaData";

export default function About() {
  return (
    <Layout>
      <Typography level="h2" sx={{ mb: "1rem", fontWeight: 'bold', fontSize: '1.75rem' }}>
        About
      </Typography>
      <Avatar
        sx={{
          "--Avatar-size": "100px",
          mb: "1rem",
        }}
        alt="Opher Signal Logo"
        src={siteMetaData.image}
      />
      <Typography sx={{ fontSize: '1rem' }}>
        Gopher Signal uses smart technology to quickly summarize important points from Hacker News articles, giving you brief and useful updates.
      </Typography>
    </Layout>
  );
}
