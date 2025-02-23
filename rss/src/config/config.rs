use std::env;

#[derive(Clone)]
pub struct AppConfig {
    pub port: String,
    pub api_url: String,
    pub database_url: String,
}

impl AppConfig {
    pub fn from_env() -> Self {
        // Only load .env if not running tests.
        #[cfg(not(test))]
        {
            dotenv::dotenv().ok();
        }
        Self {
            port: env::var("RSS_PORT").unwrap_or_else(|_| "9090".to_string()),
            api_url: env::var("API_URL")
                .unwrap_or_else(|_| "https://gophersignal.com/api/v1/articles".to_string()),
            database_url: env::var("DATABASE_URL")
                .unwrap_or_else(|_| "mysql://user:password@host:3306/gophersignal".to_string()),
        }
    }
}
