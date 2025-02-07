use crate::config::AppConfig;
use crate::models::article::{ApiResponse, Article};
use crate::routes::rss::RssQuery;
use reqwest::Client;

// Fetches articles from the backend service, applying optional query parameters.
pub async fn fetch_articles(
    query: &RssQuery,
    config: &AppConfig,
) -> Result<Vec<Article>, reqwest::Error> {
    let client = Client::new();

    // Use the API URL from the configuration
    let backend_url = config.api_url.clone();
    let mut request = client.get(&backend_url);

    // Prepare query parameters from the RssQuery struct
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

    // Attach the query parameters if any exist
    if !params.is_empty() {
        request = request.query(&params);
    }

    let response = request.send().await?;
    let api_response: ApiResponse = response.json().await?;
    Ok(api_response.articles)
}
