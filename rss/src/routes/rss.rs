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

/// Generates the RSS feed title based on the query filters.
fn generate_title(query: &RssQuery) -> String {
    let mut components = Vec::with_capacity(4);

    // Boolean filters
    [query.flagged, query.dead, query.dupe]
        .iter()
        .zip(["Flagged", "Dead", "Dupe"])
        .filter_map(|(flag, label)| flag.and_then(|f| f.then_some(label)))
        .for_each(|label| components.push(label));

    // Threshold filters
    let has_thresholds = [query.min_upvotes, query.min_comments]
        .iter()
        .any(|&v| v.filter(|&x| x > 0).is_some());

    if has_thresholds {
        components.push("Filtered");
    }

    match components.is_empty() {
        true => "Gopher Signal".into(),
        false => format!("Gopher Signal - {}", components.join(", ")),
    }
}

/// Builds an RSS item from an article, including title, description, etc.
fn build_item(article: &Article) -> rss::Item {
    let (guid_value, is_permalink) = extract_hn_guid(&article.link);
    let domain = extract_domain(&article.link);
    let summary = article.summary.as_deref().unwrap_or("No summary");

    ItemBuilder::default()
        .title(Some(article.title.clone()))
        .link(Some(article.link.clone()))
        .description(Some(format!(
            "{}<br><br><small>{}</small>",
            encode_minimal(summary),
            build_info(article, &domain)
        )))
        .pub_date(Some(compute_pub_date(article)))
        .guid(Some(Guid {
            value: guid_value,
            permalink: is_permalink,
        }))
        .build()
}

/// Computes the publication date for an article, adjusting by article ID.
fn compute_pub_date(article: &Article) -> String {
    let base_date = article
        .published_at
        .as_deref()
        .unwrap_or(&article.created_at);
    let offset = chrono::Duration::seconds(article.id.into());

    DateTime::parse_from_rfc3339(base_date)
        .unwrap_or_else(|_| Utc::now().into())
        .checked_add_signed(offset)
        .unwrap_or_else(|| Utc::now().into())
        .to_rfc2822()
}

/// Extracts the Hacker News GUID from a link if it points to a Hacker News article.
fn extract_hn_guid(link: &str) -> (String, bool) {
    Url::parse(link)
        .ok()
        .and_then(|url| {
            if url.host_str() != Some("news.ycombinator.com") {
                return None;
            }

            url.query_pairs()
                .find(|(k, _)| k == "id")
                .map(|(_, v)| (format!("hn_id={}", v), false))
        })
        .unwrap_or_else(|| (link.to_string(), false))
}

/// Extracts the domain from the article's link.
fn extract_domain(link: &str) -> String {
    Url::parse(link)
        .ok()
        .and_then(|url| url.host_str().map(ToString::to_string))
        .unwrap_or_else(|| "source".into())
}

/// Builds additional information for an RSS item, including upvotes, comments, and source.
fn build_info(article: &Article, domain: &str) -> String {
    let comment_link = article.comment_link.as_deref().unwrap_or("#");
    let comment_count = article.comment_count.unwrap_or(0);

    let comment_text = if comment_count == 0 {
        "💬 0 comments".into()
    } else {
        format!(
            "💬 <a href=\"{}\">{}</a>",
            encode_minimal(comment_link),
            comment_count
        )
    };

    [
        format!("▲ {}", article.upvotes.unwrap_or(0)),
        comment_text,
        format!(
            "via <a href=\"{}\">{}</a>",
            encode_minimal(&article.link),
            encode_minimal(domain)
        ),
    ]
    .join(" · ")
}

/// Main function to generate the RSS feed using articles fetched via `ArticlesClient`.
pub async fn generate_rss_feed<T: ArticlesClient + Clone>(
    Query(query): Query<RssQuery>,
    Extension(config): Extension<AppConfig>,
    Extension(client): Extension<T>,
) -> Result<Response<String>, AppError> {
    let mut articles = client.fetch_articles(&query, &config).await?;
    articles.sort_by_key(|a| std::cmp::Reverse(a.id));

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
