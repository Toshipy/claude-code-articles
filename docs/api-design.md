# ClaudeCode Articles Platform - API設計書

## 1. 設計方針

### 基本原則
- **RESTful API**: リソース指向のURL設計、適切なHTTPメソッドの使用
- **バージョニング**: URLパスベース `/api/v1/`
- **レスポンス形式**: JSON（`Content-Type: application/json`）
- **文字エンコーディング**: UTF-8
- **日時形式**: ISO 8601（`2026-02-19T12:00:00Z`）
- **ID形式**: UUIDv4

### HTTPメソッド
| メソッド | 用途 |
|---------|------|
| GET | リソースの取得 |
| POST | リソースの作成 |
| PUT | リソースの全体更新 |
| PATCH | リソースの部分更新 |
| DELETE | リソースの削除 |

### HTTPステータスコード
| コード | 用途 |
|--------|------|
| 200 | 成功（データ返却あり） |
| 201 | 作成成功 |
| 204 | 成功（データ返却なし） |
| 400 | リクエスト不正 |
| 401 | 認証エラー |
| 403 | 権限エラー |
| 404 | リソース未検出 |
| 409 | 競合（重複等） |
| 422 | バリデーションエラー |
| 429 | レート制限超過 |
| 500 | サーバーエラー |

---

## 2. 認証方式

### JWT Bearer Token
- ログイン後に発行されるJWTをAuthorizationヘッダーに付与
- トークン有効期限: アクセストークン 1時間、リフレッシュトークン 30日

```
Authorization: Bearer <jwt_token>
```

#### JWTペイロード
```json
{
  "sub": "user-uuid",
  "email": "user@example.com",
  "role": "user",
  "iat": 1708300800,
  "exp": 1708304400
}
```

### Google OAuth2 フロー
1. クライアントがGoogleからIDトークンを取得
2. `POST /api/v1/auth/google` にIDトークンを送信
3. サーバーがIDトークンを検証し、JWTを発行
4. 初回ログイン時はユーザーを自動作成

### 認証が必要なエンドポイント
| エンドポイント | 認証 |
|---------------|------|
| 記事一覧・詳細・検索 | 不要 |
| タグ一覧・タグ別記事 | 不要 |
| ソース一覧 | 不要 |
| ユーザープロフィール | 必要 |
| ブックマーク操作 | 必要 |
| 管理者API | 必要（admin権限） |

---

## 3. 共通レスポンス形式

### 成功レスポンス（単一リソース）
```json
{
  "success": true,
  "data": { ... }
}
```

### 成功レスポンス（リスト・ページネーション付き）
```json
{
  "success": true,
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

### エラーレスポンス
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "リクエストパラメータが不正です",
    "details": [
      {
        "field": "email",
        "message": "有効なメールアドレスを入力してください"
      }
    ]
  }
}
```

---

## 4. エラーコード体系

| コード | HTTPステータス | 説明 |
|--------|--------------|------|
| `UNAUTHORIZED` | 401 | 認証が必要 |
| `INVALID_TOKEN` | 401 | トークンが無効または期限切れ |
| `FORBIDDEN` | 403 | 権限不足 |
| `NOT_FOUND` | 404 | リソースが見つからない |
| `VALIDATION_ERROR` | 422 | バリデーションエラー |
| `DUPLICATE_RESOURCE` | 409 | リソースが既に存在 |
| `RATE_LIMITED` | 429 | レート制限超過 |
| `INTERNAL_ERROR` | 500 | サーバー内部エラー |
| `BAD_REQUEST` | 400 | リクエスト形式が不正 |

---

## 5. エンドポイント一覧

| メソッド | パス | 説明 | 認証 |
|---------|------|------|------|
| GET | `/api/v1/articles` | 記事一覧 | 不要 |
| GET | `/api/v1/articles/:id` | 記事詳細 | 不要 |
| GET | `/api/v1/articles/search` | 全文検索 | 不要 |
| GET | `/api/v1/tags` | タグ一覧 | 不要 |
| GET | `/api/v1/tags/:slug/articles` | タグ別記事一覧 | 不要 |
| GET | `/api/v1/sources` | ソース一覧 | 不要 |
| POST | `/api/v1/auth/google` | Google OAuth認証 | 不要 |
| GET | `/api/v1/users/me` | プロフィール取得 | 必要 |
| POST | `/api/v1/bookmarks` | ブックマーク追加 | 必要 |
| DELETE | `/api/v1/bookmarks/:article_id` | ブックマーク削除 | 必要 |
| GET | `/api/v1/bookmarks` | ブックマーク一覧 | 必要 |
| POST | `/api/v1/admin/collect` | 記事収集トリガー | 必要(admin) |

---

## 6. エンドポイント詳細仕様

### 6.1 GET /api/v1/articles - 記事一覧

記事をページネーション付きで取得する。

#### クエリパラメータ
| パラメータ | 型 | 必須 | デフォルト | 説明 |
|-----------|-----|------|-----------|------|
| page | integer | No | 1 | ページ番号 |
| per_page | integer | No | 20 | 1ページあたりの件数（最大100） |
| tag | string | No | - | タグslugでフィルタ |
| source_id | string | No | - | ソースIDでフィルタ |
| sort | string | No | `published_at` | ソート項目: `published_at`, `collected_at` |
| order | string | No | `desc` | ソート順: `asc`, `desc` |

#### レスポンス 200
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "title": "Claude Codeで始めるAI駆動開発",
      "summary": "Claude Codeの基本的な使い方と実践的なワークフローを紹介します...",
      "url": "https://example.com/articles/claude-code-intro",
      "thumbnail_url": "https://example.com/images/thumb.jpg",
      "published_at": "2026-02-18T09:00:00Z",
      "collected_at": "2026-02-18T10:30:00Z",
      "source": {
        "id": "550e8400-e29b-41d4-a716-446655440010",
        "name": "Zenn",
        "icon_url": "https://example.com/icons/zenn.png"
      },
      "tags": [
        { "id": "tag-uuid-1", "name": "Claude Code", "slug": "claude-code" },
        { "id": "tag-uuid-2", "name": "AI開発", "slug": "ai-development" }
      ]
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

---

### 6.2 GET /api/v1/articles/:id - 記事詳細

指定IDの記事の詳細情報を取得する。

#### パスパラメータ
| パラメータ | 型 | 説明 |
|-----------|-----|------|
| id | string (UUID) | 記事ID |

#### レスポンス 200
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "title": "Claude Codeで始めるAI駆動開発",
    "summary": "Claude Codeの基本的な使い方と実践的なワークフローを紹介します...",
    "content": "## はじめに\nClaude Codeは...",
    "url": "https://example.com/articles/claude-code-intro",
    "thumbnail_url": "https://example.com/images/thumb.jpg",
    "author": "tech_writer",
    "published_at": "2026-02-18T09:00:00Z",
    "collected_at": "2026-02-18T10:30:00Z",
    "source": {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "name": "Zenn",
      "url": "https://zenn.dev",
      "icon_url": "https://example.com/icons/zenn.png"
    },
    "tags": [
      { "id": "tag-uuid-1", "name": "Claude Code", "slug": "claude-code" },
      { "id": "tag-uuid-2", "name": "AI開発", "slug": "ai-development" }
    ],
    "related_articles": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440002",
        "title": "Claude Code Tips & Tricks",
        "thumbnail_url": "https://example.com/images/thumb2.jpg"
      }
    ]
  }
}
```

#### レスポンス 404
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "指定された記事が見つかりません"
  }
}
```

---

### 6.3 GET /api/v1/articles/search - 全文検索

記事をキーワードで全文検索する。

#### クエリパラメータ
| パラメータ | 型 | 必須 | デフォルト | 説明 |
|-----------|-----|------|-----------|------|
| q | string | Yes | - | 検索キーワード（2文字以上） |
| page | integer | No | 1 | ページ番号 |
| per_page | integer | No | 20 | 1ページあたりの件数（最大100） |

#### レスポンス 200
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "title": "Claude Codeで始めるAI駆動開発",
      "summary": "Claude Codeの基本的な使い方と...",
      "url": "https://example.com/articles/claude-code-intro",
      "thumbnail_url": "https://example.com/images/thumb.jpg",
      "published_at": "2026-02-18T09:00:00Z",
      "source": {
        "id": "550e8400-e29b-41d4-a716-446655440010",
        "name": "Zenn",
        "icon_url": "https://example.com/icons/zenn.png"
      },
      "tags": [
        { "id": "tag-uuid-1", "name": "Claude Code", "slug": "claude-code" }
      ],
      "highlight": "...実践的な<mark>Claude Code</mark>の活用法を..."
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 5,
    "total_pages": 1
  }
}
```

#### レスポンス 422
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "検索キーワードは2文字以上入力してください",
    "details": [
      { "field": "q", "message": "2文字以上入力してください" }
    ]
  }
}
```

---

### 6.4 GET /api/v1/tags - タグ一覧

全タグを記事数付きで取得する。

#### レスポンス 200
```json
{
  "success": true,
  "data": [
    {
      "id": "tag-uuid-1",
      "name": "Claude Code",
      "slug": "claude-code",
      "article_count": 45
    },
    {
      "id": "tag-uuid-2",
      "name": "AI開発",
      "slug": "ai-development",
      "article_count": 30
    },
    {
      "id": "tag-uuid-3",
      "name": "プロンプトエンジニアリング",
      "slug": "prompt-engineering",
      "article_count": 22
    }
  ]
}
```

---

### 6.5 GET /api/v1/tags/:slug/articles - タグ別記事一覧

指定タグに紐づく記事をページネーション付きで取得する。

#### パスパラメータ
| パラメータ | 型 | 説明 |
|-----------|-----|------|
| slug | string | タグのslug |

#### クエリパラメータ
| パラメータ | 型 | 必須 | デフォルト | 説明 |
|-----------|-----|------|-----------|------|
| page | integer | No | 1 | ページ番号 |
| per_page | integer | No | 20 | 1ページあたりの件数（最大100） |
| sort | string | No | `published_at` | ソート項目 |
| order | string | No | `desc` | ソート順 |

#### レスポンス 200
記事一覧と同一形式。

#### レスポンス 404
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "指定されたタグが見つかりません"
  }
}
```

---

### 6.6 GET /api/v1/sources - ソース一覧

記事収集元のソース一覧を取得する。

#### レスポンス 200
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "name": "Zenn",
      "url": "https://zenn.dev",
      "icon_url": "https://example.com/icons/zenn.png",
      "article_count": 80,
      "last_collected_at": "2026-02-19T06:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440011",
      "name": "Qiita",
      "url": "https://qiita.com",
      "icon_url": "https://example.com/icons/qiita.png",
      "article_count": 55,
      "last_collected_at": "2026-02-19T06:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440012",
      "name": "note",
      "url": "https://note.com",
      "icon_url": "https://example.com/icons/note.png",
      "article_count": 15,
      "last_collected_at": "2026-02-19T06:00:00Z"
    }
  ]
}
```

---

### 6.7 POST /api/v1/auth/google - Google OAuth認証

GoogleのIDトークンを検証し、JWTアクセストークンを発行する。

#### リクエストボディ
```json
{
  "id_token": "eyJhbGciOiJSUzI1NiIs..."
}
```

#### レスポンス 200
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": "user-uuid-001",
      "email": "user@gmail.com",
      "name": "テストユーザー",
      "avatar_url": "https://lh3.googleusercontent.com/...",
      "role": "user",
      "created_at": "2026-02-19T12:00:00Z"
    }
  }
}
```

#### レスポンス 401
```json
{
  "success": false,
  "error": {
    "code": "INVALID_TOKEN",
    "message": "Google IDトークンの検証に失敗しました"
  }
}
```

---

### 6.8 GET /api/v1/users/me - プロフィール取得

認証済みユーザーの情報を取得する。

#### ヘッダー
```
Authorization: Bearer <jwt_token>
```

#### レスポンス 200
```json
{
  "success": true,
  "data": {
    "id": "user-uuid-001",
    "email": "user@gmail.com",
    "name": "テストユーザー",
    "avatar_url": "https://lh3.googleusercontent.com/...",
    "role": "user",
    "bookmark_count": 12,
    "created_at": "2026-02-19T12:00:00Z",
    "updated_at": "2026-02-19T12:00:00Z"
  }
}
```

#### レスポンス 401
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "認証が必要です"
  }
}
```

---

### 6.9 POST /api/v1/bookmarks - ブックマーク追加

記事をブックマークに追加する。

#### ヘッダー
```
Authorization: Bearer <jwt_token>
```

#### リクエストボディ
```json
{
  "article_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

#### レスポンス 201
```json
{
  "success": true,
  "data": {
    "article_id": "550e8400-e29b-41d4-a716-446655440001",
    "created_at": "2026-02-19T12:30:00Z"
  }
}
```

#### レスポンス 409
```json
{
  "success": false,
  "error": {
    "code": "DUPLICATE_RESOURCE",
    "message": "この記事は既にブックマークされています"
  }
}
```

---

### 6.10 DELETE /api/v1/bookmarks/:article_id - ブックマーク削除

ブックマークを削除する。

#### パスパラメータ
| パラメータ | 型 | 説明 |
|-----------|-----|------|
| article_id | string (UUID) | 記事ID |

#### ヘッダー
```
Authorization: Bearer <jwt_token>
```

#### レスポンス 204
レスポンスボディなし。

#### レスポンス 404
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "指定されたブックマークが見つかりません"
  }
}
```

---

### 6.11 GET /api/v1/bookmarks - ブックマーク一覧

認証ユーザーのブックマーク一覧をページネーション付きで取得する。

#### ヘッダー
```
Authorization: Bearer <jwt_token>
```

#### クエリパラメータ
| パラメータ | 型 | 必須 | デフォルト | 説明 |
|-----------|-----|------|-----------|------|
| page | integer | No | 1 | ページ番号 |
| per_page | integer | No | 20 | 1ページあたりの件数（最大100） |

#### レスポンス 200
```json
{
  "success": true,
  "data": [
    {
      "article_id": "550e8400-e29b-41d4-a716-446655440001",
      "bookmarked_at": "2026-02-19T12:30:00Z",
      "article": {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "title": "Claude Codeで始めるAI駆動開発",
        "summary": "Claude Codeの基本的な使い方と...",
        "url": "https://example.com/articles/claude-code-intro",
        "thumbnail_url": "https://example.com/images/thumb.jpg",
        "published_at": "2026-02-18T09:00:00Z",
        "source": {
          "id": "550e8400-e29b-41d4-a716-446655440010",
          "name": "Zenn",
          "icon_url": "https://example.com/icons/zenn.png"
        },
        "tags": [
          { "id": "tag-uuid-1", "name": "Claude Code", "slug": "claude-code" }
        ]
      }
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 12,
    "total_pages": 1
  }
}
```

---

### 6.12 POST /api/v1/admin/collect - 記事収集トリガー

記事収集バッチを手動実行する（管理者専用）。

#### ヘッダー
```
Authorization: Bearer <jwt_token>
```

#### リクエストボディ（任意）
```json
{
  "source_id": "550e8400-e29b-41d4-a716-446655440010"
}
```
`source_id`を省略した場合、全ソースから収集を実行する。

#### レスポンス 200
```json
{
  "success": true,
  "data": {
    "job_id": "job-uuid-001",
    "status": "started",
    "target_sources": ["Zenn", "Qiita", "note"],
    "started_at": "2026-02-19T13:00:00Z"
  }
}
```

#### レスポンス 403
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "管理者権限が必要です"
  }
}
```

---

## 7. レート制限方針

### 制限値
| 対象 | 制限 | ウィンドウ |
|------|------|-----------|
| 未認証ユーザー | 60リクエスト | 1分あたり |
| 認証済みユーザー | 120リクエスト | 1分あたり |
| 検索エンドポイント | 30リクエスト | 1分あたり |
| 管理者API | 10リクエスト | 1分あたり |

### レスポンスヘッダー
```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1708300860
```

### 制限超過時のレスポンス 429
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMITED",
    "message": "リクエスト数の制限を超えました。しばらく待ってから再試行してください",
    "details": [
      { "field": "retry_after", "message": "30秒後に再試行してください" }
    ]
  }
}
```

---

## 8. OpenAPI 3.0 仕様（主要エンドポイント）

```yaml
openapi: 3.0.3
info:
  title: ClaudeCode Articles Platform API
  description: ClaudeCode関連記事を収集・配信するプラットフォームのAPI
  version: 1.0.0

servers:
  - url: https://api.claudecode-articles.example.com/api/v1
    description: Production
  - url: http://localhost:8080/api/v1
    description: Development

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  schemas:
    Source:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
        url:
          type: string
          format: uri
        icon_url:
          type: string
          format: uri

    Tag:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
        slug:
          type: string

    ArticleSummary:
      type: object
      properties:
        id:
          type: string
          format: uuid
        title:
          type: string
        summary:
          type: string
        url:
          type: string
          format: uri
        thumbnail_url:
          type: string
          format: uri
        published_at:
          type: string
          format: date-time
        collected_at:
          type: string
          format: date-time
        source:
          $ref: '#/components/schemas/Source'
        tags:
          type: array
          items:
            $ref: '#/components/schemas/Tag'

    Pagination:
      type: object
      properties:
        page:
          type: integer
        per_page:
          type: integer
        total:
          type: integer
        total_pages:
          type: integer

    ErrorResponse:
      type: object
      properties:
        success:
          type: boolean
          example: false
        error:
          type: object
          properties:
            code:
              type: string
            message:
              type: string
            details:
              type: array
              items:
                type: object
                properties:
                  field:
                    type: string
                  message:
                    type: string

paths:
  /articles:
    get:
      summary: 記事一覧取得
      tags: [Articles]
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            default: 1
        - name: per_page
          in: query
          schema:
            type: integer
            default: 20
            maximum: 100
        - name: tag
          in: query
          schema:
            type: string
        - name: source_id
          in: query
          schema:
            type: string
            format: uuid
        - name: sort
          in: query
          schema:
            type: string
            enum: [published_at, collected_at]
            default: published_at
        - name: order
          in: query
          schema:
            type: string
            enum: [asc, desc]
            default: desc
      responses:
        '200':
          description: 記事一覧
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/ArticleSummary'
                  pagination:
                    $ref: '#/components/schemas/Pagination'

  /articles/{id}:
    get:
      summary: 記事詳細取得
      tags: [Articles]
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: 記事詳細
        '404':
          description: 記事未検出
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /articles/search:
    get:
      summary: 記事全文検索
      tags: [Articles]
      parameters:
        - name: q
          in: query
          required: true
          schema:
            type: string
            minLength: 2
        - name: page
          in: query
          schema:
            type: integer
            default: 1
        - name: per_page
          in: query
          schema:
            type: integer
            default: 20
      responses:
        '200':
          description: 検索結果
        '422':
          description: バリデーションエラー
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /tags:
    get:
      summary: タグ一覧取得
      tags: [Tags]
      responses:
        '200':
          description: タグ一覧

  /auth/google:
    post:
      summary: Google OAuth認証
      tags: [Auth]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [id_token]
              properties:
                id_token:
                  type: string
      responses:
        '200':
          description: 認証成功
        '401':
          description: トークン検証失敗

  /users/me:
    get:
      summary: ユーザープロフィール取得
      tags: [Users]
      security:
        - BearerAuth: []
      responses:
        '200':
          description: ユーザー情報
        '401':
          description: 認証エラー

  /bookmarks:
    get:
      summary: ブックマーク一覧取得
      tags: [Bookmarks]
      security:
        - BearerAuth: []
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            default: 1
        - name: per_page
          in: query
          schema:
            type: integer
            default: 20
      responses:
        '200':
          description: ブックマーク一覧
        '401':
          description: 認証エラー
    post:
      summary: ブックマーク追加
      tags: [Bookmarks]
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [article_id]
              properties:
                article_id:
                  type: string
                  format: uuid
      responses:
        '201':
          description: ブックマーク作成成功
        '409':
          description: 重複エラー

  /bookmarks/{article_id}:
    delete:
      summary: ブックマーク削除
      tags: [Bookmarks]
      security:
        - BearerAuth: []
      parameters:
        - name: article_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: 削除成功
        '404':
          description: ブックマーク未検出

  /admin/collect:
    post:
      summary: 記事収集トリガー
      tags: [Admin]
      security:
        - BearerAuth: []
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                source_id:
                  type: string
                  format: uuid
      responses:
        '200':
          description: 収集ジョブ開始
        '403':
          description: 権限エラー
```
