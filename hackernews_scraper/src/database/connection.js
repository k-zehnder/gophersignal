// Handles database connection, disconnection, and initialization operations.

const mysql = require('mysql2/promise');

// Establishes a connection to the database using environment variables.
const connectToDatabase = async () => {
  const connection = await mysql.createConnection({
    host: process.env.MYSQL_HOST,
    port: process.env.MYSQL_PORT,
    user: process.env.MYSQL_USER,
    password: process.env.MYSQL_PASSWORD,
    database: process.env.MYSQL_DATABASE,
  });

  await connection.query('USE gophersignal');
  console.log('Database connected successfully');
  return connection;
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

// Initializes the database by connecting and optionally resetting it based on the configuration.
const initializeDatabase = async (config) => {
  const connection = await connectToDatabase();
  if (config.debugMode) {
    await resetDatabase(connection);
  }
  return connection;
};

module.exports = {
  connectToDatabase,
  closeDatabaseConnection,
  resetDatabase,
  initializeDatabase,
};
