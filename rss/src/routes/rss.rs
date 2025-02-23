use crate::{
    config::config::AppConfig, errors::errors::AppError, models::article::Article,
    services::articles::ArticlesClient,
};
use axum::{
    extract::{Extension, Query},
    http::{header, StatusCode},
    response::Response,
};
use chrono::{DateTime, Utc};
use rss::{ChannelBuilder, Guid, ItemBuilder};
use serde::Deserialize;

#[derive(Deserialize, Debug, Clone)]
pub struct RssQuery {
    pub flagged: Option<bool>,
    pub dead: Option<bool>,
    pub dupe: Option<bool>,
}

pub async fn generate_rss_feed<T: ArticlesClient + Clone>(
    Query(query): Query<RssQuery>,
    Extension(config): Extension<AppConfig>,
    Extension(client): Extension<T>,
) -> Result<Response<String>, AppError> {
    // Fetch articles and sort by ID descending
    let mut articles = client.fetch_articles(&query, &config).await?;
    articles.sort_by(|a, b| b.id.cmp(&a.id));

    // Build RSS items
    let items: Vec<_> = articles.iter().map(|article| build_item(article)).collect();

    // Build RSS channel
    let channel = ChannelBuilder::default()
        .title("GopherSignal RSS Feed")
        .link("https://gophersignal.com")
        .description("Latest articles from GopherSignal")
        .last_build_date(Utc::now().to_rfc2822())
        .items(items)
        .build();

    Ok(Response::builder()
        .status(StatusCode::OK)
        .header(header::CONTENT_TYPE, "application/rss+xml")
        .body(channel.to_string())?)
}

// Build RSS item with GUID and dates
fn build_item(article: &Article) -> rss::Item {
    let pub_date = DateTime::parse_from_rfc3339(
        &article
            .published_at
            .clone()
            .unwrap_or(article.created_at.clone()),
    )
    .unwrap_or_else(|_| Utc::now().into())
    .to_rfc2822();

    ItemBuilder::default()
        .title(Some(article.title.clone()))
        .link(Some(article.link.clone()))
        .description(Some(build_description(article)))
        .pub_date(Some(pub_date))
        .guid(Some(Guid {
            value: article.link.clone(),
            permalink: true,
        }))
        .build()
}

// Build description for the RSS item
fn build_description(article: &Article) -> String {
    format!(
        "<strong>Summary:</strong> {}<br><br>\
         <strong>Upvotes:</strong> {}<br><br>\
         <strong>Comments:</strong> {} [<a href=\"{}\">View Comments</a>]<br><br>\
         <strong>Link:</strong> <a href=\"{}\">{}</a>",
        article.summary.as_deref().unwrap_or("No summary"),
        article.upvotes.unwrap_or(0),
        article.comment_count.unwrap_or(0),
        article.comment_link.as_deref().unwrap_or("#"),
        article.link,
        article.title
    )
}
