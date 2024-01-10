import React from 'react';
import Layout from '../components/Layout';
import Description from '../components/Description';
import ArticleList from '../components/ArticleList';
import Typography from '@mui/joy/Typography';

function Index() {
  return (
    <Layout>
      <Description />
      <Typography
        level="h2"
        component="h2"
        sx={{ fontWeight: 'bold', mb: 4, fontSize: '2rem' }}
      >
        Latest Articles
      </Typography>
      <ArticleList />
    </Layout>
  );
}

export default Index;
