import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Typography from '@mui/joy/Typography';
import List from '@mui/joy/List';
import ListItem from '@mui/joy/ListItem';
import Layout from "../components/Layout"; 
import Description from "../components/Description"; 

interface Article {
  id: number;
  title: string;
  source: string;
  scrapedAt: string;
  summary: { String: string; Valid: boolean };
  link: string;
  isOnHomepage: boolean; 
}

function Articles() {
  const [articles, setArticles] = useState<Article[]>([]);

  useEffect(() => {
    const apiUrl = process.env.NEXT_PUBLIC_ENV === "development" 
      ? "http://localhost:8080/api/v1/articles" 
      : "https://gophersignal.com/api/v1/articles";

    // Fetch articles
    fetch(apiUrl)
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then(data => setArticles(data))
      .catch(error => console.error('Error fetching articles:', error));
  }, []);

  // Function to format date
  const formatDate = (dateStr: string): string => {
    const date = new Date(dateStr);
    return date.toLocaleDateString(undefined, {
      year: 'numeric', month: 'long', day: 'numeric'
    });
  };

  return (
    <Layout>
      <Description />
      <Typography level="h2" component="h2" sx={{ fontWeight: 'bold', mb: 4, fontSize: '2rem' }}>
        Latest Articles
      </Typography>

      <List sx={{ display: "flex", flexDirection: 'column', gap: "1.5rem" }}>
        {articles.filter(article => article.isOnHomepage).map((article, idx) => (
          <ListItem key={idx} sx={{ display: "flex", flexDirection: "column", alignItems: "flex-start" }}>
            <Typography level="body2" sx={{ color: 'text.secondary', mb: '0.5rem', fontSize: '0.875rem' }}>
              {formatDate(article.scrapedAt)} ⋅ {article.source}
            </Typography>

            <Typography level="h3" component="h3" sx={{ mb: '0.5rem', fontWeight: 'medium', fontSize: '1.5rem' }}>
              <Link legacyBehavior href={article.link} passHref>
                <a target="_blank" rel="noopener noreferrer" style={{ textDecoration: 'none', color: '#007bff' }}>
                  {article.title}
                </a>
              </Link>
            </Typography>

            <Typography level="body2" sx={{ fontSize: '1rem' }}>
              {article.summary.Valid ? article.summary.String : 'No summary available'}
            </Typography>
          </ListItem>
        ))}
      </List>
    </Layout>
  );
}

export default Articles;

