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

/// Builds an RSS item from an article, including title, description, etc.
fn build_item(article: &Article) -> rss::Item {
    // HN link yields hn<id>; else use SHA-1 of the lower-cased title for stable GUID.
    let (guid_value, is_permalink) = extract_hn_guid(&article.link).unwrap_or_else(|| {
        let hash = format!(
            "{:x}",
            Sha1::digest(article.title.to_lowercase().as_bytes())
        );
        (format!("title:{hash}"), false)
    });

    let domain = extract_domain(&article.link);
    let summary = article.summary.as_deref().unwrap_or("No summary");
    let pub_date = compute_pub_date(article);

    ItemBuilder::default()
        .title(Some(article.title.clone()))
        .link(Some(article.link.clone()))
        .description(Some(format!(
            "{}<br><br><small>{}</small>",
            encode_minimal(summary),
            build_info(article, &domain),
        )))
        .pub_date(Some(pub_date))
        .guid(Some(Guid {
            value: guid_value,
            permalink: is_permalink,
        }))
        .build()
}

/// Returns RFC-2822 `<pubDate>` = parse(created_at) + article.id seconds
fn compute_pub_date(article: &Article) -> String {
    let base = DateTime::parse_from_rfc3339(&article.created_at)
        .unwrap_or_else(|_| Utc::now().into())
        .with_timezone(&Utc);

    let offset = chrono::Duration::seconds(article.id as i64);
    base.checked_add_signed(offset).unwrap_or(base).to_rfc2822()
}

/// If link is a Hacker News discussion, returns (hn<id>, false)
fn extract_hn_guid(link: &str) -> Option<(String, bool)> {
    Url::parse(link).ok().and_then(|url| {
        (url.host_str() == Some("news.ycombinator.com"))
            .then(|| {
                url.query_pairs()
                    .find(|(k, _)| k == "id")
                    .map(|(_, v)| (format!("hn:{v}"), false))
            })
            .flatten()
    })
}

/// Extracts the domain from the article's link.
fn extract_domain(link: &str) -> String {
    Url::parse(link)
        .ok()
        .and_then(|u| u.host_str().map(ToString::to_string))
        .unwrap_or_else(|| "source".into())
}

/// Builds the footer shown in the description.
fn build_info(article: &Article, domain: &str) -> String {
    let comments = match article.comment_count.unwrap_or(0) {
        0 => "💬 0 comments".into(),
        n => format!(
            "💬 <a href=\"{}\">{}</a>",
            encode_minimal(article.comment_link.as_deref().unwrap_or("#")),
            n
        ),
    };

    [
        format!("▲ {}", article.upvotes.unwrap_or(0)),
        comments,
        format!(
            "via <a href=\"{}\">{}</a>",
            encode_minimal(&article.link),
            encode_minimal(domain),
        ),
    ]
    .join(" · ")
}

/// `GET /rss` Fetch articles, build channel, return RSS XML.
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
