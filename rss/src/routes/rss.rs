//! RSS feed generator for GopherSignal using in-memory deduplication

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
use htmlescape::encode_minimal;
use rss::{ChannelBuilder, Guid, ItemBuilder};
use serde::Deserialize;
use std::collections::HashSet;
use url::Url;

#[derive(Deserialize, Debug, Clone)]
pub struct RssQuery {
    pub flagged: Option<bool>,
    pub dead: Option<bool>,
    pub dupe: Option<bool>,
}

/// Generate the RSS feed
pub async fn generate_rss_feed<T: ArticlesClient + Clone>(
    Query(query): Query<RssQuery>,
    Extension(config): Extension<AppConfig>,
    Extension(client): Extension<T>,
) -> Result<Response<String>, AppError> {
    // Fetch articles and sort in descending order by ID
    let mut articles = client.fetch_articles(&query, &config).await?;
    articles.sort_by(|a, b| b.id.cmp(&a.id));

    // Deduplicate articles using a HashSet
    let mut seen = HashSet::new();
    articles.retain(|article| seen.insert(article.link.clone()));

    // Build RSS items
    let items: Vec<_> = articles.iter().map(build_item).collect();

    // Construct RSS channel
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

fn build_item(article: &Article) -> rss::Item {
    // Compute a unique publication date using article ID as an offset
    let id_offset = chrono::Duration::seconds(article.id as i64);
    let pub_date =
        DateTime::parse_from_rfc3339(&article.published_at.as_ref().unwrap_or(&article.created_at))
            .unwrap_or_else(|_| Utc::now().into())
            .checked_add_signed(id_offset)
            .unwrap()
            .to_rfc2822();

    // Get the domain from the article link
    let domain = Url::parse(&article.link)
        .ok()
        .and_then(|url| url.host_str().map(|h| h.to_string()))
        .unwrap_or_else(|| "source".to_string());

    let summary = encode_minimal(article.summary.as_deref().unwrap_or("No summary"));

    // Render comment text. Clickable link if comments exist; plain text if 0
    let comment_count = article.comment_count.unwrap_or(0);
    let comment_text = if comment_count > 0 {
        format!(
            "💬 <a href=\"{}\">{}</a>",
            encode_minimal(article.comment_link.as_deref().unwrap_or("#")),
            comment_count
        )
    } else {
        "💬 0 comments".to_string()
    };

    // Build info string with upvotes, comment text, and source info
    let info = vec![
        format!("▲ {}", article.upvotes.unwrap_or(0)),
        comment_text,
        format!(
            "via <a href=\"{}\">{}</a>",
            encode_minimal(&article.link),
            encode_minimal(&domain)
        ),
    ]
    .join(" · ");

    let description = format!("{}<br><br><small>{}</small>", summary, info);

    ItemBuilder::default()
        .title(Some(article.title.clone()))
        .link(Some(article.link.clone()))
        .description(Some(description))
        .pub_date(Some(pub_date))
        .guid(Some(Guid {
            value: article.link.clone(),
            permalink: true,
        }))
        .build()
}
