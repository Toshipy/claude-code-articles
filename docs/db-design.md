# ClaudeCode Articles Platform - DB設計書

## 1. 設計方針

### 基本方針
- **RDBMS**: PostgreSQL 16+
- **命名規則**: スネークケース（snake_case）、テーブル名は複数形
- **主キー**: BIGSERIAL（`id`）を基本とする。中間テーブルは複合主キー
- **タイムスタンプ**: `created_at`/`updated_at` は `TIMESTAMPTZ`（タイムゾーン付き）で管理
- **ソフトデリート**: 採用しない。不要データは物理削除
- **文字列型**: 可変長が基本。短い固定長にはVARCHAR(n)、長文にはTEXT
- **NULL許容**: 必須項目はNOT NULL制約を付与。任意項目のみNULL許容

### パフォーマンス方針
- 検索・ソートに使うカラムにはインデックスを付与
- articlesテーブルは`published_at`による時系列レンジパーティショニングを採用
- 全文検索にはPostgreSQLの`tsvector`/GINインデックスを活用

### マイグレーション方針
- **ツール**: [golang-migrate](https://github.com/golang-migrate/migrate)
- マイグレーションファイルは `db/migrations/` に `NNNNNN_description.up.sql` / `NNNNNN_description.down.sql` の形式で管理
- 全マイグレーションはトランザクション内で実行し、失敗時にロールバック可能とする

---

## 2. テーブル一覧

| テーブル名 | 概要 |
|---|---|
| `sources` | 記事取得元（RSSフィード・スクレイピング対象サイト） |
| `articles` | 収集した記事本体（パーティションテーブル） |
| `tags` | 記事に付与するタグマスタ |
| `article_tags` | 記事とタグの多対多中間テーブル |
| `users` | ユーザー情報 |
| `bookmarks` | ユーザーの記事ブックマーク |
| `search_history` | ユーザーの検索履歴 |

---

## 3. テーブル詳細定義

### 3.1 sources（記事取得元）

| カラム名 | 型 | 制約 | 説明 |
|---|---|---|---|
| `id` | `BIGSERIAL` | PRIMARY KEY | 取得元ID |
| `name` | `VARCHAR(255)` | NOT NULL | サイト名 |
| `url` | `TEXT` | NOT NULL, UNIQUE | サイトURL |
| `feed_url` | `TEXT` | | RSSフィードURL（NULLならスクレイピング対象） |
| `last_fetched_at` | `TIMESTAMPTZ` | | 最終取得日時 |
| `is_active` | `BOOLEAN` | NOT NULL DEFAULT true | 有効フラグ |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT NOW() | 作成日時 |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT NOW() | 更新日時 |

```sql
CREATE TABLE sources (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    url           TEXT NOT NULL UNIQUE,
    feed_url      TEXT,
    last_fetched_at TIMESTAMPTZ,
    is_active     BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 3.2 articles（記事）

| カラム名 | 型 | 制約 | 説明 |
|---|---|---|---|
| `id` | `BIGSERIAL` | PRIMARY KEY | 記事ID |
| `source_id` | `BIGINT` | NOT NULL, FK → sources(id) | 取得元ID |
| `title` | `VARCHAR(500)` | NOT NULL | 記事タイトル |
| `url` | `TEXT` | NOT NULL, UNIQUE | 記事URL（重複取り込み防止） |
| `thumbnail_url` | `TEXT` | | サムネイル画像URL |
| `summary` | `TEXT` | | 記事要約 |
| `source` | `VARCHAR(100)` | NOT NULL | 記事ソース名（表示用） |
| `published_at` | `TIMESTAMPTZ` | NOT NULL | 記事公開日時（パーティションキー） |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT NOW() | 取り込み日時 |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT NOW() | 更新日時 |

```sql
CREATE TABLE articles (
    id            BIGSERIAL,
    source_id     BIGINT NOT NULL,
    title         VARCHAR(500) NOT NULL,
    url           TEXT NOT NULL,
    thumbnail_url TEXT,
    summary       TEXT,
    source        VARCHAR(100) NOT NULL,
    published_at  TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, published_at),
    UNIQUE (url, published_at),
    FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE
) PARTITION BY RANGE (published_at);
```

### 3.3 tags（タグ）

| カラム名 | 型 | 制約 | 説明 |
|---|---|---|---|
| `id` | `BIGSERIAL` | PRIMARY KEY | タグID |
| `name` | `VARCHAR(100)` | NOT NULL, UNIQUE | タグ表示名 |
| `slug` | `VARCHAR(100)` | NOT NULL, UNIQUE | URLスラッグ |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT NOW() | 作成日時 |

```sql
CREATE TABLE tags (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL UNIQUE,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 3.4 article_tags（記事タグ中間テーブル）

| カラム名 | 型 | 制約 | 説明 |
|---|---|---|---|
| `article_id` | `BIGINT` | NOT NULL, FK → articles(id) | 記事ID |
| `tag_id` | `BIGINT` | NOT NULL, FK → tags(id) | タグID |

```sql
CREATE TABLE article_tags (
    article_id BIGINT NOT NULL,
    tag_id     BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, tag_id)
);
```

> **Note**: `article_id` の外部キー制約はパーティションテーブルへの参照のため、アプリケーション層で整合性を担保する。

### 3.5 users（ユーザー）

| カラム名 | 型 | 制約 | 説明 |
|---|---|---|---|
| `id` | `BIGSERIAL` | PRIMARY KEY | ユーザーID |
| `email` | `VARCHAR(255)` | NOT NULL, UNIQUE | メールアドレス |
| `display_name` | `VARCHAR(100)` | NOT NULL | 表示名 |
| `avatar_url` | `TEXT` | | アバター画像URL |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT NOW() | 登録日時 |

```sql
CREATE TABLE users (
    id           BIGSERIAL PRIMARY KEY,
    email        VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    avatar_url   TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 3.6 bookmarks（ブックマーク）

| カラム名 | 型 | 制約 | 説明 |
|---|---|---|---|
| `user_id` | `BIGINT` | NOT NULL, FK → users(id) | ユーザーID |
| `article_id` | `BIGINT` | NOT NULL | 記事ID |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT NOW() | ブックマーク日時 |

```sql
CREATE TABLE bookmarks (
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    article_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, article_id)
);
```

### 3.7 search_history（検索履歴）

| カラム名 | 型 | 制約 | 説明 |
|---|---|---|---|
| `id` | `BIGSERIAL` | PRIMARY KEY | 検索履歴ID |
| `user_id` | `BIGINT` | NOT NULL, FK → users(id) | ユーザーID |
| `query` | `VARCHAR(500)` | NOT NULL | 検索クエリ |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT NOW() | 検索日時 |

```sql
CREATE TABLE search_history (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    query      VARCHAR(500) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 4. ER図（Mermaid）

```mermaid
erDiagram
    sources ||--o{ articles : "has many"
    articles ||--o{ article_tags : "has many"
    tags ||--o{ article_tags : "has many"
    users ||--o{ bookmarks : "has many"
    articles ||--o{ bookmarks : "has many"
    users ||--o{ search_history : "has many"

    sources {
        bigint id PK
        varchar name
        text url UK
        text feed_url
        timestamptz last_fetched_at
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    articles {
        bigint id PK
        bigint source_id FK
        varchar title
        text url UK
        text thumbnail_url
        text summary
        varchar source
        timestamptz published_at PK
        timestamptz created_at
        timestamptz updated_at
    }

    tags {
        bigint id PK
        varchar name UK
        varchar slug UK
        timestamptz created_at
    }

    article_tags {
        bigint article_id PK,FK
        bigint tag_id PK,FK
    }

    users {
        bigint id PK
        varchar email UK
        varchar display_name
        text avatar_url
        timestamptz created_at
    }

    bookmarks {
        bigint user_id PK,FK
        bigint article_id PK,FK
        timestamptz created_at
    }

    search_history {
        bigint id PK
        bigint user_id FK
        varchar query
        timestamptz created_at
    }
```

---

## 5. インデックス設計

### articles

```sql
-- 公開日時での降順ソート（一覧表示用）
CREATE INDEX idx_articles_published_at ON articles (published_at DESC);

-- ソース別記事一覧
CREATE INDEX idx_articles_source_id ON articles (source_id);

-- 全文検索用（タイトル + 要約）
ALTER TABLE articles ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(summary, ''))
    ) STORED;
CREATE INDEX idx_articles_search ON articles USING GIN (search_vector);
```

### article_tags

```sql
-- タグから記事を引く（タグ絞り込み用）
CREATE INDEX idx_article_tags_tag_id ON article_tags (tag_id);
```

### tags

```sql
-- スラッグによるルックアップ（URLルーティング用）
-- UNIQUE制約により自動的にインデックスが作成される
```

### bookmarks

```sql
-- 記事のブックマーク数集計用
CREATE INDEX idx_bookmarks_article_id ON bookmarks (article_id);
```

### search_history

```sql
-- ユーザー別検索履歴（最近の検索を表示）
CREATE INDEX idx_search_history_user_id_created ON search_history (user_id, created_at DESC);
```

### sources

```sql
-- アクティブなソースのフィルタリング
CREATE INDEX idx_sources_is_active ON sources (is_active) WHERE is_active = true;
```

---

## 6. パーティショニング方針

### 対象テーブル
`articles` テーブルに対し、`published_at` カラムによるレンジパーティショニングを適用する。

### パーティション粒度
- **月次パーティション** を採用
- 記事量の増加に応じて適切な粒度を維持

### パーティション作成例

```sql
-- 2025年のパーティション
CREATE TABLE articles_2025_01 PARTITION OF articles
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
CREATE TABLE articles_2025_02 PARTITION OF articles
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
-- ... 以降の月も同様

-- 2026年のパーティション
CREATE TABLE articles_2026_01 PARTITION OF articles
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE articles_2026_02 PARTITION OF articles
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
-- ... 以降の月も同様
```

### パーティション自動作成
- cronジョブまたはアプリケーション側で、3ヶ月先のパーティションを事前に作成する
- パーティションが存在しないとINSERTが失敗するため、余裕を持って作成する

### デフォルトパーティション

```sql
CREATE TABLE articles_default PARTITION OF articles DEFAULT;
```

> 想定外の日付データをキャッチするためのデフォルトパーティションを設置する。

---

## 7. マイグレーション管理方針

### ツール
[golang-migrate](https://github.com/golang-migrate/migrate) を使用。

### ディレクトリ構成

```
db/
└── migrations/
    ├── 000001_create_sources.up.sql
    ├── 000001_create_sources.down.sql
    ├── 000002_create_articles.up.sql
    ├── 000002_create_articles.down.sql
    ├── 000003_create_tags.up.sql
    ├── 000003_create_tags.down.sql
    ├── 000004_create_article_tags.up.sql
    ├── 000004_create_article_tags.down.sql
    ├── 000005_create_users.up.sql
    ├── 000005_create_users.down.sql
    ├── 000006_create_bookmarks.up.sql
    ├── 000006_create_bookmarks.down.sql
    ├── 000007_create_search_history.up.sql
    ├── 000007_create_search_history.down.sql
    └── 000008_create_indexes.up.sql
    └── 000008_create_indexes.down.sql
```

### 実行コマンド

```bash
# マイグレーション実行（全て適用）
migrate -path db/migrations -database "postgres://user:pass@localhost:5432/claude_code_articles?sslmode=disable" up

# 1つロールバック
migrate -path db/migrations -database "postgres://..." down 1

# バージョン確認
migrate -path db/migrations -database "postgres://..." version
```

### 運用ルール
1. マイグレーションファイルは一度適用したら**変更禁止**。修正は新しいマイグレーションで行う
2. `down.sql` は必ず作成し、ロールバック可能な状態を維持する
3. テーブル作成順序はFK依存関係を考慮（sources → articles → tags → article_tags → users → bookmarks → search_history）
4. 本番適用前にステージング環境で検証する
