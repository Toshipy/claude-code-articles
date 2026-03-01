# ClaudeCode Articles Platform

Claude Code に関する最新技術記事（リンク・サムネイル）を自動収集・集約し、開発者が一箇所で効率的に情報収集できる専用プラットフォームです。

## 主要機能

- **記事自動収集**: RSS / Web スクレイピングによる Claude Code 関連記事の定期収集（Zenn, Qiita, note, dev.to, Medium 等）
- **記事一覧表示**: サムネイル・タイトル・概要・ソース付きのカード形式表示
- **検索・フィルタリング**: キーワード検索、ソース別・日付別フィルタ
- **タグ分類**: 記事の自動カテゴリ分類（入門、Tips、比較、アップデート等）
- **ブックマーク**: ユーザーごとのお気に入り記事保存（Google OAuth 認証）

## 技術スタック

| レイヤー | 技術 |
|---------|------|
| Frontend | Next.js 15 + shadcn/ui + Tailwind CSS 4 |
| Backend API | Go 1.23 + Echo v4 |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Infra | AWS（ECS Fargate / RDS / CloudFront / S3） |

## ディレクトリ構成

```
claude-code-articles/
├── frontend/          # Next.js アプリ
│   ├── app/           # App Router（ページ）
│   ├── components/    # UI コンポーネント
│   └── lib/           # API クライアント・型定義
├── backend/           # Go API サーバー
│   ├── internal/
│   │   ├── handler/   # HTTP ハンドラー
│   │   ├── service/   # ビジネスロジック
│   │   ├── repository/# データアクセス
│   │   └── domain/    # ドメインモデル
│   └── migrations/    # DB マイグレーション
├── infra/             # Terraform（AWS インフラ）
└── docs/              # 設計ドキュメント
```

## 設計ドキュメント

| ドキュメント | 説明 |
|-------------|------|
| [アーキテクチャ設計書](docs/architecture.md) | システム全体構成・技術スタック・データフロー |
| [API 設計書](docs/api-design.md) | REST API エンドポイント仕様・認証方式・OpenAPI 定義 |
| [DB 設計書](docs/db-design.md) | テーブル定義・ER 図・インデックス設計・マイグレーション方針 |
| [UI 設計書](docs/ui-design.md) | デザインコンセプト・カラーパレット・コンポーネント設計・ワイヤーフレーム |
| [インフラ設計書](docs/infra-design.md) | AWS 構成・VPC 設計・CI/CD パイプライン・セキュリティ・コスト見積もり |

## 開発環境のセットアップ

```bash
# リポジトリのクローン
git clone https://github.com/Toshipy/claude-code-articles.git
cd claude-code-articles

# フロントエンド依存関係インストール
cd frontend && npm install

# バックエンド依存関係インストール
cd ../backend && go mod download

# ローカル環境起動（docker-compose）
docker-compose up -d
```

## API エンドポイント（主要）

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/v1/articles` | 記事一覧（ページネーション付き） |
| GET | `/api/v1/articles/search?q=` | 全文検索 |
| GET | `/api/v1/tags` | タグ一覧 |
| GET | `/api/v1/tags/:slug/articles` | タグ別記事一覧 |
| POST | `/api/v1/auth/google` | Google OAuth 認証 |
| GET | `/api/v1/bookmarks` | ブックマーク一覧（要認証） |

詳細は [API 設計書](docs/api-design.md) を参照してください。
