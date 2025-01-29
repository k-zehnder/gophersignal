use crate::services::fetch_articles::fetch_articles;
use axum::{http::StatusCode, response::IntoResponse};
use rss::{ChannelBuilder, ItemBuilder};

pub async fn generate_rss_feed() -> impl IntoResponse {
    // Fetch articles from external API
    let articles = match fetch_articles().await {
        Ok(data) => data,
        Err(_) => return StatusCode::INTERNAL_SERVER_ERROR.into_response(),
    };

    // Map articles to RSS items
    let items: Vec<_> = articles
        .into_iter()
        .map(|article| {
            let description = format!(
                "{}\n\nUpvotes: {}\nComments: {} [View Comments]({})",
                article.summary, article.upvotes, article.comment_count, article.comment_link
            );

            ItemBuilder::default()
                .title(Some(article.title))
                .link(Some(article.link))
                .description(Some(description))
                .pub_date(Some(article.created_at))
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

    axum::response::Html(channel.to_string()).into_response()
}
