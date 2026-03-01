import type {
  ApiResponse,
  PaginatedResponse,
  ArticleSummary,
  ArticleDetail,
  Tag,
  Source,
  Bookmark,
  User,
} from "./types";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api/v1";

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(path: string, init?: RequestInit): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      ...init,
      headers: {
        "Content-Type": "application/json",
        ...init?.headers,
      },
    });

    const json = await res.json();

    if (!res.ok || json.success === false) {
      throw new ApiRequestError(
        json.error?.message ?? "Unknown error",
        json.error?.code ?? "UNKNOWN",
        res.status,
      );
    }

    return json as T;
  }

  private authHeaders(token: string): HeadersInit {
    return { Authorization: `Bearer ${token}` };
  }

  // --- Articles ---

  async getArticles(params?: {
    page?: number;
    per_page?: number;
    tag?: string;
    source_id?: string;
    sort?: "published_at" | "collected_at";
    order?: "asc" | "desc";
  }): Promise<PaginatedResponse<ArticleSummary>> {
    const qs = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined) qs.set(k, String(v));
      });
    }
    const query = qs.toString();
    return this.request(`/articles${query ? `?${query}` : ""}`);
  }

  async getArticle(id: string): Promise<ApiResponse<ArticleDetail>> {
    return this.request(`/articles/${encodeURIComponent(id)}`);
  }

  async searchArticles(params: {
    q: string;
    page?: number;
    per_page?: number;
  }): Promise<PaginatedResponse<ArticleSummary>> {
    const qs = new URLSearchParams({ q: params.q });
    if (params.page) qs.set("page", String(params.page));
    if (params.per_page) qs.set("per_page", String(params.per_page));
    return this.request(`/articles/search?${qs}`);
  }

  // --- Tags ---

  async getTags(): Promise<ApiResponse<Tag[]>> {
    return this.request("/tags");
  }

  async getTagArticles(
    slug: string,
    params?: {
      page?: number;
      per_page?: number;
      sort?: string;
      order?: string;
    },
  ): Promise<PaginatedResponse<ArticleSummary>> {
    const qs = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined) qs.set(k, String(v));
      });
    }
    const query = qs.toString();
    return this.request(`/tags/${encodeURIComponent(slug)}/articles${query ? `?${query}` : ""}`);
  }

  // --- Sources ---

  async getSources(): Promise<ApiResponse<Source[]>> {
    return this.request("/sources");
  }

  // --- Bookmarks ---

  async getBookmarks(
    token: string,
    params?: { page?: number; per_page?: number },
  ): Promise<PaginatedResponse<Bookmark>> {
    const qs = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined) qs.set(k, String(v));
      });
    }
    const query = qs.toString();
    return this.request(`/bookmarks${query ? `?${query}` : ""}`, {
      headers: this.authHeaders(token),
    });
  }

  async addBookmark(token: string, articleId: string): Promise<ApiResponse<{ article_id: string; created_at: string }>> {
    return this.request("/bookmarks", {
      method: "POST",
      headers: this.authHeaders(token),
      body: JSON.stringify({ article_id: articleId }),
    });
  }

  async removeBookmark(token: string, articleId: string): Promise<void> {
    await this.request(`/bookmarks/${encodeURIComponent(articleId)}`, {
      method: "DELETE",
      headers: this.authHeaders(token),
    });
  }

  // --- Auth ---

  async loginWithGoogle(idToken: string) {
    return this.request<ApiResponse<{ access_token: string; refresh_token: string; user: User }>>(
      "/auth/google",
      {
        method: "POST",
        body: JSON.stringify({ id_token: idToken }),
      },
    );
  }

  // --- User ---

  async getMe(token: string): Promise<ApiResponse<User>> {
    return this.request("/users/me", {
      headers: this.authHeaders(token),
    });
  }
}

export class ApiRequestError extends Error {
  code: string;
  status: number;

  constructor(message: string, code: string, status: number) {
    super(message);
    this.name = "ApiRequestError";
    this.code = code;
    this.status = status;
  }
}

export const api = new ApiClient(API_BASE);
