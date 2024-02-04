/* 
 * Articles Component: This component is responsible for displaying a list of articles fetched from an API.
*/
import React from 'react';
import Link from 'next/link';
import Typography from '@mui/joy/Typography';
import List from '@mui/joy/List';
import ListItem from '@mui/joy/ListItem';
import Layout from '../components/Layout';
import Description from '../components/Description';
import { formatDate, processSummary } from '../lib/stringUtils';

// Define the Article interface
interface Article {
  id: number;
  title: string;
  source: string;
  createdAt: string;
  updatedAt: string;
  formattedDate: string;
  summary: string | { String: string; Valid: boolean };
  link: string;
}

// Articles Component: Displays a list of articles fetched from an API
function Articles({ articles }: { articles: Article[] }) {
  return (
    <Layout>
      <Description />
      <Typography level="h2" component="h2" sx={{ fontWeight: 'bold', mb: 4, fontSize: '2rem' }}>
        Latest Articles
      </Typography>
      <List sx={{ display: "flex", flexDirection: 'column', gap: "1.5rem" }}>
        {articles.map((article, idx) => (
          <ListItem key={idx} sx={{ display: "flex", flexDirection: "column", alignItems: "flex-start" }}>
            {/* Display the formatted date */}
            <Typography sx={{ color: 'text.secondary', mb: '0.5rem', fontSize: '0.875rem' }}>
              {article.formattedDate} ⋅ {article.source}
            </Typography>
            <Typography level="h3" component="h3" sx={{ mb: '0.5rem', fontWeight: 'medium', fontSize: '1.5rem' }}>
              {/* Link to article with title */}
              <Link legacyBehavior href={article.link.replace(/\/\//g, '/')} passHref>
                <a target="_blank" rel="noopener noreferrer" style={{ textDecoration: 'none', color: '#007bff' }}>
                  {article.title}
                </a>
              </Link>
            </Typography>
            {/* Render article summary using processSummary function */}
            <Typography sx={{ fontSize: '1rem' }}>
              {typeof article.summary === 'string' ? article.summary : processSummary(article.summary)}
            </Typography>
          </ListItem>
        ))}
      </List>
    </Layout>
  );
}

// getServerSideProps: Server-side function for fetching articles data on each request
export async function getServerSideProps() {
  const apiUrl = process.env.NEXT_PUBLIC_ENV === "development"
    ? "http://localhost:8080/api/v1/articles"
    : "https://gophersignal.com/api/v1/articles";

  try {
    const response = await fetch(apiUrl);
    if (!response.ok) {
      throw new Error(`Failed to fetch articles: ${response.status}`);
    }
    const articlesData = await response.json();

    // Format the date for each article and add a new property 'formattedDate'
    const articles: Article[] = articlesData.map((article: any) => ({
      ...article,
      formattedDate: formatDate(article.created_at),
    }));

    return { props: { articles } };
  } catch (error) {
    console.error('Error fetching articles:', error);
    return { props: { articles: [] } };
  }
}

export default Articles;
