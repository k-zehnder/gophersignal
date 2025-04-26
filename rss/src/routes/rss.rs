use crate::{
    config::config::AppConfig, errors::errors::AppError, models::article::Article,
    services::articles::ArticlesClient,
};
use axum::{
    extract::{Extension, Query},
    http::{header, StatusCode},
    response::Response,
};
use chrono::{DateTime, Duration, Utc};
use htmlescape::encode_minimal;
use rss::{ChannelBuilder, GuidBuilder, ItemBuilder};
use serde::Deserialize;
use url::Url;

// Filters for RSS feed.
#[derive(Deserialize, Debug, Clone)]
pub struct RssQuery {
    pub flagged: Option<bool>,
    pub dead: Option<bool>,
    pub dupe: Option<bool>,
    pub min_upvotes: Option<u32>,
    pub min_comments: Option<u32>,
}

/// Generate the RSS feed based on query filters.
pub async fn generate_rss_feed<T>(
    Query(query): Query<RssQuery>,
    Extension(config): Extension<AppConfig>,
    Extension(client): Extension<T>,
) -> Result<Response<String>, AppError>
where
    T: ArticlesClient + Clone,
{
    let articles = client.fetch_articles(&query, &config).await?;
    let channel = build_rss_channel(articles, &query);
    build_rss_response(channel)
}

// Sort by ID descending to have newest first.
// idx 0 has highest ID giving newest pubDate.
fn build_rss_channel(mut articles: Vec<Article>, query: &RssQuery) -> rss::Channel {
    articles.sort_by(|a, b| b.id.cmp(&a.id));
    let total = articles.len();

    let items: Vec<_> = articles
        .iter()
        .enumerate()
        .map(|(i, art)| build_rss_item(art, total, i))
        .collect();

    ChannelBuilder::default()
        .title(build_feed_title(query))
        .link("https://gophersignal.com")
        .description("Latest articles from Gopher Signal")
        .last_build_date(Utc::now().to_rfc2822())
        .items(items)
        .build()
}

// Build RSS title from active filters.
fn build_feed_title(query: &RssQuery) -> String {
    let mut parts = Vec::new();

    if query.flagged.unwrap_or(false) {
        parts.push("Flagged");
    }
    if query.dead.unwrap_or(false) {
        parts.push("Dead");
    }
    if query.dupe.unwrap_or(false) {
        parts.push("Dupe");
    }
    if query.min_upvotes.unwrap_or(0) > 0 || query.min_comments.unwrap_or(0) > 0 {
        parts.push("Filtered");
    }

    if parts.is_empty() {
        "Gopher Signal".into()
    } else {
        format!("Gopher Signal - {}", parts.join(", "))
    }
}

// Convert article to RSS item with staggered pubDate.
fn build_rss_item(article: &Article, total: usize, idx: usize) -> rss::Item {
    let offset = (total - 1 - idx) as i64;
    ItemBuilder::default()
        .title(Some(article.title.clone()))
        .description(Some(build_item_description(article)))
        .pub_date(Some(format_pub_date(&article.created_at, offset)))
        .guid(Some(build_item_guid(article)))
        .build()
}

// Escape summary and append footer.
fn build_item_description(article: &Article) -> String {
    let sum = article.summary.as_deref().unwrap_or("No summary");
    format!(
        "{}<br><br><small>{}</small>",
        encode_minimal(sum),
        build_item_footer(article)
    )
}

// Format creation date with offset as RFC2822.
fn format_pub_date(ts: &str, off: i64) -> String {
    let base = DateTime::parse_from_rfc3339(ts)
        .unwrap_or_else(|_| Utc::now().into())
        .with_timezone(&Utc);
    let dt = base
        .checked_add_signed(Duration::seconds(off))
        .unwrap_or(base);
    dt.to_rfc2822()
}

// Convert article to a <guid>.
fn build_item_guid(article: &Article) -> rss::Guid {
    let (value, is_permalink) = match article.hn_id {
        Some(id) if id > 0 => (format!("https://news.ycombinator.com/item?id={}", id), true),
        _ => (article.link.clone(), false),
    };
    GuidBuilder::default()
        .value(value)
        .permalink(is_permalink)
        .build()
}

// Footer shows upvotes comments link and source domain.
fn build_item_footer(article: &Article) -> String {
    let upvotes = article.upvotes.unwrap_or(0);
    let count = article.comment_count.unwrap_or(0);
    let clink = article.comment_link.as_deref().unwrap_or("#");
    let comments = format!("💬 <a href=\"{}\">{}</a>", encode_minimal(clink), count);
    let domain = Url::parse(&article.link)
        .ok()
        .and_then(|u| u.host_str().map(str::to_string))
        .unwrap_or_else(|| "source".into());
    format!(
        "▲ {} · {} · via <a href=\"{}\">{}</a>",
        upvotes,
        comments,
        encode_minimal(&article.link),
        encode_minimal(&domain)
    )
}

// Wrap channel in HTTP response.
fn build_rss_response(channel: rss::Channel) -> Result<Response<String>, AppError> {
    Response::builder()
        .status(StatusCode::OK)
        .header(header::CONTENT_TYPE, "application/rss+xml; charset=utf-8")
        .body(channel.to_string())
        .map_err(Into::into)
}
