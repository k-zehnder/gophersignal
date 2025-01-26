import { useState, useEffect } from 'react';
import { processSummary, formatDate } from '../lib/stringUtils';
import { Article, ArticlesResponseSchema } from '../types';

// Custom React hook to fetch and manage a list of articles
const useArticles = () => {
  const [articles, setArticles] = useState<Article[]>([]);

  useEffect(() => {
    // Determine the API URL based on the environment
    const apiUrl =
      process.env.NEXT_PUBLIC_ENV === 'development'
        ? 'http://localhost:8080/api/v1/articles'
        : 'https://gophersignal.com/api/v1/articles';

    const fetchArticles = async () => {
      try {
        const response = await fetch(apiUrl);

        if (!response.ok) {
          throw new Error('Network response was not ok');
        }

        // Parse the API response
        const jsonResponse = await response.json();
        const apiResponse = ArticlesResponseSchema.parse(jsonResponse);

        // Transform articles data into the desired format
        const articlesData: Article[] = apiResponse.articles.map(
          (item: Article) => ({
            id: item.id,
            title: item.title,
            source: item.source,
            created_at: item.created_at
              ? formatDate(item.created_at)
              : undefined,
            updated_at: item.updated_at
              ? formatDate(item.updated_at)
              : undefined,
            summary: processSummary(
              item.summary ? { String: item.summary, Valid: true } : null
            ),
            link: item.link,
            upvotes: item.upvotes,
            comment_count: item.comment_count,
            comment_link: item.comment_link,
          })
        );

        setArticles(articlesData);
      } catch (error) {
        console.error('Error fetching articles:', error);
      }
    };

    fetchArticles();
  }, []);

  return articles;
};

export default useArticles;
