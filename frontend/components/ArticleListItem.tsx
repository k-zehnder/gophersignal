import React from 'react';
import { Article } from '../types';
import Link from 'next/link';
import Typography from '@mui/joy/Typography';
import ListItem from '@mui/joy/ListItem';

// Define the props interface for ArticleListItem.
interface ArticleListItemProps {
  article: Article;
}

// ArticleListItem component displays a single article as a list item.
const ArticleListItem: React.FC<ArticleListItemProps> = ({ article }) => {
  return (
    <ListItem
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'flex-start',
      }}
    >
      {/* Display article metadata (updated date and source). */}
      <Typography
        sx={{ color: 'text.secondary', mb: '0.5rem', fontSize: '0.875rem' }}
      >
        {article.updatedAt} ⋅ {article.source}
      </Typography>

      {/* Display the article title as a link to the article. */}
      <Typography
        level="h3"
        component="h3"
        sx={{ mb: '0.5rem', fontWeight: 'medium', fontSize: '1.5rem' }}
      >
        {/* Use 'next/link' to create a link to the article. */}
        <Link legacyBehavior href={article.link} passHref>
          <a
            target="_blank"
            rel="noopener noreferrer"
            style={{ textDecoration: 'none', color: '#007bff' }}
          >
            {article.title}
          </a>
        </Link>
      </Typography>

      {/* Display the article summary. */}
      <Typography sx={{ fontSize: '1rem' }}>{article.summary}</Typography>
    </ListItem>
  );
};

export default ArticleListItem;
