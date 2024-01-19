import { useState, useEffect } from 'react';
import { ApiArticle, Article } from '../types';
import { processSummary, formatDate } from '../lib/stringUtils';

const useArticles = () => {
  const [articles, setArticles] = useState<Article[]>([]);

  useEffect(() => {
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
        const apiArticles = (await response.json()) as ApiArticle[];

        const articlesData: Article[] = apiArticles.map((item) => ({
          id: item.id,
          title: item.title,
          source: item.source,
          createdAt: formatDate(item.createdAt),
          updatedAt: formatDate(item.updatedAt),
          summary: processSummary(item.summary),
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

  return articles;
};

export default useArticles;
