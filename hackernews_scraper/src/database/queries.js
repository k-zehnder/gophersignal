// Provides functions for database operations related to articles, including saving articles,
// updating article summaries, and fetching unsummarized articles.

// Inserts a new article into the database with the given connection and article details.
const saveArticle = async (connection, article) => {
  const query = `
    INSERT INTO articles (title, link, content, source, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?)
  `;

  // Format the current timestamp for SQL insertion
  const currentTimestamp = new Date()
    .toISOString()
    .slice(0, 19)
    .replace('T', ' ');

  // Execute the SQL query to save the article
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
const updateArticleSummary = async (connection, id, summary) => {
  await connection.execute('UPDATE articles SET summary = ? WHERE id = ?', [
    summary,
    id,
  ]);
};

// Retrieves articles from the database that do not have summaries, using the given connection.
const fetchUnsummarizedArticles = async (connection) => {
  const [rows] = await connection.execute(
    "SELECT id, content FROM articles WHERE (summary IS NULL OR summary = '') ORDER BY ID DESC LIMIT 40;"
  );
  return rows;
};

module.exports = {
  saveArticle,
  updateArticleSummary,
  fetchUnsummarizedArticles,
};
