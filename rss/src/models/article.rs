use serde::Deserialize;

#[derive(Deserialize, Debug)]
pub struct Article {
    pub id: u32,
    pub title: String,
    pub link: String,
    pub content: Option<String>,
    pub summary: Option<String>,
    pub source: String,
    pub created_at: String,
    pub updated_at: String,
    pub upvotes: Option<u32>,
    pub comment_count: Option<u32>,
    pub comment_link: Option<String>,
    pub flagged: Option<bool>,
    pub dead: Option<bool>,
    pub dupe: Option<bool>,
}

#[derive(Deserialize, Debug)]
pub struct ApiResponse {
    pub code: u32,
    pub status: String,
    pub total_count: u32,
    pub articles: Option<Vec<Article>>,
}
