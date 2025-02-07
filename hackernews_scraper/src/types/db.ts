import { Connection } from 'mysql2/promise';
import { Article } from './article';

export interface DBClient {
  saveArticles: (articles: Article[]) => Promise<void>;
  updateArticleSummary: (id: number, summary: string) => Promise<void>;
  markArticleAsDead: (id: number) => Promise<void>;
  markArticleAsDuplicate: (id: number) => Promise<void>;
  closeDatabaseConnection: () => Promise<void>;
  connection: Connection;
}
