use super::article::Article;
use serde::Deserialize;

#[derive(Deserialize)]
pub struct ApiResponse {
    pub code: u32,
    pub status: String,
    pub total_count: u32,
    pub articles: Vec<Article>,
}
