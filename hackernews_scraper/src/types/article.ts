// Defines the Article interface used for Hacker News articles
export interface Article {
  title: string;
  link: string;
  hnId: number;
  articleRank: number;
  flagged: boolean;
  dead: boolean;
  dupe: boolean;
  upvotes: number;
  commentCount: number;
  commentLink: string;
  content?: string;
  summary?: string;
  category?: string;
}
