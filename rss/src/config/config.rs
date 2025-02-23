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

#[cfg(test)]
mod tests {
    use super::*;
    use serial_test::serial;
    use std::env;

    #[test]
    fn test_app_config_defaults() {
        env::remove_var("RSS_PORT");
        env::remove_var("API_URL");
        env::remove_var("DATABASE_URL");

        let config = AppConfig::from_env();
        assert_eq!(config.port, "9090");
        assert_eq!(config.api_url, "https://gophersignal.com/api/v1/articles");
        assert_eq!(
            config.database_url,
            "mysql://user:password@host:3306/gophersignal"
        );
    }

    #[test]
    #[serial]
    fn test_app_config_custom() {
        env::remove_var("RSS_PORT");
        env::remove_var("API_URL");
        env::remove_var("DATABASE_URL");

        env::set_var("RSS_PORT", "8000");
        env::set_var("API_URL", "http://example.com/api");
        env::set_var("DATABASE_URL", "mysql://custom:secret@dbhost:3306/customdb");

        let config = AppConfig::from_env();
        assert_eq!(config.port, "8000");
        assert_eq!(config.api_url, "http://example.com/api");
        assert_eq!(
            config.database_url,
            "mysql://custom:secret@dbhost:3306/customdb"
        );

        env::remove_var("RSS_PORT");
        env::remove_var("API_URL");
        env::remove_var("DATABASE_URL");
    }
}
