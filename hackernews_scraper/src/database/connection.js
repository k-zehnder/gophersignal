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

  // Select the database to use.
  await connection.query('USE gophersignal');

  // If in debug mode, reset the database.
  if (config.debugMode === 'true') {
    await resetDatabase(connection);
  }

  console.log('Database connected successfully');

  // Inserts a new article into the database with the given connection and article details.
  const saveArticle = async (article) => {
    const query = `
      INSERT INTO articles (title, link, content, source, created_at, updated_at)
      VALUES (?, ?, ?, ?, ?, ?)
    `;

    // Format the current timestamp for SQL insertion.
    const currentTimestamp = new Date()
      .toISOString()
      .slice(0, 19)
      .replace('T', ' ');

    // Truncate the content if it exceeds the column limit for TEXT.
    const maxContentLength = 20000;
    if (article.content.length > maxContentLength) {
      article.content = article.content.slice(0, maxContentLength);
    }

    // Execute the SQL query to save the article.
    await connection.execute(query, [
      article.title,
      article.link,
      article.content,
      'Hacker News',
      currentTimestamp,
      currentTimestamp,
    ]);
  };

  // Updates the summary of an article in the database with the given connection, article ID, and summary text.
  const updateArticleSummary = async (id, summary) => {
    await connection.execute('UPDATE articles SET summary = ? WHERE id = ?', [
      summary,
      id,
    ]);
  };

  // Retrieves articles from the database that do not have summaries.
  const fetchUnsummarizedArticles = async () => {
    const [rows] = await connection.execute(
      "SELECT id, content FROM articles WHERE (summary IS NULL OR summary = '') ORDER BY ID DESC LIMIT 40;"
    );
    return rows;
  };

  // Closes the database connection.
  const closeDatabaseConnection = async (connection) => {
    if (connection) {
      await connection.end();
      console.log('Database connection closed');
    }
  };

  // Resets the database by dropping and recreating the articles table.
  const resetDatabase = async (connection) => {
    await connection.query('DROP TABLE IF EXISTS articles');
    await connection.query(`
      CREATE TABLE articles (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        link VARCHAR(512) NOT NULL,
        content TEXT,
        summary VARCHAR(2000),
        source VARCHAR(100) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
      )
    `);
    console.log('Database reset successfully');
  };

  return {
    saveArticle,
    updateArticleSummary,
    fetchUnsummarizedArticles,
    closeDatabaseConnection,
    resetDatabase,
    connection,
  };
};

module.exports = {
  connectToDatabase,
};
