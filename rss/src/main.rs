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

    // Get the PORT value from the environment or default to 9090
    let port = env::var("RSS_PORT").unwrap_or_else(|_| "9090".to_string());

    // Create the Axum router
    let app = Router::new().route("/rss", get(generate_rss_feed));

    println!("RSS Service is running on port: {}", port);

    // Start the server
    axum::Server::bind(&format!("0.0.0.0:{}", port).parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}

#[derive(Deserialize)]
struct ApiResponse {
    code: u32,
    status: String,
    total_count: u32,
    articles: Vec<Article>,
}

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

// Generate RSS feed
async fn generate_rss_feed() -> Result<impl IntoResponse, StatusCode> {
    let api_url = "http://backend:8080/api/v1/articles";

    // Fetch data from the API
    let client = Client::new();
    let response = client.get(api_url).send().await.map_err(|err| {
        eprintln!("Failed to fetch articles: {}", err);
        StatusCode::INTERNAL_SERVER_ERROR
    })?;

    // Parse the JSON response
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
                "ID: {}\nTitle: {}\nSummary: {}\nCreated At: {}\nUpvotes: {}\nComments: {} [View Comments]({})\nLink: {}",
                article.id,
                article.title,
                article.summary,
                article.created_at,
                article.upvotes,
                article.comment_count,
                article.comment_link,
                article.link,
            );

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

    // Return RSS XML
    Ok(axum::response::Html(channel.to_string()))
}

// Test for the RSS feed
#[tokio::test]
async fn test_rss_feed() {
    use hyper::Server;
    use std::net::SocketAddr;

    // Create the router
    let app = Router::new().route("/rss", get(generate_rss_feed));

    // Bind to a random available port for testing
    let addr = SocketAddr::from(([127, 0, 0, 1], 0));
    let server = Server::bind(&addr).serve(app.into_make_service());
    let bound_addr = server.local_addr();

    // Run the server in the background
    tokio::spawn(server);

    // Send a request to the test server
    let response = reqwest::get(&format!("http://{}/rss", bound_addr))
        .await
        .unwrap();

    // Assert the response is successful and contains RSS XML
    assert_eq!(response.status(), 200);
    let body = response.text().await.unwrap();
    assert!(body.contains("<rss"));
    assert!(body.contains("Summary:"));
    assert!(body.contains("Upvotes:"));
    assert!(body.contains("Comments:"));
}
