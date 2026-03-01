import Link from "next/link";
import { notFound } from "next/navigation";
import { api } from "@/lib/api";
import { ArticleGrid } from "@/components/ArticleGrid";
import { Button } from "@/components/ui/button";
import { ChevronRight, Home } from "lucide-react";

type Props = {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ page?: string; sort?: string; order?: string }>;
};

export default async function TagPage({ params, searchParams }: Props) {
  const { slug } = await params;
  const sp = await searchParams;
  const page = Number(sp.page) || 1;
  const sort = sp.sort ?? "published_at";
  const order = sp.order ?? "desc";

  const res = await api.getTagArticles(slug, { page, sort, order }).catch(() => null);
  if (!res) notFound();

  const articles = res.data;
  const pagination = res.pagination;

  // Decode tag name from slug for display
  const tagName = slug.replace(/-/g, " ");

  return (
    <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      {/* Breadcrumb */}
      <nav aria-label="パンくずリスト" className="mb-6 flex items-center gap-1 text-sm text-[#a1a1aa]">
        <Link href="/" className="hover:text-amber-500 transition-colors">
          <Home className="h-4 w-4" />
        </Link>
        <ChevronRight className="h-3 w-3" />
        <span>Tags</span>
        <ChevronRight className="h-3 w-3" />
        <span className="text-[#fafafa]">{tagName}</span>
      </nav>

      {/* Header */}
      <div className="mb-8">
        <h1 className="text-2xl font-bold sm:text-3xl">
          {tagName} の記事
          <span className="ml-2 text-lg text-[#a1a1aa]">
            ({pagination.total}件)
          </span>
        </h1>
        <p className="mt-2 text-[#a1a1aa]">
          {tagName}に関する技術記事の一覧です
        </p>
      </div>

      {/* Sort Controls */}
      <div className="mb-6 flex items-center gap-2">
        <span className="text-sm text-[#a1a1aa]">並び替え:</span>
        <div className="flex gap-1">
          <Button
            variant={sort === "published_at" ? "secondary" : "ghost"}
            size="sm"
            asChild
          >
            <Link href={`/tags/${slug}?sort=published_at&order=desc`}>
              新着順
            </Link>
          </Button>
          <Button
            variant={sort === "collected_at" ? "secondary" : "ghost"}
            size="sm"
            asChild
          >
            <Link href={`/tags/${slug}?sort=collected_at&order=desc`}>
              収集順
            </Link>
          </Button>
        </div>
      </div>

      {/* Article Grid */}
      {articles.length > 0 ? (
        <ArticleGrid articles={articles} />
      ) : (
        <p className="text-[#a1a1aa]">このタグに関する記事はまだありません。</p>
      )}

      {/* Pagination */}
      {pagination.total_pages > 1 && (
        <nav aria-label="ページネーション" className="mt-8 flex items-center justify-center gap-2">
          {page > 1 && (
            <Button variant="outline" size="sm" asChild>
              <Link href={`/tags/${slug}?page=${page - 1}&sort=${sort}&order=${order}`}>
                &lt; 前へ
              </Link>
            </Button>
          )}
          {Array.from({ length: pagination.total_pages }, (_, i) => i + 1).map(
            (p) => (
              <Button
                key={p}
                variant={p === page ? "default" : "outline"}
                size="sm"
                asChild
              >
                <Link href={`/tags/${slug}?page=${p}&sort=${sort}&order=${order}`}>
                  {p}
                </Link>
              </Button>
            ),
          )}
          {page < pagination.total_pages && (
            <Button variant="outline" size="sm" asChild>
              <Link href={`/tags/${slug}?page=${page + 1}&sort=${sort}&order=${order}`}>
                次へ &gt;
              </Link>
            </Button>
          )}
        </nav>
      )}
    </div>
  );
}
