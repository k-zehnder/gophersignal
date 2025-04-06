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
use url::Url;

#[derive(Deserialize, Debug, Clone)]
pub struct RssQuery {
    pub flagged: Option<bool>,
    pub dead: Option<bool>,
    pub dupe: Option<bool>,
    // Optional threshold parameters:
    pub min_upvotes: Option<u32>,
    pub min_comments: Option<u32>,
}

/// Generate the RSS feed.
pub async fn generate_rss_feed<T: ArticlesClient + Clone>(
    Query(query): Query<RssQuery>,
    Extension(config): Extension<AppConfig>,
    Extension(client): Extension<T>,
) -> Result<Response<String>, AppError> {
    // Fetch articles from the backend API using the query.
    let mut articles = client.fetch_articles(&query, &config).await?;
    // Sort articles in descending order by id.
    articles.sort_by(|a, b| b.id.cmp(&a.id));

    // Build RSS items.
    let items: Vec<_> = articles.iter().map(build_item).collect();

    // Construct the RSS channel.
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
    // Compute a unique publication date using the article ID as an offset.
    let id_offset = chrono::Duration::seconds(article.id);
    let pub_date =
        DateTime::parse_from_rfc3339(&article.published_at.as_ref().unwrap_or(&article.created_at))
            .unwrap_or_else(|_| Utc::now().into())
            .checked_add_signed(id_offset)
            .unwrap()
            .to_rfc2822();

    // Parse the domain from the article link.
    let domain = Url::parse(&article.link)
        .ok()
        .and_then(|url| url.host_str().map(|h| h.to_string()))
        .unwrap_or_else(|| "source".to_string());

    let summary = encode_minimal(article.summary.as_deref().unwrap_or("No summary"));

    // Render comment text.
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

    // Build an info string with upvotes, comment text, and source.
    let upvotes = article.upvotes.unwrap_or(0);
    let info = vec![
        format!("▲ {}", upvotes),
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
