/**
 * useArticles: A custom React hook for fetching and managing a list of articles.
 * This hook encapsulates the logic for asynchronously loading articles from a backend API,
 * including error handling and state management. It leverages the useEffect hook to fetch
 * data when the component mounts and useState to store the fetched articles.
 */
import { useState, useEffect } from 'react';
import { Article } from '../types'; 
import { processSummary, formatDate } from '../lib/stringUtils';

// Custom React hook to fetch and manage a list of articles
const useArticles = () => {
  // State to store the list of articles
  const [articles, setArticles] = useState<Article[]>([]);

  // Use effect to fetch articles when the component mounts
  useEffect(() => {
    // Determine the API URL for fetching articles
    const apiUrl =
      process.env.NEXT_PUBLIC_ENV === 'development'
        ? 'http://localhost:8080/api/v1/articles'
        : 'https://gophersignal.com/api/v1/articles';

    // Function to fetch articles from the backend
    const fetchArticles = async () => {
      try {
        // Fetch articles from the backend
        const response = await fetch(apiUrl);
        
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        
        // Parse backend response into articles data
        const apiArticles = await response.json();
        const articlesData: Article[] = apiArticles.map((item: any) => ({
          id: item.id,
          title: item.title,
          source: item.source,
          createdAt: formatDate(item.created_at),
          updatedAt: formatDate(item.updated_at),
          summary: processSummary(item.summary),
          link: item.link,
        }));

        // Update the state with the fetched articles
        setArticles(articlesData);
      } catch (error) {
        // Handle any errors that occur during fetching
        console.error('Error fetching articles:', error);
      }
    };

    // Call the fetchArticles function
    fetchArticles();
  }, []);

  // Return the list of articles
  return articles;
};

export default useArticles;
