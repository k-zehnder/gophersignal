// Handles database connection, disconnection, and initialization operations.

import mysql, { Connection } from 'mysql2/promise';
import { Article, MySQLConfig } from '../types';

// Handles database connection, disconnection, and initialization operations
const connectToDatabase = async (mysqlConfig: MySQLConfig) => {
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
    const maxContentLength = 45000;
    const currentTimestamp = new Date()
      .toISOString()
      .slice(0, 19)
      .replace('T', ' ');

    const values = articles.map(
      ({
        title,
        link,
        content = '',
        summary,
        upvotes = 0,
        comment_count = 0,
        comment_link = '',
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
        currentTimestamp,
        currentTimestamp,
      ]
    );

    const query = `
    INSERT INTO articles (title, link, content, summary, source, upvotes, comment_count, comment_link, created_at, updated_at)
    VALUES ?
    ON DUPLICATE KEY UPDATE
      upvotes = VALUES(upvotes),
      comment_count = VALUES(comment_count),
      comment_link = VALUES(comment_link),
      updated_at = VALUES(updated_at);
  `;

    await connection.query(query, [values]);
  };

  // Updates the summary of an article
  const updateArticleSummary = async (
    id: number,
    summary: string
  ): Promise<void> => {
    await connection.execute('UPDATE articles SET summary = ? WHERE id = ?', [
      summary,
      id,
    ]);
  };

  // Closes the database connection.
  const closeDatabaseConnection = async (): Promise<void> => {
    if (connection) {
      await connection.end();
      console.log('Database connection closed');
    }
  };

  return {
    saveArticles,
    updateArticleSummary,
    closeDatabaseConnection,
    connection,
  };
};

export { connectToDatabase };
