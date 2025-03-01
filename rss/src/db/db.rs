use crate::config::config::AppConfig;
use sqlx::MySqlPool;

/// Creates and returns a MySQL connection pool.
pub async fn create_pool(config: &AppConfig) -> Result<MySqlPool, sqlx::Error> {
    MySqlPool::connect(&config.database_url).await
}
