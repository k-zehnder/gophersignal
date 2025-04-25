// Defines the Article interface used for Hacker News articles
export interface Article {
  title: string;
  link: string;
  hn_id: number;
  article_rank: number;
  flagged: boolean;
  dead: boolean;
  dupe: boolean;
  upvotes: number;
  comment_count: number;
  comment_link: string;
  content?: string;
  summary?: string;
  category?: string;
}
