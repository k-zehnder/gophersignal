import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Typography from '@mui/joy/Typography';
import List from '@mui/joy/List';
import ListItem from '@mui/joy/ListItem';
import Layout from "../components/Layout"; 

interface Article {
  id: number;
  title: string;
  source: string;
  scrapedAt: string;
  summary: { String: string; Valid: boolean }; // Updated summary field
}

function Articles() {
  const [articles, setArticles] = useState<Article[]>([]);

  useEffect(() => {
    fetch('https://gophersignal.com/articles')
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then(data => {
        setArticles(data);
      })
      .catch(error => {
        console.error('Error fetching articles:', error);
      });
  }, []);

  return (
    <Layout>
      <Typography level="h1" component="h2" sx={{ fontWeight: 'bold', mb: 2 }}>
        Latest Articles
      </Typography>

      <List sx={{ display: "flex", flexDirection: 'column', gap: "1rem" }}>
        {articles.map(({ id, title, source, scrapedAt, summary }, index) => (
          <ListItem key={index} sx={{ display: "flex", flexDirection: "column", alignItems: "flex-start" }}>
            <Typography level="body2" sx={{ color: 'text.secondary', mb: '0.5rem' }}>
              "Today" ⋅ {source}
            </Typography>
            
            <Typography level="h4" component="h3" sx={{ mb: '0.5rem', fontWeight: 'medium' }}>
              <Link legacyBehavior href={`/article/${id}`} passHref>
                <a>{title}</a>
              </Link>
            </Typography>
            
            <Typography level="body2">
              {summary.Valid ? summary.String : 'No summary available'}
            </Typography>
          </ListItem>
        ))}
      </List>
    </Layout>
  );
}

export default Articles;
