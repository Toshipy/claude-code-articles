# ClaudeCode Articles - CLAUDE.md

## プロジェクト概要

ClaudeCode に関する技術記事を集約するプラットフォーム。
Next.js フロントエンド + Go バックエンドで構成。

## 技術スタック

| レイヤー | 技術 |
|---------|------|
| フロントエンド | Next.js 15 (App Router), TypeScript, Tailwind CSS v4, Shadcn/ui |
| バックエンド | Go |
| スタイリング | ダークテーマ（`#0a0a0a` ベース）、アンバー（`#f59e0b`）アクセント |

## ディレクトリ構成

```
.
├── frontend/          # Next.js アプリケーション
│   ├── app/           # App Router ページ
│   ├── components/    # React コンポーネント
│   └── lib/           # ユーティリティ・型定義
├── backend/           # Go API サーバー
├── docs/              # 設計ドキュメント
│   ├── ui-design.md   # UI設計書（カラー・コンポーネント仕様）
│   ├── api-design.md  # API設計書
│   └── db-design.md   # DB設計書
└── CLAUDE.md          # このファイル
```

---

## Figma Code to Canvas セットアップ

### 概要

既存の React コンポーネント（`frontend/components/`）を Figma キャンバスにデザインとして起こすワークフローです。
Figma MCP サーバーを通じて Claude Code が直接 Figma ファイルを操作します。

### 1. 事前準備

#### Figma API トークンの取得

1. [Figma](https://www.figma.com) にログイン
2. アカウント設定 > **Personal access tokens** へ移動
3. トークンを生成してコピー

#### 環境変数の設定

```bash
export FIGMA_API_KEY="your-figma-api-token"
```

または `.env.local`（git 管理外）に記載:

```
FIGMA_API_KEY=your-figma-api-token
```

### 2. MCP サーバーの起動確認

`.claude/settings.json` に Figma MCP サーバーの設定が含まれています。
Claude Code 起動時に自動的に MCP サーバーが接続されます。

```bash
# Claude Code を MCP 有効で起動
claude
```

### 3. Code to Canvas ワークフロー

#### 基本的な使い方

Claude Code のプロンプトで以下のように依頼します：

```
以下のコンポーネントを Figma ファイル [FILE_ID] に変換してください：
- frontend/components/ArticleCard.tsx
- frontend/components/Header.tsx
```

#### コンポーネント → Figma マッピング

| コンポーネント | Figma フレーム名 | 備考 |
|---------------|-----------------|------|
| `ArticleCard` (default) | `Card/Default` | アスペクト比 16:9 サムネイル付きカード |
| `ArticleCard` (compact) | `Card/Compact` | サイドバー用テキストのみ |
| `ArticleCard` (horizontal) | `Card/Horizontal` | 検索結果用横型 |
| `Header` | `Layout/Header` | ロゴ・検索・ユーザーメニュー |
| `TagBadge` | `Badge/Tag` | アウトラインバリアント |
| `SearchBar` | `Search/CommandPalette` | Cmd+K モーダル |
| `Sidebar` | `Layout/Sidebar` | 人気記事・タグクラウド |

#### デザイントークン（カラーパレット）

Figma にカラースタイルを作成する際は以下の値を使用してください：

```
Background/Base     #0a0a0a
Background/Card     #141414
Background/Elevated #1a1a1a
Border/Default      #262626
Text/Primary        #fafafa
Text/Secondary      #a1a1aa
Accent/Primary      #f59e0b  (Amber-500)
Accent/Hover        #fbbf24  (Amber-400)
Accent/Muted        #78350f  (Amber-900)
Success             #22c55e
Destructive         #ef4444
Code/Background     #0d1117
```

#### タイポグラフィ

```
Font: Inter (UI・見出し) / JetBrains Mono (コード)
H1: 30px / Bold / line-height 1.3
H2: 24px / SemiBold / line-height 1.35
H3: 20px / SemiBold / line-height 1.4
Body: 16px / Regular / line-height 1.7
Small: 14px / Regular / line-height 1.5
XSmall: 12px / Regular / line-height 1.4
```

---

## 開発ガイドライン

- コンポーネントの変更は `docs/ui-design.md` の仕様に従う
- 新しいカラーは必ず CSS 変数（`--variable-name`）で定義する
- Figma のデザインは `docs/ui-design.md` のワイヤーフレームと整合性を保つ
