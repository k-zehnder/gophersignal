import { z } from 'zod';

export const ArticleSchema = z.object({
  title: z.string(),
  link: z.string().url(),
  article_rank: z.number().int(),
  flagged: z.boolean(),
  dead: z.boolean(),
  dupe: z.boolean(),
  upvotes: z.number().int(),
  comment_count: z.number().int(),
  comment_link: z.string().url(),
  content: z.string().optional(),
  summary: z.string().optional(),
  category: z.string().optional(),
});

export type Article = z.infer<typeof ArticleSchema>;
