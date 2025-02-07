export interface Article {
  title: string;
  link: string;
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
