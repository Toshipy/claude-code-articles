import Link from "next/link";
import { api } from "@/lib/api";
import { ArticleGrid } from "@/components/ArticleGrid";
import { TagCloud } from "@/components/TagCloud";
import { Sidebar } from "@/components/Sidebar";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ArrowRight, Terminal } from "lucide-react";

export default async function HomePage() {
  const [articlesRes, tagsRes] = await Promise.all([
    api.getArticles({ per_page: 6 }).catch(() => null),
    api.getTags().catch(() => null),
  ]);

  const articles = articlesRes?.data ?? [];
  const tags = tagsRes?.data ?? [];

  // Split into recent and popular (use first 6 as recent, could be separate API call for popular)
  const recentArticles = articles.slice(0, 6);
  const sidebarArticles = articles.slice(0, 5);

  return (
    <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      {/* Hero Section */}
      <section className="mb-12 rounded-lg border border-[#262626] bg-[#141414] p-8 sm:p-12">
        <div className="flex items-center gap-2 text-amber-500 mb-4">
          <Terminal className="h-5 w-5" />
          <span className="font-mono text-sm">~/claude-code-articles</span>
        </div>
        <h1 className="text-3xl font-bold leading-tight sm:text-4xl">
          <span className="text-amber-500">&gt;</span> ClaudeCodeを使いこなすための
          <br />
          <span className="ml-4">技術記事プラットフォーム</span>
          <span className="animate-pulse text-amber-500">_</span>
        </h1>
        <p className="mt-4 max-w-2xl text-[#a1a1aa]">
          Zenn・Qiita・note・dev.toなど、各メディアに散らばるClaude Code関連の技術記事を自動収集。
          最新のTips、活用事例、アップデート情報をまとめてチェックできます。
        </p>
        <div className="mt-6 flex flex-wrap gap-3">
          <Button asChild>
            <Link href="/search">記事を探す</Link>
          </Button>
          <Button variant="outline" asChild>
            <Link href="#tags">
              人気のタグ <ArrowRight className="ml-1 h-4 w-4" />
            </Link>
          </Button>
        </div>
      </section>

      <div className="flex gap-8">
        {/* Main Content */}
        <div className="min-w-0 flex-1">
          {/* Recent Articles */}
          <section className="mb-12">
            <div className="mb-6 flex items-center justify-between">
              <h2 className="text-xl font-bold">新着記事</h2>
              <Link
                href="/search"
                className="text-sm text-amber-500 hover:text-amber-400 transition-colors"
              >
                もっと見る →
              </Link>
            </div>
            {recentArticles.length > 0 ? (
              <ArticleGrid articles={recentArticles} />
            ) : (
              <p className="text-[#a1a1aa]">記事を読み込み中...</p>
            )}
          </section>

          <Separator className="mb-12" />

          {/* Tag Cloud */}
          <section id="tags" className="mb-12">
            <h2 className="mb-6 text-xl font-bold">タグクラウド</h2>
            {tags.length > 0 ? (
              <TagCloud tags={tags} />
            ) : (
              <p className="text-[#a1a1aa]">タグを読み込み中...</p>
            )}
          </section>
        </div>

        {/* Sidebar */}
        <Sidebar popularArticles={sidebarArticles} tags={tags.slice(0, 10)} />
      </div>
    </div>
  );
}
