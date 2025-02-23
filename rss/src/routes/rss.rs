use crate::db::db::{load_published_articles, update_published_articles};
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
use rss::{ChannelBuilder, Guid, ItemBuilder};
use serde::Deserialize;
use std::collections::HashSet;

#[derive(Deserialize, Debug, Clone)]
pub struct RssQuery {
    pub flagged: Option<bool>,
    pub dead: Option<bool>,
    pub dupe: Option<bool>,
}

/// Generate the RSS feed by including all articles but indicating new ones.
pub async fn generate_rss_feed<T: ArticlesClient + Clone>(
    Query(query): Query<RssQuery>,
    Extension(config): Extension<AppConfig>,
    Extension(client): Extension<T>,
) -> Result<Response<String>, AppError> {
    // Fetch articles and sort by ID descending.
    let mut articles = client.fetch_articles(&query, &config).await?;
    articles.sort_by(|a, b| b.id.cmp(&a.id));

    // Deduplicate articles using article.link as a unique key.
    let mut seen = HashSet::new();
    articles.retain(|article| {
        if seen.contains(&article.link) {
            false
        } else {
            seen.insert(article.link.clone());
            true
        }
    });

    // Load the set of previously published articles.
    let mut published = load_published_articles();

    // Build RSS items for every article.
    let items: Vec<_> = articles
        .iter()
        .map(|article| {
            // If the article link is not in our published set, it's new.
            let is_new = !published.contains(&article.link);
            if is_new {
                published.insert(article.link.clone());
            }
            build_item(article, is_new)
        })
        .collect();

    // Persist the updated published articles set.
    update_published_articles(&published);

    // Build the RSS channel.
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

/// Build an RSS item with GUID, publication date, and a custom <isNew> element.
fn build_item(article: &Article, is_new: bool) -> rss::Item {
    let pub_date = DateTime::parse_from_rfc3339(
        &article
            .published_at
            .clone()
            .unwrap_or(article.created_at.clone()),
    )
    .unwrap_or_else(|_| Utc::now().into())
    .to_rfc2822();

    let description = format!(
        "<isNew>{}</isNew>\
         <strong>Summary:</strong> {}<br><br>\
         <strong>Upvotes:</strong> {}<br><br>\
         <strong>Comments:</strong> {} [<a href=\"{}\">View Comments</a>]<br><br>\
         <strong>Link:</strong> <a href=\"{}\">{}</a>",
        if is_new { "true" } else { "false" },
        article.summary.as_deref().unwrap_or("No summary"),
        article.upvotes.unwrap_or(0),
        article.comment_count.unwrap_or(0),
        article.comment_link.as_deref().unwrap_or("#"),
        article.link,
        article.title
    );

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
