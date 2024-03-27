import Link from "next/link";
import List from "@mui/joy/List";
import ListItem from "@mui/joy/ListItem";
import Typography from "@mui/joy/Typography";
import ModeButton from "./ModeButton";
import { siteMetaData } from "../lib/siteMetaData";

// Define navigation links for the NavBar.
const navLinks = [
  {
    name: "Home",
    path: "/",
  },
  {
    name: "About",
    path: "/about",
  },
];

// NavBar component for rendering the navigation bar.
export default function NavBar() {
  return (
    <>
      {/* Render the red alert bar to notify users that the API is down */}
      <div style={{ backgroundColor: "red", color: "white", padding: "10px", textAlign: "center" }}>
        The HuggingFace API is currently down and summaries are unavailable. We apologize for the inconvenience.
      </div>
      <Link href="/">
        {/* Render the site title as a link to the home page */}
        <Typography component="h1" level="display2" fontSize="xl">
          {siteMetaData.title}
        </Typography>
      </Link>
      <nav>
        <List
          sx={{
            display: "flex",
            flexDirection: { xs: "row", md: "column" },
          }}
        >
          {navLinks.map(({ path, name }) => (
            // Render a link to a local route for each navigation item
            <ListItem key={path}>
              <Link href={path}>{name}</Link>
            </ListItem>
          ))}
          <ListItem>
            {/* Render the ModeButton component for light/dark mode toggle */}
            <ModeButton />
          </ListItem>
        </List>
      </nav>
    </>
  );
}
