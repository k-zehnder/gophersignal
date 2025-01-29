use crate::models::api_response::ApiResponse;
use crate::models::article::Article;
use reqwest::Client;

pub async fn fetch_articles() -> Result<Vec<Article>, reqwest::Error> {
    let client = Client::new();
    let response = client
        .get("https://gophersignal.com/api/v1/articles")
        .send()
        .await?;
    let api_response: ApiResponse = response.json().await?;
    Ok(api_response.articles)
}
