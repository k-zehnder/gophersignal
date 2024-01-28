import React from 'react';
import useArticles from '../hooks/useArticles';
import ArticleListItem from './ArticleListItem';
import List from '@mui/joy/List';

// This component displays a list of articles fetched using the 'useArticles' hook.
function ArticleList() {
  // Fetch the list of articles using the 'useArticles' hook.
  const articles = useArticles();

  return (
    // Render a list of articles using MUI List component.
    <List sx={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      {/* Map through the articles and render each article as an 'ArticleListItem' component. */}
      {articles.map(article => (
        <ArticleListItem key={article.id} article={article} />
      ))}
    </List>
  );
}

export default ArticleList;
