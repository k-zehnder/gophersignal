/* 
 * Articles Component: This component is responsible for displaying a list of articles fetched from an API.
*/
import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Typography from '@mui/joy/Typography';
import List from '@mui/joy/List';
import ListItem from '@mui/joy/ListItem';
import Layout from '../components/Layout';
import Description from '../components/Description';

interface Article {
  id: number;
  title: string;
  source: string;
  createdAt: string;
  updatedAt: string;
  summary: string;
  link: string;
  isOnHomepage: boolean;
}

// Articles Component: This component is responsible for displaying a list of articles fetched from an API.
function Articles() {
  const [articles, setArticles] = useState<Article[]>([]);

  useEffect(() => {
    const fetchArticles = async () => {
      // Fetching articles based on the environment (development or production)
      const apiUrl = process.env.NEXT_PUBLIC_ENV === 'development'
        ? 'http://localhost:8080/api/v1/articles'
        : 'https://gophersignal.com/api/v1/articles';

      try {
        const response = await fetch(apiUrl);
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const data = await response.json();

        // Mapping API data to Article objects
        const articlesData: Article[] = data.map((item: any) => ({
          id: item.id,
          title: item.title,
          source: item.source,
          createdAt: item.createdAt,
          updatedAt: item.updatedAt,
          summary: item.summary && item.summary.Valid ? item.summary.String : 'No summary available',
          link: item.link,
          isOnHomepage: item.is_on_homepage,
        }));

        setArticles(articlesData);
      } catch (error) {
        console.error('Error fetching articles:', error);
      }
    };

    fetchArticles();
  }, []);

  const formatDate = (dateStr: string): string => {
    // Formatting date strings for better readability
    if (!dateStr) return 'Date not available';
    const date = new Date(dateStr);
    return isNaN(date.getTime()) ? 'Invalid Date' : date.toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' });
  };

  return (
    <Layout>
      <Description />
      <Typography level="h2" component="h2" sx={{ fontWeight: 'bold', mb: 4, fontSize: '2rem' }}>
        Latest Articles
      </Typography>

      <List sx={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
        {articles.map((article, idx) => (
          <ListItem key={idx} sx={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start' }}>
            <Typography sx={{ color: 'text.secondary', mb: '0.5rem', fontSize: '0.875rem' }}>
              {formatDate(article.createdAt)} ⋅ {article.source}
            </Typography>

            <Typography level="h3" component="h3" sx={{ mb: '0.5rem', fontWeight: 'medium', fontSize: '1.5rem' }}>
              <Link legacyBehavior href={article.link} passHref>
                <a target="_blank" rel="noopener noreferrer" style={{ textDecoration: 'none', color: '#007bff' }}>
                  {article.title}
                </a>
              </Link>
            </Typography>

            <Typography sx={{ fontSize: '1rem' }}>
              {article.summary !== 'No summary available' ? article.summary : <em>No summary available</em>}
            </Typography>
          </ListItem>
        ))}
      </List>
    </Layout>
  );
}

export default Articles;
