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
      {/* Wrap everything in a container that ensures content pushes the footer down */}
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          minHeight: '100vh',
        }}
      >
        {/* Main content fills available space */}
        <div style={{ flex: 1 }}>
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
        </div>

        {/* Footer stays at the bottom */}
        <Footer />
      </div>
    </Layout>
  );
};

export default Index;
