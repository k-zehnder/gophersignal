import React from 'react';
import Layout from '../components/Layout';
import Typography from '@mui/joy/Typography';
import Avatar from '@mui/joy/Avatar';
import { siteMetaData } from '../lib/siteMetaData';

export default function About() {
  return (
    <Layout>
      {/* Display the heading "About". */}
      <Typography
        level="h2"
        sx={{ mb: '1rem', fontWeight: 'bold', fontSize: '1.75rem' }}
      >
        About
      </Typography>

      {/* Display the Gopher Signal logo as an Avatar. */}
      <Avatar
        sx={{
          '--Avatar-size': '100px',
          mb: '1rem',
        }}
        alt="Gopher Signal Logo"
        src={siteMetaData.image}
      />

      {/* Provide information about Gopher Signal. */}
      <Typography sx={{ fontSize: '1rem', mb: '0.5rem' }}>
        Gopher Signal uses smart technology to quickly summarize important
        points from{' '}
        <a
          href="https://news.ycombinator.com"
          target="_blank"
          rel="noopener noreferrer"
          style={{
            color: '#007bff',
            textDecoration: 'none',
          }}
          onMouseEnter={(e) =>
            (e.currentTarget.style.textDecoration = 'underline')
          }
          onMouseLeave={(e) => (e.currentTarget.style.textDecoration = 'none')}
        >
          Hacker News
        </a>{' '}
        articles, giving you brief and useful updates.
      </Typography>

      {/* Provide GitHub link. */}
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
          onMouseEnter={(e) =>
            (e.currentTarget.style.textDecoration = 'underline')
          }
          onMouseLeave={(e) => (e.currentTarget.style.textDecoration = 'none')}
        >
          GitHub
        </a>
        .
      </Typography>
    </Layout>
  );
}
