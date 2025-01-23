import React from 'react';
import { Article } from '../types';
import Link from 'next/link';
import Typography from '@mui/joy/Typography';
import ListItem from '@mui/joy/ListItem';
import Box from '@mui/joy/Box';

interface ArticleListItemProps {
  article: Article;
}

const ArticleListItem: React.FC<ArticleListItemProps> = ({ article }) => {
  return (
    <ListItem
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'flex-start',
        gap: '0.5rem',
      }}
    >
      {/* Date and Source */}
      <Typography sx={{ color: 'text.secondary', fontSize: '0.875rem' }}>
        {article.updated_at} &middot; {article.source}
      </Typography>

      {/* Title */}
      <Typography
        level="h3"
        component="h3"
        sx={{ fontWeight: 'medium', fontSize: '1.5rem', mb: 0 }}
      >
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

      {/* Summary */}
      <Typography sx={{ fontSize: '1rem' }}>{article.summary}</Typography>

      {/* Upvotes & Comments Row */}
      {(article.upvotes || article.comment_count) && (
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: '1rem',
            mt: '0.25rem',
          }}
        >
          {/* Upvotes */}
          {typeof article.upvotes === 'number' && (
            <Typography level="body2" sx={{ color: 'text.secondary' }}>
              {article.upvotes} upvotes
            </Typography>
          )}

          {/* Comments (Linked to HN comment thread) */}
          {article.comment_count && article.comment_link && (
            <Link legacyBehavior href={article.comment_link} passHref>
              <Typography
                component="a"
                sx={{
                  fontSize: '0.875rem',
                  color: '#007bff',
                  textDecoration: 'none',
                  '&:hover': {
                    textDecoration: 'underline',
                  },
                }}
              >
                {article.comment_count} comments
              </Typography>
            </Link>
          )}
        </Box>
      )}
    </ListItem>
  );
};

export default ArticleListItem;
