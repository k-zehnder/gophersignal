use axum::{http::StatusCode, response::IntoResponse, routing::get, Router};
use chrono::{Duration, Utc};
use dotenv::dotenv;
use reqwest::Client;
use rss::{ChannelBuilder, ItemBuilder};
use serde::Deserialize;
use std::env;

// Full API response with metadata and articles
#[derive(Deserialize)]
struct ApiResponse {
    code: u32,
    status: String,
    total_count: u32,
    articles: Option<Vec<Article>>,
}

// Data for an individual article
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

// Generate the RSS feed from the backend API
async fn generate_rss_feed() -> Result<impl IntoResponse, StatusCode> {
    let api_url = "http://backend:8080/api/v1/articles";

    let client = Client::new();
    let response = client.get(api_url).send().await.map_err(|err| {
        eprintln!("Failed to fetch articles: {}", err);
        StatusCode::INTERNAL_SERVER_ERROR
    })?;

    let api_response: ApiResponse = response.json().await.map_err(|err| {
        eprintln!("Failed to parse JSON response: {}", err);
        StatusCode::INTERNAL_SERVER_ERROR
    })?;

    let mut articles = api_response.articles.unwrap_or_else(Vec::new);

    // Sort articles by `id` DESC
    articles.sort_by(|a, b| b.id.cmp(&a.id));

    // Assign unique pubDates in descending order
    let now = Utc::now();

    let items: Vec<_> = articles
        .into_iter()
        .enumerate()
        .map(|(i, article)| {
            // Offset in whole MINUTES for each subsequent item
            // i=0 => now, i=1 => now - 1 minute, etc.
            let pub_date = (now - Duration::minutes(i as i64)).to_rfc2822();

            let description = format!(
                "Summary: {}<br><br>Upvotes: {}<br><br>Comments: {} [<a href=\"{}\">View Comments</a>]<br><br>\
                 Link: <a href=\"{}\">{}</a>",
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
                .pub_date(Some(pub_date))
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
