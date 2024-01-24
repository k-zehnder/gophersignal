export interface ApiArticle {
    id: number;
    title: string;
    source: string;
    created_at: string;
    updated_at: string;
    summary: { String: string; Valid: boolean } | null;
    link: string;
    is_on_homepage: boolean;
}
  
export interface Article {
    id: number;
    title: string;
    source: string;
    createdAt: string;
    updatedAt: string;
    summary: string;
    link: string;
    isOnHomepage: boolean;
}
