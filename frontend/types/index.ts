import { z } from 'zod';

// Define the Zod schema for a single article
export const ArticleSchema = z.object({
  id: z.union([z.string(), z.number()]).transform((value) => value.toString()),
  title: z.string(),
  source: z.string().optional(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
  summary: z
    .object({
      String: z.string(),
      Valid: z.boolean(),
    })
    .transform((val) => (val.Valid ? val.String : '')),
  link: z.string().url(),
  upvotes: z
    .object({
      Int64: z.number(),
      Valid: z.boolean(),
    })
    .transform((val) => (val.Valid ? val.Int64 : 0)),
  comment_count: z
    .object({
      Int64: z.number(),
      Valid: z.boolean(),
    })
    .transform((val) => (val.Valid ? val.Int64 : 0)),
  comment_link: z
    .object({
      String: z.string(),
      Valid: z.boolean(),
    })
    .transform((val) => (val.Valid ? val.String : '')),
});

// Define the Zod schema for the API response
export const ArticlesResponseSchema = z.object({
  code: z.number(),
  status: z.string(),
  articles: z.array(ArticleSchema),
});

// Define TypeScript types from Zod schemas
export type Article = z.infer<typeof ArticleSchema>;
export type ArticlesResponse = z.infer<typeof ArticlesResponseSchema>;
