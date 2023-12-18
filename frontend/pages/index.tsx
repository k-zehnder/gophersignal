import Layout from "../components/Layout";
import List from "@mui/joy/List";
import ListItem from "@mui/joy/ListItem";
import Typography from "@mui/joy/Typography";

export default function Index() {
  return (
    <Layout>
      <Typography level="h2" sx={{ mb: "1rem" }}>
        Minimalist Joy UI Blog
      </Typography>
      <Typography sx={{ mb: "1rem" }}>
        Welcome to your sleek new Joy UI blog. ✨
      </Typography>
      <Typography component="h2" level="h3">
        Features
      </Typography>
      <List>
        <ListItem>✓ Built with TypeScript</ListItem>
        <ListItem>✓ Designed with Joy UI's default styles</ListItem>
        <ListItem>✓ Ready to publish with Next.js Markdown blog</ListItem>
        <ListItem>✓ Light and dark modes with toggle button</ListItem>
        <ListItem>✓ Includes Prettier for code formatting</ListItem>
      </List>
      <Typography>
        View it on{" "}
        <a href="https://github.com/samuelsycamore/joy-next-blog/">GitHub</a>.
        Created with 💙 by <a href="https://mui.com/">MUI</a>.
      </Typography>
    </Layout>
  );
}
