/**
 * NavBar: A component for rendering the navigation bar of the application.
 * It includes links to various sections like Home, About, and API. The NavBar
 * also contains the ModeButton component to toggle between light and dark modes.
 * The API URL is dynamically determined based on the environment.
 */
import Link from "next/link";
import List from "@mui/joy/List";
import ListItem from "@mui/joy/ListItem";
import Typography from "@mui/joy/Typography";
import ModeButton from "./ModeButton";
import { siteMetaData } from "../lib/siteMetaData";

// Define the API URL based on the environment.
const apiUrl =
  process.env.NEXT_PUBLIC_ENV === "development"
    ? "http://localhost:8080/swagger/index.html#/"
    : "https://gophersignal.com/swagger/index.html#/";

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
  {
    name: "API",
    path: apiUrl,
  },
];

// NavBar: A component for rendering the navigation bar of the application.
export default function NavBar() {
  return (
    <>
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
          {navLinks.map(({ path, name }) => {
            if (path.startsWith("http")) {
              {/* Render an external link for URLs starting with "http" */}
              return (
                <ListItem key={path}>
                  <a href={path} target="_blank" rel="noopener noreferrer">
                    {name}
                  </a>
                </ListItem>
              );
            } else {
              {/* Render a link to a local route for other paths */}
              return (
                <ListItem key={path}>
                  <Link href={path}>{name}</Link>
                </ListItem>
              );
            }
          })}
          <ListItem>
            {/* Render the ModeButton component for light/dark mode toggle */}
            <ModeButton />
          </ListItem>
        </List>
      </nav>
    </>
  );
}
