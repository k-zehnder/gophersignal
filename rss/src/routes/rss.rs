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

/// Builds an RSS item with immutable GUID.
fn build_item(article: &Article) -> rss::Item {
    // GUID is SHA-1(host + path) of the canonical article URL.
    let canonical = Url::parse(&article.link)
        .map(|u| format!("{}{}", u.host_str().unwrap_or("").to_lowercase(), u.path()))
        .unwrap_or_else(|_| article.link.to_lowercase());
    let guid_value = format!("sha1:{:x}", Sha1::digest(canonical.as_bytes()));

    // Optional HN discussion thread URL is used in the footer.
    let hn_url = article.comment_link.as_deref();

    // Misc metadata.
    let domain = Url::parse(&article.link)
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
            build_info(article, &domain, hn_url)
        )))
        .pub_date(Some(pub_date))
        .guid(Some(Guid {
            value: guid_value,
            permalink: false,
        }))
        .build()
}

/// Compute RFC-2822 pubDate (unique by adding article.id seconds).
fn compute_pub_date(article: &Article) -> String {
    let base = DateTime::parse_from_rfc3339(&article.created_at)
        .unwrap_or_else(|_| Utc::now().into())
        .with_timezone(&Utc);
    let offset = chrono::Duration::seconds(article.id as i64);
    base.checked_add_signed(offset).unwrap_or(base).to_rfc2822()
}

/// Footer: "▲ upvotes · 💬 comments · via <domain>".
fn build_info(article: &Article, domain: &str, hn_url: Option<&str>) -> String {
    let comments = match (article.comment_count.unwrap_or(0), hn_url) {
        (0, _) => "💬 0 comments".into(),
        (n, Some(url)) => format!("💬 <a href=\"{}\">{}</a>", encode_minimal(url), n),
        (n, None) => format!("💬 {n} comments"),
    };

    format!(
        "▲ {} · {} · via <a href=\"{}\">{}</a>",
        article.upvotes.unwrap_or(0),
        comments,
        encode_minimal(&article.link),
        encode_minimal(domain)
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
        .header(header::CONTENT_TYPE, "application/rss+xml; charset=utf-8")
        .body(channel.to_string())?)
}
