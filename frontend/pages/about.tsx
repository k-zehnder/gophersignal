import React from 'react';
import Layout from '../components/Layout';
import Typography from '@mui/joy/Typography';
import Avatar from '@mui/joy/Avatar';
import { siteMetaData } from '../lib/siteMetaData';

export default function About() {
  return (
    <Layout>
      {/* Heading */}
      <Typography
        level="h2"
        sx={{ mb: '1rem', fontWeight: 'bold', fontSize: '1.75rem' }}
      >
        About
      </Typography>

      {/* Logo */}
      <Avatar
        sx={{
          '--Avatar-size': '100px',
          mb: '1rem',
        }}
        alt="Gopher Signal Logo"
        src={siteMetaData.image}
      />

      {/* Description */}
      <Typography sx={{ fontSize: '1rem', mb: '0.5rem' }}>
        Gopher Signal uses smart technology to summarize key points from{' '}
        <a
          href="https://news.ycombinator.com"
          target="_blank"
          rel="noopener noreferrer"
          style={{
            color: '#007bff',
            textDecoration: 'none',
          }}
        >
          Hacker News
        </a>{' '}
        articles, providing concise and useful updates.
      </Typography>

      {/* GitHub Link */}
      <Typography sx={{ fontSize: '1rem', mt: '1rem' }}>
        Check out the project on{' '}
        <a
          href="https://github.com/k-zehnder/gophersignal"
          target="_blank"
          rel="noopener noreferrer"
          style={{
            color: '#007bff',
            textDecoration: 'none',
          }}
        >
          GitHub
        </a>{' '}
        — we're seeking contributors!
      </Typography>
    </Layout>
  );
}
