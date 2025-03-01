mod config;
mod errors;
mod models;
mod routes;
mod services;

use axum::{routing::get, Extension, Router};
use config::config::AppConfig;
use errors::errors::AppError;
use routes::rss::generate_rss_feed;
use services::articles::HttpArticlesClient;
use std::net::SocketAddr;

#[tokio::main]
async fn main() -> Result<(), AppError> {
    let config = AppConfig::from_env();
    let client = HttpArticlesClient;

    let app = Router::new()
        .route("/rss", get(generate_rss_feed::<HttpArticlesClient>))
        .layer(Extension(config.clone()))
        .layer(Extension(client));

    println!("Server running on port: {}", config.port);

    let addr: SocketAddr = format!("0.0.0.0:{}", config.port).parse()?;
    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}
