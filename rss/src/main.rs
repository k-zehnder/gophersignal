mod config;
mod models;
mod routes;
mod services;

use crate::config::AppConfig;
use crate::routes::rss::generate_rss_feed;
use crate::services::articles::HttpArticlesClient;
use axum::{routing::get, Extension, Router};

#[tokio::main]
async fn main() {
    let config = AppConfig::from_env();

    // Create an instance of the client.
    let client = HttpArticlesClient;

    // Create the router and inject the dependencies.
    let app = Router::new()
        .route("/rss", get(generate_rss_feed::<HttpArticlesClient>))
        .layer(Extension(config.clone()))
        .layer(Extension(client));

    println!("Server running on port: {}", config.port);

    axum::Server::bind(&format!("0.0.0.0:{}", config.port).parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
