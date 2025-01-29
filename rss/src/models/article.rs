use serde::Deserialize;

#[derive(Deserialize)]
pub struct Article {
    pub id: u32,
    pub title: String,
    pub link: String,
    pub summary: String,
    pub created_at: String,
    pub upvotes: u32,
    pub comment_count: u32,
    pub comment_link: String,
}
