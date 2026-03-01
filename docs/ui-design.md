# ClaudeCode Articles - UI設計書

## 1. デザインコンセプト

**"Terminal-Inspired Intelligence"**

ClaudeCodeのCLI体験を想起させる、ターミナルライクでありながら洗練されたモダンUI。
技術者が「使い慣れた空間」と感じるダークテーマをベースに、AIツールとしての先進性をアクセントカラーで表現する。

### デザイン原則

- **Clarity**: 情報密度を高く保ちつつ、視覚的階層を明確にする
- **Speed**: 高速な記事探索・閲覧体験を最優先する
- **Craft**: コードブロック・技術コンテンツの表示品質にこだわる
- **Warmth**: ダークテーマの冷たさをアンバー系アクセントで緩和する

---

## 2. カラーパレット

### ダークテーマ（デフォルト）

| 用途 | カラー | CSS変数 |
|------|--------|---------|
| 背景（Base） | `#0a0a0a` | `--background` |
| 背景（Card） | `#141414` | `--card` |
| 背景（Elevated） | `#1a1a1a` | `--muted` |
| ボーダー | `#262626` | `--border` |
| テキスト（Primary） | `#fafafa` | `--foreground` |
| テキスト（Secondary） | `#a1a1aa` | `--muted-foreground` |
| アクセント（Primary） | `#f59e0b` | `--primary` (Amber-500) |
| アクセント（Hover） | `#fbbf24` | `--primary-hover` (Amber-400) |
| アクセント（Subtle） | `#78350f` | `--primary-muted` (Amber-900) |
| Success | `#22c55e` | `--success` |
| Destructive | `#ef4444` | `--destructive` |
| コードブロック背景 | `#0d1117` | `--code-bg` |

### Tailwind CSS設定

```ts
// tailwind.config.ts
const config = {
  theme: {
    extend: {
      colors: {
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
        border: "hsl(var(--border))",
      },
    },
  },
};
```

---

## 3. タイポグラフィ

### フォントファミリー

| 用途 | フォント | Tailwindクラス |
|------|----------|----------------|
| 見出し・UI | Inter | `font-sans` |
| 本文 | Noto Sans JP + Inter | `font-sans` |
| コード | JetBrains Mono | `font-mono` |

### サイズ体系

| レベル | サイズ | 行間 | 用途 |
|--------|--------|------|------|
| Display | 36px / 2.25rem | 1.2 | ヒーローセクション |
| H1 | 30px / 1.875rem | 1.3 | ページタイトル |
| H2 | 24px / 1.5rem | 1.35 | セクション見出し |
| H3 | 20px / 1.25rem | 1.4 | カード見出し |
| Body | 16px / 1rem | 1.7 | 記事本文 |
| Small | 14px / 0.875rem | 1.5 | メタ情報・補足 |
| XSmall | 12px / 0.75rem | 1.4 | タグ・バッジ |

---

## 4. 使用Shadcn/uiコンポーネント一覧

| コンポーネント | 用途 |
|---------------|------|
| `Button` | CTA、ブックマーク、ナビゲーション |
| `Card` | 記事カード、関連記事 |
| `Badge` | タグ表示 |
| `Input` | 検索バー |
| `Dialog` | ログインモーダル |
| `DropdownMenu` | ユーザーメニュー、ソート |
| `Avatar` | ユーザーアイコン |
| `Separator` | セクション区切り |
| `Skeleton` | ローディング状態 |
| `Toast` | 通知（ブックマーク追加等） |
| `ScrollArea` | サイドバー、長いリスト |
| `Tooltip` | アイコンボタンの説明 |
| `Pagination` | 記事一覧ページネーション |
| `Breadcrumb` | パンくずナビ |
| `Sheet` | モバイルサイドバー |
| `Command` | 検索コマンドパレット（Cmd+K） |

---

## 5. ページ一覧・ワイヤーフレーム

### 5.1 トップページ (`/`)

```
┌──────────────────────────────────────────────────────────┐
│ [Logo] ClaudeCode Articles    [Search] [🔖] [User▼]     │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │  > ClaudeCodeを使いこなすための                     │  │
│  │    技術記事プラットフォーム_                         │  │
│  │                                                    │  │
│  │  [記事を探す]  [人気のタグ →]                       │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  新着記事                                     [もっと見る]│
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                 │
│  │ OGP画像  │ │ OGP画像  │ │ OGP画像  │   ┌──────────┐ │
│  │          │ │          │ │          │   │ Sidebar  │ │
│  │ タイトル │ │ タイトル │ │ タイトル │   │          │ │
│  │ 概要...  │ │ 概要...  │ │ 概要...  │   │ 人気記事 │ │
│  │ [tag]    │ │ [tag]    │ │ [tag]    │   │ 1. ----  │ │
│  │ 日付  🔖 │ │ 日付  🔖 │ │ 日付  🔖 │   │ 2. ----  │ │
│  └──────────┘ └──────────┘ └──────────┘   │ 3. ----  │ │
│                                            │          │ │
│  人気記事                                  │ タグ     │ │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐   │ [Claude] │ │
│  │ OGP画像  │ │ OGP画像  │ │ OGP画像  │   │ [MCP]    │ │
│  │ タイトル │ │ タイトル │ │ タイトル │   │ [Hooks]  │ │
│  │ 概要...  │ │ 概要...  │ │ 概要...  │   │ [Tips]   │ │
│  └──────────┘ └──────────┘ └──────────┘   └──────────┘ │
│                                                          │
│  タグクラウド                                            │
│  [ClaudeCode] [MCP] [Hooks] [プロンプト] [Tips]          │
│  [自動化] [設定] [ワークフロー] [CLI] [Agent]            │
│                                                          │
├──────────────────────────────────────────────────────────┤
│  Footer: © 2026 ClaudeCode Articles                      │
└──────────────────────────────────────────────────────────┘
```

### 5.2 記事詳細ページ (`/articles/[id]`)

```
┌──────────────────────────────────────────────────────────┐
│ [Logo] ClaudeCode Articles    [Search] [🔖] [User▼]     │
├──────────────────────────────────────────────────────────┤
│  Home > タグ名 > 記事タイトル                            │
│                                                          │
│  ┌──────────────────────────────────┐   ┌──────────────┐│
│  │                                  │   │   Sidebar    ││
│  │         OGP画像                  │   │              ││
│  │                                  │   │  目次(TOC)   ││
│  └──────────────────────────────────┘   │  ├ 見出し1   ││
│                                         │  ├ 見出し2   ││
│  [ClaudeCode] [MCP]                     │  └ 見出し3   ││
│  # 記事タイトル                          │              ││
│  2026-02-19 · 読了 5分                   │──────────────││
│                                         │              ││
│  記事本文...                             │  関連記事    ││
│  テキストテキストテキスト                │  ┌────────┐  ││
│                                         │  │ Card   │  ││
│  ```typescript                          │  └────────┘  ││
│  const example = "code block";          │  ┌────────┐  ││
│  ```                                    │  │ Card   │  ││
│                                         │  └────────┘  ││
│  テキストテキストテキスト                │              ││
│                                         │              ││
│  ┌────────────────────────────────┐     └──────────────┘│
│  │  [🔖 ブックマーク] [シェア ↗]  │                     │
│  └────────────────────────────────┘                     │
│                                                          │
│  関連記事（モバイル用・下部表示）                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                 │
│  │ Card     │ │ Card     │ │ Card     │                 │
│  └──────────┘ └──────────┘ └──────────┘                 │
├──────────────────────────────────────────────────────────┤
│  Footer                                                  │
└──────────────────────────────────────────────────────────┘
```

### 5.3 タグ別記事一覧 (`/tags/[slug]`)

```
┌──────────────────────────────────────────────────────────┐
│ [Logo] ClaudeCode Articles    [Search] [🔖] [User▼]     │
├──────────────────────────────────────────────────────────┤
│  Home > Tags > ClaudeCode                                │
│                                                          │
│  # ClaudeCode の記事 (24件)                              │
│  ClaudeCodeに関する技術記事の一覧です                     │
│                                                          │
│  並び替え: [新着順 ▼]                                    │
│                                                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                 │
│  │ OGP画像  │ │ OGP画像  │ │ OGP画像  │                 │
│  │ タイトル │ │ タイトル │ │ タイトル │                 │
│  │ 概要...  │ │ 概要...  │ │ 概要...  │                 │
│  │ [tag]    │ │ [tag]    │ │ [tag]    │                 │
│  └──────────┘ └──────────┘ └──────────┘                 │
│                                                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                 │
│  │ Card     │ │ Card     │ │ Card     │                 │
│  └──────────┘ └──────────┘ └──────────┘                 │
│                                                          │
│  [< 前へ]  1  2  3  ...  8  [次へ >]                     │
├──────────────────────────────────────────────────────────┤
│  Footer                                                  │
└──────────────────────────────────────────────────────────┘
```

### 5.4 検索結果ページ (`/search`)

```
┌──────────────────────────────────────────────────────────┐
│ [Logo] ClaudeCode Articles    [Search] [🔖] [User▼]     │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  ┌──────────────────────────────────────────────────┐    │
│  │ 🔍  "MCP サーバー"                          [×] │    │
│  └──────────────────────────────────────────────────┘    │
│                                                          │
│  "MCP サーバー" の検索結果: 12件                         │
│  フィルター: [すべて] [ClaudeCode] [MCP] [Tips]          │
│                                                          │
│  ┌──────────────────────────────────────────────────┐    │
│  │ 記事タイトル（マッチ部分ハイライト）              │    │
│  │ ...テキスト中の「MCP サーバー」がハイライト...     │    │
│  │ [MCP] [ClaudeCode]  2026-02-10                   │    │
│  └──────────────────────────────────────────────────┘    │
│  ┌──────────────────────────────────────────────────┐    │
│  │ 記事タイトル                                      │    │
│  │ ...スニペットテキスト...                          │    │
│  │ [MCP]  2026-02-05                                │    │
│  └──────────────────────────────────────────────────┘    │
│  ...                                                     │
│                                                          │
│  [< 前へ]  1  2  [次へ >]                                │
├──────────────────────────────────────────────────────────┤
│  Footer                                                  │
└──────────────────────────────────────────────────────────┘
```

### 5.5 ブックマーク一覧 (`/bookmarks`) ※要ログイン

```
┌──────────────────────────────────────────────────────────┐
│ [Logo] ClaudeCode Articles    [Search] [🔖] [User▼]     │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  # ブックマーク (8件)                                    │
│  保存した記事の一覧です                                   │
│                                                          │
│  並び替え: [保存日順 ▼]                                  │
│                                                          │
│  ┌──────────────────────────────────────────────────┐    │
│  │ □  記事タイトル                                   │    │
│  │    概要テキスト...                                │    │
│  │    [tag] [tag]  保存日: 2026-02-18    [🗑 削除]   │    │
│  └──────────────────────────────────────────────────┘    │
│  ┌──────────────────────────────────────────────────┐    │
│  │ □  記事タイトル                                   │    │
│  │    概要テキスト...                                │    │
│  │    [tag]  保存日: 2026-02-15           [🗑 削除]  │    │
│  └──────────────────────────────────────────────────┘    │
│  ...                                                     │
│                                                          │
│  選択した記事: [一括削除]                                │
│                                                          │
│  ─── 未ログイン時 ───                                    │
│  ┌──────────────────────────────────────────────────┐    │
│  │  ブックマーク機能を使うにはログインが必要です      │    │
│  │  [ログイン]  [新規登録]                           │    │
│  └──────────────────────────────────────────────────┘    │
├──────────────────────────────────────────────────────────┤
│  Footer                                                  │
└──────────────────────────────────────────────────────────┘
```

---

## 6. コンポーネント設計

### 6.1 Header

```tsx
// components/layout/Header.tsx
type HeaderProps = {
  user?: User | null;
};
```

- ロゴ（左）：ClaudeCode Articlesテキストロゴ + ターミナルカーソルアニメーション
- 検索（中央）：`Cmd+K`でCommand paletteを起動するトリガーボタン
- 右側：ブックマークアイコン、ユーザーメニュー（DropdownMenu）
- モバイル：ハンバーガーメニュー（Sheet）で展開
- スクロール時に`backdrop-blur`で半透明固定ヘッダー

### 6.2 ArticleCard

```tsx
// components/article/ArticleCard.tsx
type ArticleCardProps = {
  article: {
    id: string;
    title: string;
    excerpt: string;
    ogpImageUrl: string;
    tags: Tag[];
    publishedAt: string;
    readingTime: number;
  };
  variant?: "default" | "compact" | "horizontal";
};
```

- **default**: 縦型カード（OGP画像上部、テキスト下部）。トップページのグリッドで使用
- **compact**: 画像なし・テキストのみ。サイドバーの人気記事で使用
- **horizontal**: 横型（画像左・テキスト右）。検索結果で使用
- ホバー時: ボーダーがアンバーに変化 + 微小な`translateY(-2px)`

### 6.3 ArticleGrid

```tsx
// components/article/ArticleGrid.tsx
type ArticleGridProps = {
  articles: Article[];
  variant?: "default" | "compact";
  columns?: 2 | 3;
};
```

- CSS Grid使用: `grid-cols-1 md:grid-cols-2 lg:grid-cols-3`
- `gap-6`で統一的な余白

### 6.4 TagBadge

```tsx
// components/ui/TagBadge.tsx
type TagBadgeProps = {
  tag: {
    slug: string;
    name: string;
    count?: number;
  };
  size?: "sm" | "md";
  interactive?: boolean;
};
```

- Shadcn/uiの`Badge`をベースにカスタマイズ
- `variant="outline"`にアンバーのボーダー
- ホバー時に背景がアンバー系に変化
- `interactive=true`でクリック可能（タグページへ遷移）

### 6.5 SearchBar

```tsx
// components/search/SearchBar.tsx
type SearchBarProps = {
  defaultValue?: string;
  onSearch: (query: string) => void;
};
```

- Shadcn/uiの`Command`コンポーネントベース
- `Cmd+K` / `Ctrl+K`でどこからでも起動
- リアルタイムサジェスト（記事タイトル + タグ）
- 検索履歴の表示

### 6.6 Sidebar

```tsx
// components/layout/Sidebar.tsx
type SidebarProps = {
  popularArticles?: Article[];
  tags?: Tag[];
};
```

- デスクトップ: 右サイドバー（`w-80`固定幅）
- 人気記事ランキング（compact ArticleCard使用）
- タグクラウド（TagBadge使用）
- `sticky top-20`で追従

### 6.7 TableOfContents（目次）

```tsx
// components/article/TableOfContents.tsx
type TOCProps = {
  headings: { id: string; text: string; level: number }[];
};
```

- 記事詳細ページのサイドバーに表示
- Intersection Observerで現在位置をハイライト
- スムーススクロールでジャンプ

---

## 7. レスポンシブ対応方針

### ブレークポイント

| 名称 | 幅 | レイアウト |
|------|-----|-----------|
| Mobile | < 768px | 1カラム、ハンバーガーメニュー |
| Tablet | 768px - 1023px | 2カラムグリッド、サイドバー非表示 |
| Desktop | 1024px - 1279px | 3カラムグリッド + サイドバー |
| Wide | >= 1280px | max-width: 1280px 中央配置 |

### レスポンシブ詳細

- **Header**: モバイルではロゴ短縮 + Sheet（ドロワー）でナビ展開
- **ArticleGrid**: `grid-cols-1` → `md:grid-cols-2` → `lg:grid-cols-3`
- **Sidebar**: `lg:block hidden`。モバイル・タブレットでは非表示、コンテンツは下部に移動
- **記事詳細**: TOCはモバイルで`Sheet`として左からスライドイン
- **検索**: モバイルではフルスクリーンのCommand palette
- **フォントサイズ**: モバイルでBody 15px、Display 28px に縮小

---

## 8. アニメーション・インタラクション方針

### 基本方針

`prefers-reduced-motion`を尊重し、すべてのアニメーションを無効化できるようにする。

### アニメーション一覧

| 要素 | トリガー | アニメーション | Duration |
|------|----------|--------------|----------|
| ページ遷移 | ルート変更 | `opacity 0→1` + `translateY(8px→0)` | 200ms |
| ArticleCard | ホバー | `translateY(-2px)` + ボーダー色変化 | 150ms |
| ArticleCard | 表示 | Staggered fade-in（グリッド順） | 300ms |
| TagBadge | ホバー | 背景色変化 | 100ms |
| ヒーロー | 初回表示 | タイプライター風テキスト表示 | 1500ms |
| ヘッダー | スクロール | `backdrop-blur`適用 | 200ms |
| ブックマーク | クリック | スケールバウンス `1→1.2→1` | 300ms |
| Toast | 表示/非表示 | スライドイン/アウト（右下から） | 200ms |
| Command palette | 開閉 | `scale(0.95→1)` + `opacity` | 150ms |
| Skeleton | ローディング | シマーアニメーション | 1500ms loop |

### CSS設定

```css
/* Tailwind: transition-all duration-150 ease-out */
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## 9. アクセシビリティ方針

### WCAG 2.1 AA準拠

| 項目 | 対応方針 |
|------|----------|
| **色コントラスト** | テキスト/背景のコントラスト比 4.5:1以上を確保。`#a1a1aa` on `#0a0a0a` = 7.5:1 |
| **キーボード操作** | すべてのインタラクティブ要素にフォーカスリング（`ring-2 ring-amber-500`） |
| **スクリーンリーダー** | 適切な`aria-label`、`aria-live`リージョン（Toast、検索結果数） |
| **見出し階層** | h1→h2→h3の正しい順序を維持 |
| **画像alt** | OGP画像に記事タイトルをalt属性として設定 |
| **フォーカス管理** | モーダル開閉時のフォーカストラップ、ページ遷移時のフォーカスリセット |
| **ランドマーク** | `<header>`, `<main>`, `<nav>`, `<aside>`, `<footer>`を適切に使用 |
| **Skip Link** | ページ上部に「メインコンテンツへスキップ」リンク |
| **リンク** | リンクテキストは内容を説明するもの（「こちら」を避ける） |
| **フォーム** | Input要素にlabelを関連付け。エラーは`aria-describedby`で通知 |

---

## 10. レイアウトコンポーネントツリー

```
<RootLayout>                           # app/layout.tsx
├── <Header />                         # 固定ヘッダー
│   ├── <Logo />
│   ├── <SearchTrigger />              # Cmd+K起動ボタン
│   ├── <BookmarkLink />
│   └── <UserMenu />                   # DropdownMenu
├── <CommandPalette />                 # グローバル検索モーダル
├── <main>
│   ├── <ContentArea />                # 各ページのコンテンツ
│   └── <Sidebar />                    # lg以上で表示
└── <Footer />

# トップページ
<HomePage>
├── <HeroSection />                    # ターミナル風ヒーロー
├── <ArticleSection title="新着記事">
│   └── <ArticleGrid articles={recent} />
├── <ArticleSection title="人気記事">
│   └── <ArticleGrid articles={popular} />
└── <TagCloudSection />

# 記事詳細ページ
<ArticlePage>
├── <Breadcrumb />
├── <OGPImage />
├── <ArticleMeta />                    # タグ、日付、読了時間
├── <ArticleBody />                    # Markdownレンダリング
├── <ArticleActions />                 # ブックマーク、シェア
└── <RelatedArticles />

# タグ別一覧
<TagPage>
├── <Breadcrumb />
├── <TagHeader />                      # タグ名、件数、説明
├── <SortSelect />
├── <ArticleGrid />
└── <Pagination />
```
