import Link from "next/link";
import List from "@mui/joy/List";
import ListItem from "@mui/joy/ListItem";
import Typography from "@mui/joy/Typography";
import ModeButton from "./ModeButton";
import { siteMetaData } from "../lib/siteMetaData";

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
    path: "https://gophersignal.com/swagger/index.html#/",
  },
];

export default function NavBar() {
  return (
    <>
      <Link href="/">
        <Typography component="h1" level="display2" fontSize="xl">
          {siteMetaData.title}
        </Typography>
      </Link>
      <nav>
        <List
          sx={{ display: "flex", flexDirection: { xs: "row", md: "column" } }}
        >
          {navLinks.map(({ path, name }) => {
            if (path.startsWith("http")) {
              return (
                <ListItem key={path}>
                  <a href={path} target="_blank" rel="noopener noreferrer">
                    {name}
                  </a>
                </ListItem>
              );
            } else {
              return (
                <ListItem key={path}>
                  <Link href={path}>{name}</Link>
                </ListItem>
              );
            }
          })}
          <ListItem>
            <ModeButton />
          </ListItem>
        </List>
      </nav>
    </>
  );
}
