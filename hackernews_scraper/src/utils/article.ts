import { Article } from '../types';

const createArticleHelpers = () => {
  const categorizeArticles = (articles: Article[]) => {
    return articles.reduce(
      (acc, article) => {
        if (article.flagged) acc.flagged.push(article);
        if (article.dead) acc.dead.push(article);
        if (article.dupe) acc.dupe.push(article);
        return acc;
      },
      { flagged: [] as Article[], dead: [] as Article[], dupe: [] as Article[] }
    );
  };

  const getTopArticlesWithContent = (
    processedArticles: Article[],
    topArticles: Article[]
  ): Required<Article>[] => {
    return processedArticles.filter(
      (article): article is Required<Article> =>
        topArticles.some((top) => top.link === article.link) &&
        article.content !== undefined
    );
  };

  return {
    categorizeArticles,
    getTopArticlesWithContent,
  };
};

export type ArticleHelpers = ReturnType<typeof createArticleHelpers>;

export default createArticleHelpers;
