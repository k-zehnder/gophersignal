// Handles database connection, disconnection, and initialization operations.

const mysql = require('mysql2/promise');

// Handles database connection, disconnection, and initialization operations.
const connectToDatabase = async (config) => {
  // Establish a new connection to the MySQL database using the provided configuration.
  const connection = await mysql.createConnection({
    host: config.mysql.host,
    port: config.mysql.port,
    user: config.mysql.user,
    password: config.mysql.password,
    database: config.mysql.database,
  });

  console.log('Database connected successfully');

  // Inserts multiple articles into the database in bulk.
  const saveArticles = async (articles) => {
    const maxContentLength = 45000;
    const currentTimestamp = new Date()
      .toISOString()
      .slice(0, 19)
      .replace('T', ' ');
    const values = articles.map(({ title, link, content, summary }) => [
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
  const updateArticleSummary = async (id, summary) => {
    await connection.execute('UPDATE articles SET summary = ? WHERE id = ?', [
      summary,
      id,
    ]);
  };

  // Closes the database connection.
  const closeDatabaseConnection = async () => {
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

module.exports = {
  connectToDatabase,
};
