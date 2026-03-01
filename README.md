# ClaudeCode Articles Platform

Claude Code に関する最新技術記事（リンク・サムネイル）を自動収集・集約し、開発者が一箇所で効率的に情報収集できる専用プラットフォームです。

## 主要機能

- **記事自動収集**: RSS/Web スクレイピングによる Claude Code 関連記事の定期収集
- **記事一覧表示**: サムネイル・タイトル・概要・ソース付きのカード形式表示
- **検索・フィルタリング**: キーワード検索、ソース別・日付別フィルタ
- **タグ分類**: 記事の自動カテゴリ分類（入門、Tips、比較、アップデート等）
- **トレンド表示**: 人気記事・最新記事のランキング
- **ブックマーク**: ユーザーごとのお気に入り記事保存

## 技術スタック

| レイヤー | 技術 | バージョン |
|---------|------|-----------|
| Frontend | Next.js | 15.x |
| UI Components | shadcn/ui | latest |
| CSS | Tailwind CSS | 4.x |
| Backend API | Go | 1.23+ |
| Web Framework | Echo | v4 |
| Database | PostgreSQL | 16 |
| Cache | Redis | 7.x |
| Infra | AWS | - |
| Container | Docker | - |
| CI/CD | GitHub Actions | - |

## ディレクトリ構成

```
claude-code-articles/
├── frontend/                 # Next.js アプリ
│   ├── app/                  # App Router
│   ├── components/           # UIコンポーネント
│   └── lib/                  # ユーティリティ
├── backend/                  # Go API サーバー
│   ├── internal/
│   │   ├── handler/          # HTTPハンドラー
│   │   ├── service/          # ビジネスロジック
│   │   ├── repository/       # データアクセス
│   │   └── domain/model/     # データモデル
│   └── migrations/           # DBマイグレーション
├── docs/                     # 設計ドキュメント
└── .github/workflows/        # CI/CD
```

## 設計ドキュメント

| ドキュメント | 概要 |
|---|---|
| [アーキテクチャ設計書](docs/architecture.md) | システム全体の構成・技術スタック・データフロー |
| [API設計書](docs/api-design.md) | RESTful APIエンドポイント仕様・認証方式 |
| [DB設計書](docs/db-design.md) | テーブル定義・ER図・インデックス・マイグレーション方針 |
| [UI設計書](docs/ui-design.md) | デザインコンセプト・コンポーネント・ワイヤーフレーム |
| [インフラ設計書](docs/infra-design.md) | AWS構成・VPC設計・CI/CDパイプライン・コスト見積もり |

## 開発環境セットアップ

```bash
# リポジトリのクローン
git clone https://github.com/Toshipy/claude-code-articles.git
cd claude-code-articles

# フロントエンド
cd frontend
npm install
npm run dev

# バックエンド
cd backend
go mod download
go run ./cmd/api
```

## APIエンドポイント一覧

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/v1/articles` | 記事一覧 |
| GET | `/api/v1/articles/:id` | 記事詳細 |
| GET | `/api/v1/articles/search` | 全文検索 |
| GET | `/api/v1/tags` | タグ一覧 |
| GET | `/api/v1/tags/:slug/articles` | タグ別記事一覧 |
| GET | `/api/v1/sources` | ソース一覧 |
| POST | `/api/v1/auth/google` | Google OAuth認証 |
| GET | `/api/v1/users/me` | プロフィール取得（要認証） |
| GET | `/api/v1/bookmarks` | ブックマーク一覧（要認証） |
| POST | `/api/v1/bookmarks` | ブックマーク追加（要認証） |
| DELETE | `/api/v1/bookmarks/:article_id` | ブックマーク削除（要認証） |
