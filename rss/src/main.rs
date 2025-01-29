mod config;
mod models;
mod routes;
mod services;

use crate::config::AppConfig;
use axum::{routing::get, Router};

#[tokio::main]
async fn main() {
    // Load environment variables and app configuration
    let config = AppConfig::from_env();

    // Create Axum router
    let app = Router::new().route("/rss", get(routes::rss::generate_rss_feed));

    println!("Server running on port: {}", config.port);

    // Start server
    axum::Server::bind(&format!("0.0.0.0:{}", config.port).parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
