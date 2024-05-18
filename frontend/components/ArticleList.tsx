import React from 'react';
import useArticles from '../hooks/useArticles';
import ArticleListItem from './ArticleListItem';
import List from '@mui/joy/List';
import { Article } from '../types';

// This component displays a list of articles fetched using the 'useArticles' hook.
function ArticleList() {
  const articles: Article[] = useArticles();

  if (!Array.isArray(articles)) {
    return <div>No articles available</div>;
  }

  return (
    <List sx={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      {/* Map through the articles and render each article as an 'ArticleListItem' component. */}
      {articles.map((article: Article) => (
        <ArticleListItem key={article.id} article={article} />
      ))}
    </List>
  );
}

export default ArticleList;
