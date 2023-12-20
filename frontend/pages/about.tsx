import Layout from "../components/Layout";
import Typography from "@mui/joy/Typography";
import Avatar from "@mui/joy/Avatar";
import { siteMetaData } from "../lib/siteMetaData";

export default function About() {
  return (
    <Layout>
      <Typography level="h2" sx={{ mb: "1rem" }}>
        About
      </Typography>
      <Avatar
        sx={{
          "--Avatar-size": "100px",
          mb: "1rem",
        }}
        alt="this image does not exist"
        src={siteMetaData.image}
      >
        JN
      </Avatar>
      <Typography>
        This application turns HackerNews articles into short code snippets.
      </Typography>
    </Layout>
  );
}
