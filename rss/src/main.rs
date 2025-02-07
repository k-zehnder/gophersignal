mod config;
mod models;
mod routes;
mod services;

use crate::config::AppConfig;
use axum::{routing::get, Extension, Router};

#[tokio::main]
async fn main() {
    // Load environment variables and app configuration
    let config = AppConfig::from_env();

    // Create Axum router and add the config as an extension
    let app = Router::new()
        .route("/rss", get(routes::rss::generate_rss_feed))
        .layer(Extension(config.clone()));

    println!("Server running on port: {}", config.port);

    // Start server
    axum::Server::bind(&format!("0.0.0.0:{}", config.port).parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
