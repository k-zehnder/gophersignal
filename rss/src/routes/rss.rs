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
use sha1::{Digest, Sha1};
use url::Url;

#[derive(Deserialize, Debug, Clone)]
pub struct RssQuery {
    pub flagged: Option<bool>,
    pub dead: Option<bool>,
    pub dupe: Option<bool>,
    pub min_upvotes: Option<u32>,
    pub min_comments: Option<u32>,
}

/// Generates the RSS feed title based on the query filters.
fn generate_title(query: &RssQuery) -> String {
    let mut parts = Vec::with_capacity(4);

    // Boolean filters
    [query.flagged, query.dead, query.dupe]
        .iter()
        .zip(["Flagged", "Dead", "Dupe"])
        .filter_map(|(opt, label)| opt.filter(|&flag| flag).map(|_| label))
        .for_each(|lbl| parts.push(lbl));

    // Threshold filters
    if [query.min_upvotes, query.min_comments]
        .iter()
        .any(|&opt| opt.filter(|&v| v > 0).is_some())
    {
        parts.push("Filtered");
    }

    if parts.is_empty() {
        "Gopher Signal".into()
    } else {
        format!("Gopher Signal – {}", parts.join(", "))
    }
}

/// Builds an RSS item with a GUID that never changes.
fn build_item(article: &Article) -> rss::Item {
    // Stable GUID: SHA-1 of canonical external link (host + path)
    let canonical = Url::parse(&article.link)
        .map(|u| format!("{}{}", u.host_str().unwrap_or("").to_lowercase(), u.path()))
        .unwrap_or_else(|_| article.link.to_lowercase());
    let guid_value = format!("sha1:{:x}", Sha1::digest(canonical.as_bytes()));
    let is_permalink = false;

    // Click target: HN thread if present, else the article URL
    let link_target = article
        .comment_link
        .as_deref()
        .unwrap_or(&article.link)
        .to_string();

    // Metadata
    let domain = Url::parse(&link_target)
        .ok()
        .and_then(|u| u.host_str().map(ToString::to_string))
        .unwrap_or_else(|| "source".into());
    let summary = article.summary.as_deref().unwrap_or("No summary");
    let pub_date = compute_pub_date(article);

    ItemBuilder::default()
        .title(Some(article.title.clone()))
        .description(Some(format!(
            "{}<br><br><small>{}</small>",
            encode_minimal(summary),
            build_info(article, &domain)
        )))
        .pub_date(Some(pub_date))
        .guid(Some(Guid {
            value: guid_value,
            permalink: is_permalink,
        }))
        .build()
}

/// Compute RFC-2822 pubDate by adding article.id seconds to created_at, preventing identical timestamps.
fn compute_pub_date(article: &Article) -> String {
    let base = DateTime::parse_from_rfc3339(&article.created_at)
        .unwrap_or_else(|_| Utc::now().into())
        .with_timezone(&Utc);
    let offset = chrono::Duration::seconds(article.id as i64);
    base.checked_add_signed(offset).unwrap_or(base).to_rfc2822()
}

/// Builds the footer: "▲ upvotes · 💬 comments · via <domain>"
fn build_info(article: &Article, domain: &str) -> String {
    let comments = match article.comment_count.unwrap_or(0) {
        0 => "💬 0 comments".into(),
        n => format!(
            "💬 <a href=\"{}\">{}</a>",
            encode_minimal(article.comment_link.as_deref().unwrap_or("#")),
            n
        ),
    };
    let link = encode_minimal(&article.link);
    let domain_escaped = encode_minimal(domain);

    format!(
        "▲ {} · {} · via <a href=\"{}\">{}</a>",
        article.upvotes.unwrap_or(0),
        comments,
        link,
        domain_escaped
    )
}

/// `GET /rss` – fetch articles, build channel, return RSS XML.
pub async fn generate_rss_feed<T: ArticlesClient + Clone>(
    Query(query): Query<RssQuery>,
    Extension(cfg): Extension<AppConfig>,
    Extension(client): Extension<T>,
) -> Result<Response<String>, AppError> {
    let articles = client.fetch_articles(&query, &cfg).await?;

    let channel = ChannelBuilder::default()
        .title(generate_title(&query))
        .link("https://gophersignal.com")
        .description("Latest articles from Gopher Signal")
        .last_build_date(Utc::now().to_rfc2822())
        .items(articles.iter().map(build_item).collect::<Vec<_>>())
        .build();

    Ok(Response::builder()
        .status(StatusCode::OK)
        .header(header::CONTENT_TYPE, "application/rss+xml")
        .body(channel.to_string())?)
}
