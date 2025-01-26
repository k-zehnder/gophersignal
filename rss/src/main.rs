use axum::{http::StatusCode, response::IntoResponse, routing::get, Router};
use dotenv::dotenv;
use reqwest::Client;
use rss::{ChannelBuilder, ItemBuilder};
use serde::Deserialize;
use std::env;

#[tokio::main]
async fn main() {
    // Load environment variables from the .env file
    dotenv().ok();

    // Retrieve the PORT value or default to 9090
    let port = env::var("RSS_PORT").unwrap_or_else(|_| "9090".to_string());

    // Configure the Axum router with the /rss route
    let app = Router::new().route("/rss", get(generate_rss_feed));

    println!("RSS Service is running on port: {}", port);

    // Start the server and bind it to the specified address
    axum::Server::bind(&format!("0.0.0.0:{}", port).parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}

// Structure representing the API response
#[derive(Deserialize)]
struct ApiResponse {
    code: u32,
    status: String,
    total_count: u32,
    articles: Vec<Article>,
}

// Structure representing an individual article
#[derive(Deserialize)]
struct Article {
    id: u32,
    title: String,
    link: String,
    summary: String,
    created_at: String,
    upvotes: u32,
    comment_count: u32,
    comment_link: String,
}

// Generate the RSS feed from the backend API
async fn generate_rss_feed() -> Result<impl IntoResponse, StatusCode> {
    let api_url = "http://backend:8080/api/v1/articles";

    // Fetch data from the backend API
    let client = Client::new();
    let response = client.get(api_url).send().await.map_err(|err| {
        eprintln!("Failed to fetch articles: {}", err);
        StatusCode::INTERNAL_SERVER_ERROR
    })?;

    // Deserialize JSON response into the ApiResponse struct
    let api_response: ApiResponse = response.json().await.map_err(|err| {
        eprintln!("Failed to parse JSON response: {}", err);
        StatusCode::INTERNAL_SERVER_ERROR
    })?;

    // Map articles to RSS items
    let items: Vec<_> = api_response
        .articles
        .into_iter()
        .map(|article| {
            let description = format!(
                "Summary: {}<br><br>Upvotes: {}<br><br>Comments: {} [<a href=\"{}\">View Comments</a>]<br><br>Link: <a href=\"{}\">{}</a>",
                article.summary,
                article.upvotes,
                article.comment_count,
                article.comment_link,
                article.link,
                article.link
            );

            // Build each RSS item
            ItemBuilder::default()
                .title(Some(article.title))
                .link(Some(article.link))
                .description(Some(description))
                .pub_date(Some(article.created_at))
                .build()
        })
        .collect();

    // Build the RSS channel
    let channel = ChannelBuilder::default()
        .title("GopherSignal RSS Feed")
        .link("https://gophersignal.com")
        .description("Latest articles from GopherSignal")
        .items(items)
        .build();

    // Return the RSS feed as HTML
    Ok(axum::response::Html(channel.to_string()))
}
