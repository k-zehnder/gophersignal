import React from 'react';
import Layout from '../components/Layout';
import Description from '../components/Description';
import ArticleList from '../components/ArticleList';
import Typography from '@mui/joy/Typography';
import Footer from '../components/Footer';

// Index component renders the main page layout, including description and article list.
const Index: React.FC = () => {
  return (
    <Layout>
      {/* Description of the website or application. */}
      <Description />

      {/* Heading for the latest articles section. */}
      <Typography
        level="h2"
        component="h2"
        sx={{ fontWeight: 'bold', mb: 4, fontSize: '2rem' }}
      >
        Latest Articles
      </Typography>

      {/* List of the latest articles. */}
      <div style={{ marginBottom: '16px' }}>
        <ArticleList />
      </div>

      {/* Footer */}
      <Footer />
    </Layout>
  );
};

export default Index;
