export * from './article';
export * from './config';
export * from './services';
export * from './db';

import { z } from 'zod';

export const SummaryResponseSchema = z.object({
  summary: z.string().optional(),
  response: z.string().optional(),
  _meta: z.any().optional(),
});
