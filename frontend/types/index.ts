// Represents a single article with all necessary details.
export interface Article {
  id: number;
  title: string;
  link: string;
  summary: string;
  source: string;
  createdAt: string;
  updatedAt: string;
}

// Represents the response structure for a list of articles.
export interface ArticlesResponse {
  code: number;
  status: string;
  articles: Article[];
}
