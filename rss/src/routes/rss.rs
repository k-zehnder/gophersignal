use crate::services::fetch_articles::fetch_articles;
use axum::{http::StatusCode, response::IntoResponse};
use chrono::{Duration, Utc};
use rss::{ChannelBuilder, ItemBuilder};

pub async fn generate_rss_feed() -> Result<impl IntoResponse, StatusCode> {
    // Fetch articles from external API
    let mut articles = fetch_articles().await.map_err(|err| {
        eprintln!("Failed to fetch articles: {}", err);
        StatusCode::INTERNAL_SERVER_ERROR
    })?;

    // Sort articles by `id` DESC
    articles.sort_by(|a, b| b.id.cmp(&a.id));

    // Assign unique pubDates in descending order
    let now = Utc::now();
    let items: Vec<_> = articles
        .into_iter()
        .enumerate()
        .map(|(i, article)| {
            let pub_date = (now - Duration::minutes(i as i64)).to_rfc2822();

            let description = format!(
                "Summary: {}<br><br>Upvotes: {}<br><br>Comments: {} [<a href=\"{}\">View Comments</a>]<br><br>\
                 Link: <a href=\"{}\">{}</a>",
                article.summary,
                article.upvotes,
                article.comment_count,
                article.comment_link,
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

    // Build RSS feed
    let channel = ChannelBuilder::default()
        .title("GopherSignal RSS Feed")
        .link("https://gophersignal.com")
        .description("Latest articles from GopherSignal")
        .items(items)
        .build();

    Ok(axum::response::Html(channel.to_string()))
}
