import React from 'react';
import Link from 'next/link';
import List from '@mui/joy/List';
import ListItem from '@mui/joy/ListItem';
import Typography from '@mui/joy/Typography';
import ModeButton from './ModeButton';
import { siteMetaData } from '../lib/siteMetaData';

// Define the API URL based on the environment.
const apiUrl =
  process.env.NEXT_PUBLIC_ENV === 'development'
    ? 'http://localhost:8080/swagger/index.html#/'
    : 'https://gophersignal.com/swagger/index.html#/';

// Define navigation links for the NavBar.
const navLinks = [
  {
    name: 'Home',
    path: '/',
  },
  {
    name: 'About',
    path: '/about',
  },
  {
    name: 'API',
    path: apiUrl,
  },
];

// NavBar component for rendering the navigation bar.
export default function NavBar() {
  return (
    <>
      <Link href="/">
        {/* Render the site title as a link to the home page */}
        {/* @ts-ignore */}
        <Typography component="h1" level="display2" fontSize="xl">
          {siteMetaData.title}
        </Typography>
      </Link>
      <nav>
        <List
          sx={{
            display: 'flex',
            flexDirection: { xs: 'row', md: 'column' },
          }}
        >
          {navLinks.map(({ path, name }) => {
            // Check if the path is an external URL or a local route.
            if (path.startsWith('http')) {
              return (
                // Render an external link for URLs starting with "http".
                <ListItem key={path}>
                  <a href={path} target="_blank" rel="noopener noreferrer">
                    {name}
                  </a>
                </ListItem>
              );
            } else {
              return (
                // Render a link to a local route for other paths.
                <ListItem key={path}>
                  <Link href={path}>{name}</Link>
                </ListItem>
              );
            }
          })}
          <ListItem>
            {/* Render the ModeButton component for light/dark mode toggle. */}
            <ModeButton />
          </ListItem>
        </List>
      </nav>
    </>
  );
}
