import { z } from 'zod';

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
  commitHash: string;
  modelName: string;
}

export const SummaryResponseSchema = z.object({
  summary: z
    .preprocess((raw) => {
      if (Array.isArray(raw)) {
        // Join array of lines into one string
        return (raw as string[]).join('\n');
      }
      return raw;
    }, z.string())
    .optional(),
  _meta: z.any().optional(),
});
