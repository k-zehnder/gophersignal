use crate::config::AppConfig;
use crate::services::articles::fetch_articles;
use axum::{
    extract::{Extension, Query},
    http::StatusCode,
    response::{Html, IntoResponse},
};
use chrono::{Duration, Utc};
use rss::{ChannelBuilder, ItemBuilder};
use serde::Deserialize;

#[derive(Deserialize)]
pub struct RssQuery {
    pub flagged: Option<bool>,
    pub dead: Option<bool>,
    pub dupe: Option<bool>,
}

pub async fn generate_rss_feed(
    Query(query): Query<RssQuery>,
    Extension(config): Extension<AppConfig>,
) -> Result<impl IntoResponse, StatusCode> {
    // Fetch articles
    let mut articles = fetch_articles(&query, &config).await.map_err(|err| {
        eprintln!("Failed to fetch articles: {}", err);
        StatusCode::INTERNAL_SERVER_ERROR
    })?;

    // Sort articles by id descending
    articles.sort_by(|a, b| b.id.cmp(&a.id));

    let now = Utc::now();
    let items: Vec<_> = articles
        .into_iter()
        .enumerate()
        .take(30)
        .map(|(i, article)| {
            // Calculate a pub date by subtracting i minutes from now
            let pub_date = (now - Duration::minutes(i as i64)).to_rfc2822();

            // Ensure the summary is not empty
            let summary = match &article.summary {
                Some(s) if !s.trim().is_empty() => s.clone(),
                _ => "No summary available.".to_string(),
            };

            // Wrap the description in CDATA to preserve HTML formatting
            let description = format!(
                "<![CDATA[<strong>Summary:</strong> {}<br><br>\
<strong>Upvotes:</strong> {}<br><br>\
<strong>Comments:</strong> {} [<a href=\"{}\">View Comments</a>]<br><br>\
<strong>Link:</strong> <a href=\"{}\">{}</a>]]>",
                summary,
                article.upvotes.unwrap_or(0),
                article.comment_count.unwrap_or(0),
                article
                    .comment_link
                    .unwrap_or_else(|| "No comments".to_string()),
                article.link,
                article.title
            );

            ItemBuilder::default()
                .title(Some(article.title))
                .link(Some(article.link))
                .description(Some(description))
                .pub_date(Some(pub_date))
                .build()
        })
        .collect();

    let channel = ChannelBuilder::default()
        .title("GopherSignal RSS Feed")
        .link("https://gophersignal.com")
        .description("Latest articles from GopherSignal")
        .items(items)
        .build();

    let rss_feed = channel.to_string();

    Ok(Html(rss_feed))
}
