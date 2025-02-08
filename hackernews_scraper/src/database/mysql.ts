import mysql, { Connection } from 'mysql2/promise';
import { Article, Config, DBClient } from '../types';

const createMySqlClient = async (config: Config): Promise<DBClient> => {
  const mysqlConfig = config.mysql;

  const connection: Connection = await mysql.createConnection({
    host: mysqlConfig.host,
    port: mysqlConfig.port,
    user: mysqlConfig.user,
    password: mysqlConfig.password,
    database: mysqlConfig.database,
  });

  console.log('Database connected successfully');

  // Inserts multiple articles into the database in bulk
  const saveArticles = async (articles: Article[]): Promise<void> => {
    if (articles.length === 0) return;

    const maxContentLength = 45000;
    const currentTimestamp = new Date()
      .toISOString()
      .slice(0, 19)
      .replace('T', ' ');

    // Map articles to an array-of-arrays for the insert
    const values = articles.map(
      ({
        title,
        link,
        content = '',
        summary,
        upvotes = 0,
        comment_count = 0,
        comment_link = '',
        flagged = false,
        dead = false,
        dupe = false,
      }) => [
        title,
        link,
        content.length > maxContentLength
          ? content.slice(0, maxContentLength)
          : content,
        summary,
        'Hacker News',
        upvotes,
        comment_count,
        comment_link,
        flagged,
        dead,
        dupe,
        currentTimestamp,
        currentTimestamp,
      ]
    );

    const query = `INSERT INTO articles (
      title,
      link,
      content,
      summary,
      source,
      upvotes,
      comment_count,
      comment_link,
      flagged,
      dead,
      dupe,
      created_at,
      updated_at
    ) VALUES ?`;

    await connection.query(query, [values]);
  };

  const updateArticleSummary = async (
    id: number,
    summary: string
  ): Promise<void> => {
    await connection.execute(
      'UPDATE articles SET summary = ?, updated_at = NOW() WHERE id = ?',
      [summary, id]
    );
  };

  const markArticleAsDead = async (id: number): Promise<void> => {
    await connection.execute(
      'UPDATE articles SET dead = TRUE, updated_at = NOW() WHERE id = ?',
      [id]
    );
  };

  const markArticleAsDuplicate = async (id: number): Promise<void> => {
    await connection.execute(
      'UPDATE articles SET dupe = TRUE, updated_at = NOW() WHERE id = ?',
      [id]
    );
  };

  const closeDatabaseConnection = async (): Promise<void> => {
    if (connection) {
      await connection.end();
      console.log('Database connection closed');
    }
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
