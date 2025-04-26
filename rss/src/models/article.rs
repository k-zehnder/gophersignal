use serde::Deserialize;

#[derive(Deserialize, Debug, Clone, sqlx::FromRow)]
pub struct Article {
    pub id: i32,
    pub hn_id: Option<i32>,
    pub title: String,
    pub link: String,
    #[serde(default)]
    pub article_rank: i32,
    pub content: Option<String>,
    pub summary: Option<String>,
    pub source: String,
    pub upvotes: Option<i32>,
    pub comment_count: Option<i32>,
    pub comment_link: Option<String>,
    pub flagged: bool,
    pub dead: bool,
    pub dupe: bool,
    pub created_at: String,
    pub updated_at: String,
    pub published_at: Option<String>,

    // New metadata fields
    pub commit_hash: Option<String>,
    pub model_name: Option<String>,
}

#[derive(Deserialize, Debug)]
pub struct ApiResponse {
    pub code: u32,
    pub status: String,
    pub total_count: u32,
    pub articles: Option<Vec<Article>>,
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json;

    #[test]
    fn test_deserialize_article() {
        let json = r#"
        {
            "id": 1,
            "hn_id": 42,
            "title": "Test Article",
            "link": "http://example.com",
            "content": "Some content",
            "summary": "Test summary",
            "source": "Example",
            "created_at": "2022-01-01T00:00:00Z",
            "updated_at": "2022-01-01T00:00:00Z",
            "upvotes": 10,
            "comment_count": 2,
            "comment_link": "http://example.com/comments",
            "flagged": false,
            "dead": false,
            "dupe": false,
            "published_at": null,
            "commit_hash": "abc123",
            "model_name": "premium-model"
        }
        "#;
        let article: Article = serde_json::from_str(json).unwrap();
        assert_eq!(article.id, 1);
        assert_eq!(article.hn_id, Some(42));
        assert_eq!(article.article_rank, 0);
        assert_eq!(article.commit_hash.as_deref(), Some("abc123"));
        assert_eq!(article.model_name.as_deref(), Some("premium-model"));
    }

    #[test]
    fn test_deserialize_api_response() {
        let json = r#"
        {
            "code": 200,
            "status": "OK",
            "total_count": 1,
            "articles": [{
                "id": 1,
                "hn_id": 42,
                "title": "Test Article",
                "link": "http://example.com",
                "content": "Some content",
                "summary": "Test summary",
                "source": "Example",
                "created_at": "2022-01-01T00:00:00Z",
                "updated_at": "2022-01-01T00:00:00Z",
                "upvotes": 10,
                "comment_count": 2,
                "comment_link": "http://example.com/comments",
                "flagged": false,
                "dead": false,
                "dupe": false,
                "published_at": null,
                "commit_hash": "abc123",
                "model_name": "premium-model"
            }]
        }
        "#;
        let response: ApiResponse = serde_json::from_str(json).unwrap();
        let articles = response.articles.unwrap();
        assert_eq!(articles.len(), 1);
        assert_eq!(articles[0].commit_hash.as_deref(), Some("abc123"));
        assert_eq!(articles[0].model_name.as_deref(), Some("premium-model"));
    }
}
