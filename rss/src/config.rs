use std::env;

pub struct AppConfig {
    pub port: String,
    pub api_url: String,
}

impl AppConfig {
    pub fn from_env() -> Self {
        dotenv::dotenv().ok();

        Self {
            port: env::var("RSS_PORT").unwrap_or_else(|_| "9090".to_string()),
            api_url: env::var("API_URL")
                .unwrap_or_else(|_| "https://gophersignal.com/api/v1/articles".to_string()),
        }
    }
}
