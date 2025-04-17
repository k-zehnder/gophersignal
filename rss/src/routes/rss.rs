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
    pub min_upvotes: Option<u32>,
    pub min_comments: Option<u32>,
}

/// Generate the RSS feed.
pub async fn generate_rss_feed<T: ArticlesClient + Clone>(
    Query(query): Query<RssQuery>,
    Extension(config): Extension<AppConfig>,
    Extension(client): Extension<T>,
) -> Result<Response<String>, AppError> {
    let mut articles = client.fetch_articles(&query, &config).await?;
    articles.sort_by(|a, b| b.id.cmp(&a.id));

    // Build title components.
    let mut components = Vec::new();

    // Add boolean filters.
    if query.flagged == Some(true) {
        components.push("Flagged");
    }
    if query.dead == Some(true) {
        components.push("Dead");
    }
    if query.dupe == Some(true) {
        components.push("Dupe");
    }

    // Add threshold filters.
    let has_thresholds = query.min_upvotes.filter(|&v| v > 0).is_some()
        || query.min_comments.filter(|&c| c > 0).is_some();

    if has_thresholds {
        components.push("Filtered");
    }

    // Construct final title.
    let title = match components.is_empty() {
        true => "Gopher Signal".to_string(),
        false => format!("Gopher Signal - {}", components.join(", ")),
    };

    let items: Vec<_> = articles.iter().map(build_item).collect();

    let channel = ChannelBuilder::default()
        .title(title)
        .link("https://gophersignal.com")
        .description("Latest articles from Gopher Signal")
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
    let id_offset = chrono::Duration::seconds(article.id as i64);
    let pub_date =
        DateTime::parse_from_rfc3339(&article.published_at.as_ref().unwrap_or(&article.created_at))
            .unwrap_or_else(|_| Utc::now().into())
            .checked_add_signed(id_offset)
            .unwrap()
            .to_rfc2822();

    // Extract HN guid or default link
    let guid_value = Url::parse(&article.link)
        .ok()
        .and_then(|url| {
            if url.host_str()? == "news.ycombinator.com" {
                url.query_pairs()
                    .find(|(k, _)| k == "id")
                    .map(|(_, v)| v.to_string())
            } else {
                None
            }
        })
        .unwrap_or_else(|| article.link.clone());

    // Parse domain for source
    let domain = Url::parse(&article.link)
        .ok()
        .and_then(|url| url.host_str().map(|h| h.to_string()))
        .unwrap_or_else(|| "source".to_string());

    let summary = encode_minimal(article.summary.as_deref().unwrap_or("No summary"));
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
            value: guid_value,
            permalink: false,
        }))
        .build()
}
