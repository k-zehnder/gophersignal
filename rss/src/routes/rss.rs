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
    let mut articles = client.fetch_articles(&query, &config).await?;
    // Sort newest first
    articles.sort_by(|a, b| b.id.cmp(&a.id));

    let title = build_title(&query);

    // Pass the zero-based index into build_item so we can offset by seconds reliably.
    let items: Vec<_> = articles
        .iter()
        .enumerate()
        .map(|(idx, article)| build_item(article, idx))
        .collect();

    let channel = ChannelBuilder::default()
        .title(title)
        .link("https://gophersignal.com")
        .description("Latest articles from Gopher Signal")
        .last_build_date(Utc::now().to_rfc2822())
        .items(items)
        .build();

    Ok(Response::builder()
        .status(StatusCode::OK)
        .header(header::CONTENT_TYPE, "application/rss+xml; charset=utf-8")
        .body(channel.to_string())?)
}

/// Build the RSS feed title from query filters.
fn build_title(query: &RssQuery) -> String {
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

/// Build a single RSS from an Article, using its index for the pubDate offset.
fn build_item(article: &Article, index: usize) -> rss::Item {
    ItemBuilder::default()
        .title(Some(article.title.clone()))
        .description(Some(format!(
            "{}<br><br><small>{}</small>",
            encode_minimal(article.summary.as_deref().unwrap_or("No summary")),
            build_footer(article)
        )))
        .pub_date(Some(format_pub_date(&article.created_at, index)))
        .guid(Some(build_guid(article)))
        .build()
}

/// Format pubDate by adding seconds to `created_at`.
fn format_pub_date(created_at: &str, index: usize) -> String {
    let base: DateTime<Utc> = DateTime::parse_from_rfc3339(created_at)
        .unwrap_or_else(|_| Utc::now().into())
        .with_timezone(&Utc);

    let dt = base
        .checked_add_signed(Duration::seconds(index as i64))
        .unwrap_or(base);

    dt.to_rfc2822()
}

/// Build a stable GUID
fn build_guid(article: &Article) -> rss::Guid {
    GuidBuilder::default()
        .value(match article.hn_id {
            Some(hn) => format!("https://news.ycombinator.com/item?id={}", hn),
            None => article.link.clone(),
        })
        .permalink(article.hn_id.is_some())
        .build()
}

/// Build the footer HTML of upvotes, comments, via domain.
fn build_footer(article: &Article) -> String {
    let up = article.upvotes.unwrap_or(0);
    let comments = article.comment_count.unwrap_or(0);
    let comment_html = if comments > 0 {
        format!(
            "💬 <a href=\"{}\">{}</a>",
            encode_minimal(article.comment_link.as_deref().unwrap_or("#")),
            comments
        )
    } else {
        "💬 0 comments".into()
    };
    let domain = Url::parse(&article.link)
        .ok()
        .and_then(|u| u.host_str().map(str::to_string))
        .unwrap_or_else(|| "source".into());
    format!(
        "▲ {} · {} · via <a href=\"{}\">{}</a>",
        up,
        comment_html,
        encode_minimal(&article.link),
        encode_minimal(&domain)
    )
}
