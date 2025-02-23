use serde::Deserialize;

#[derive(Deserialize, Debug, Clone, sqlx::FromRow)]
pub struct Article {
    pub id: i64,
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
            "published_at": null
        }
        "#;
        let article: Article = serde_json::from_str(json).unwrap();
        assert_eq!(article.id, 1);
        assert_eq!(article.title, "Test Article");
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
                "published_at": null
            }]
        }
        "#;
        let response: ApiResponse = serde_json::from_str(json).unwrap();
        assert_eq!(response.code, 200);
        assert_eq!(response.status, "OK");
        let articles = response.articles.unwrap();
        assert_eq!(articles.len(), 1);
        assert_eq!(articles[0].id, 1);
    }
}
