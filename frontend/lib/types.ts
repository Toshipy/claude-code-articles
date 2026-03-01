export type Source = {
  id: string;
  name: string;
  url?: string;
  icon_url: string;
  article_count?: number;
  last_collected_at?: string;
};

export type Tag = {
  id: string;
  name: string;
  slug: string;
  article_count?: number;
};

export type ArticleSummary = {
  id: string;
  title: string;
  summary: string;
  url: string;
  thumbnail_url: string;
  published_at: string;
  collected_at?: string;
  source: Source;
  tags: Tag[];
  highlight?: string;
};

export type RelatedArticle = {
  id: string;
  title: string;
  thumbnail_url: string;
};

export type ArticleDetail = ArticleSummary & {
  content: string;
  author: string;
  related_articles: RelatedArticle[];
};

export type User = {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
  role: "user" | "admin";
  bookmark_count?: number;
  created_at: string;
  updated_at?: string;
};

export type Bookmark = {
  article_id: string;
  bookmarked_at: string;
  article: ArticleSummary;
};

export type Pagination = {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
};

export type ApiResponse<T> = {
  success: true;
  data: T;
};

export type PaginatedResponse<T> = {
  success: true;
  data: T[];
  pagination: Pagination;
};

export type ApiError = {
  success: false;
  error: {
    code: string;
    message: string;
    details?: { field: string; message: string }[];
  };
};
