import mysql, { Connection } from 'mysql2/promise';
import { Article, Config, DBClient } from '../types';

const createMySqlClient = async (config: Config): Promise<DBClient> => {
  const mysqlConfig = config.mysql;

  // Create a new MySQL connection
  const connection: Connection = await mysql.createConnection({
    host: mysqlConfig.host,
    port: mysqlConfig.port,
    user: mysqlConfig.user,
    password: mysqlConfig.password,
    database: mysqlConfig.database,
  });

  console.log('Database connected successfully');

  // Inserts multiple articles into the database in bulk.
  const saveArticles = async (articles: Article[]): Promise<void> => {
    if (articles.length === 0) return;

    const maxContentLength = 45000;
    const currentTimestamp = new Date()
      .toISOString()
      .slice(0, 19)
      .replace('T', ' ');

    // Map articles to an array-of-arrays for bulk insert
    const values = articles.map(
      ({
        hnId,
        title,
        link,
        articleRank,
        content = '',
        summary,
        upvotes = 0,
        commentCount = 0,
        commentLink = '',
        flagged = false,
        dead = false,
        dupe = false,
        commitHash,
        modelName,
      }) => [
        hnId,
        title,
        link,
        articleRank,
        content.length > maxContentLength
          ? content.slice(0, maxContentLength)
          : content,
        summary,
        'Hacker News',
        upvotes,
        commentCount,
        commentLink,
        flagged,
        dead,
        dupe,
        commitHash,
        modelName,
        currentTimestamp,
        currentTimestamp,
      ]
    );

    const query = `INSERT INTO articles (
      hn_id,
      title,
      link,
      article_rank,
      content,
      summary,
      source,
      upvotes,
      comment_count,
      comment_link,
      flagged,
      dead,
      dupe,
      commit_hash,
      model_name,
      created_at,
      updated_at
    ) VALUES ?`;

    await connection.query(query, [values]);
  };

  // Updates an article's summary and timestamp.
  const updateArticleSummary = async (
    id: number,
    summary: string
  ): Promise<void> => {
    await connection.execute(
      'UPDATE articles SET summary = ?, updated_at = NOW() WHERE id = ?',
      [summary, id]
    );
  };

  // Marks an article as dead with updated timestamp.
  const markArticleAsDead = async (id: number): Promise<void> => {
    await connection.execute(
      'UPDATE articles SET dead = TRUE, updated_at = NOW() WHERE id = ?',
      [id]
    );
  };

  // Marks an article as duplicate with updated timestamp.
  const markArticleAsDuplicate = async (id: number): Promise<void> => {
    await connection.execute(
      'UPDATE articles SET dupe = TRUE, updated_at = NOW() WHERE id = ?',
      [id]
    );
  };

  // Closes the MySQL database connection.
  const closeDatabaseConnection = async (): Promise<void> => {
    await connection.end();
    console.log('Database connection closed');
  };

  return {
    saveArticles,
    updateArticleSummary,
    markArticleAsDead,
    markArticleAsDuplicate,
    closeDatabaseConnection,
    connection,
  };
};

export { createMySqlClient };
