use crate::config::AppConfig;
use crate::models::article::{ApiResponse, Article};
use crate::routes::rss::RssQuery;
use async_trait::async_trait;
use reqwest::Client;

#[async_trait]
pub trait ArticlesClient: Send + Sync {
    async fn fetch_articles(
        &self,
        query: &RssQuery,
        config: &AppConfig,
    ) -> Result<Vec<Article>, Box<dyn std::error::Error + Send + Sync>>;
}

#[derive(Clone)]
pub struct HttpArticlesClient;

#[async_trait]
impl ArticlesClient for HttpArticlesClient {
    async fn fetch_articles(
        &self,
        query: &RssQuery,
        config: &AppConfig,
    ) -> Result<Vec<Article>, Box<dyn std::error::Error + Send + Sync>> {
        let client = Client::new();
        let backend_url = config.api_url.clone();
        let mut request = client.get(&backend_url);

        // Build query parameters based on RssQuery.
        let mut params = Vec::new();
        if let Some(flagged) = query.flagged {
            params.push(("flagged", flagged.to_string()));
        }
        if let Some(dead) = query.dead {
            params.push(("dead", dead.to_string()));
        }
        if let Some(dupe) = query.dupe {
            params.push(("dupe", dupe.to_string()));
        }
        if !params.is_empty() {
            request = request.query(&params);
        }

        let response = request.send().await?;
        let api_response: ApiResponse = response.json().await?;
        Ok(api_response.articles.unwrap_or_default())
    }
}
