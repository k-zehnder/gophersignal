import React from 'react';
import Layout from '../components/Layout';
import Description from '../components/Description';
import ArticleList from '../components/ArticleList';
import Typography from '@mui/joy/Typography';

// Define the Index component.
function Index() {
  return (
    <Layout>
      {/* Render the Description component. */}
      <Description />
      {/* Display the heading "Latest Articles". */}
      <Typography
        level="h2"
        component="h2"
        sx={{ fontWeight: 'bold', mb: 4, fontSize: '2rem' }}
      >
        Latest Articles
      </Typography>
      {/* Render the ArticleList component. */}
      <ArticleList />
    </Layout>
  );
}

export default Index;
