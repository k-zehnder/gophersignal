import React from 'react';
import Layout from '../components/Layout';
import Description from '../components/Description';
import ArticleList from '../components/ArticleList';

function Articles() {
  return (
    <Layout>
      <Description />
      <ArticleList />
    </Layout>
  );
}

export default Articles;
