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

// Sort by ID descending so newest appears first.
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

// Convert an Article into an RSS <item>.
fn build_rss_item(article: &Article, total: usize, idx: usize) -> rss::Item {
    let offset = (total - 1 - idx) as i64;
    ItemBuilder::default()
        .title(Some(article.title.clone()))
        .description(Some(build_item_description(article)))
        .pub_date(Some(format_pub_date(&article.created_at, offset)))
        .guid(Some(build_item_guid(article)))
        .build()
}

// Replace newlines with <br>, escape the summary, and append the footer metadata.
fn build_item_description(article: &Article) -> String {
    // Get the raw summary, replace newlines with <br>, then escape.
    let summary_with_breaks = article
        .summary
        .as_deref()
        .unwrap_or("No summary")
        .replace('\n', "<br>");
    let mut html = encode_minimal(&summary_with_breaks);

    // Append the footer.
    html.push_str(&format!(
        "<br><br><small>{}</small>",
        build_item_footer(article)
    ));
    html
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

// Build a <guid> element for RSS.
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

// Footer metadata showing upvotes comments model commit domain.
fn build_item_footer(article: &Article) -> String {
    let upvotes_html = format!("▲ {}", article.upvotes.unwrap_or(0));

    // Comments count (link only if >0)
    let comments_html = {
        let cnt = article.comment_count.unwrap_or(0);
        let txt = cnt.to_string();
        if cnt > 0 {
            let url = article
                .comment_link
                .as_deref()
                .filter(|s| !s.is_empty())
                .map(ToString::to_string)
                .unwrap_or_else(|| {
                    article
                        .hn_id
                        .filter(|&id| id > 0)
                        .map(|id| format!("https://news.ycombinator.com/item?id={}", id))
                        .unwrap_or_else(|| "#".into())
                });
            format!(
                r#"<a href="{url}">💬 {txt}</a>"#,
                url = encode_minimal(&url),
                txt = txt
            )
        } else {
            format!("💬 {}", txt)
        }
    };

    let model_html = article
        .model_name
        .as_deref()
        .filter(|m| !m.is_empty())
        .map(|m| format!("🤖 {}", encode_minimal(m)))
        .unwrap_or_default();

    let commit_html = article
        .commit_hash
        .as_deref()
        .filter(|h| !h.is_empty())
        .map(|h| format!("🔨 {}", encode_minimal(h)))
        .unwrap_or_default();

    let domain_display = Url::parse(&article.link)
        .ok()
        .and_then(|u| u.host_str().map(str::to_string))
        .map(|host| {
            // Check case-insensitively and ensure host is longer than "www."
            // to prevent panic on slicing if host is exactly "www."
            if host.len() > 4 && host.to_lowercase().starts_with("www.") {
                host[4..].to_string() // Slice the original host string
            } else {
                host // Return the original host string if no "www." prefix or if host is just "www."
            }
        })
        .unwrap_or_else(|| "source".into());
    let domain_html = format!(
        r#"<a href="{href}">🌐 {display}</a>"#,
        href = encode_minimal(&article.link),
        display = encode_minimal(&domain_display)
    );

    vec![
        upvotes_html,
        comments_html,
        model_html,
        commit_html,
        domain_html,
    ]
    .into_iter()
    .filter(|s| !s.is_empty())
    .collect::<Vec<_>>()
    .join(" · ")
}

// Wrap channel XML in an HTTP response.
fn build_rss_response(channel: rss::Channel) -> Result<Response<String>, AppError> {
    Response::builder()
        .status(StatusCode::OK)
        .header(header::CONTENT_TYPE, "application/rss+xml; charset=utf-8")
        .body(channel.to_string())
        .map_err(Into::into)
}
