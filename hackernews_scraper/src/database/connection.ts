// Handles database connection, disconnection, and initialization operations.

import mysql, { Connection } from 'mysql2/promise';
import { Article, MySQLConfig } from '../types';

// Handles database connection, disconnection, and initialization operations.
const connectToDatabase = async (mysqlConfig: MySQLConfig) => {
  // Establish a new connection to the MySQL database using the provided configuration.
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
    const maxContentLength = 45000;
    const currentTimestamp = new Date()
      .toISOString()
      .slice(0, 19)
      .replace('T', ' ');
    const values = articles.map(({ title, link, content = '', summary }) => [
      title,
      link,
      content.length > maxContentLength
        ? content.slice(0, maxContentLength)
        : content,
      summary,
      'Hacker News',
      currentTimestamp,
      currentTimestamp,
    ]);

    const query = `
    INSERT INTO articles (title, link, content, summary, source, created_at, updated_at)
    VALUES ?
  `;

    await connection.query(query, [values]);
  };

  // Updates the summary of an article in the database with the given connection, article ID, and summary text.
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
