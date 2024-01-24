import React from 'react';
import { Article } from '../types';
import Link from 'next/link';
import Typography from '@mui/joy/Typography';
import ListItem from '@mui/joy/ListItem';

interface ArticleListItemProps {
  article: Article;
}

const ArticleListItem: React.FC<ArticleListItemProps> = ({ article }) => {
  console.log(article)
  return (
    <ListItem sx={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start' }}>
      <Typography sx={{ color: 'text.secondary', mb: '0.5rem', fontSize: '0.875rem' }}>
        {article.updatedAt} ⋅ {article.source}
      </Typography>
      <Typography level="h3" component="h3" sx={{ mb: '0.5rem', fontWeight: 'medium', fontSize: '1.5rem' }}>
        <Link legacyBehavior href={article.link} passHref>
          <a target="_blank" rel="noopener noreferrer" style={{ textDecoration: 'none', color: '#007bff' }}>
            {article.title}
          </a>
        </Link>
      </Typography>
      <Typography sx={{ fontSize: '1rem' }}>
        {article.summary}
      </Typography>
    </ListItem>
  );
};

export default ArticleListItem;
