import React from 'react';
import useArticles from '../hooks/useArticles';
import ArticleListItem from './ArticleListItem';
import List from '@mui/joy/List';

function ArticleList() {
  const articles = useArticles();

  return (
    <List sx={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      {articles.map(article => (
        <ArticleListItem key={article.id} article={article} />
      ))}
    </List>
  );
}

export default ArticleList;
