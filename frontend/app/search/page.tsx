import Link from "next/link";
import { api } from "@/lib/api";
import { ArticleCard } from "@/components/ArticleCard";
import { SearchBar } from "@/components/SearchBar";
import { Button } from "@/components/ui/button";

type Props = {
  searchParams: Promise<{ q?: string; page?: string }>;
};

export default async function SearchPage({ searchParams }: Props) {
  const sp = await searchParams;
  const query = sp.q ?? "";
  const page = Number(sp.page) || 1;

  let articles: Awaited<ReturnType<typeof api.searchArticles>> | null = null;
  if (query.length >= 2) {
    articles = await api.searchArticles({ q: query, page }).catch(() => null);
  }

  return (
    <div className="mx-auto max-w-3xl px-4 py-8 sm:px-6">
      {/* Search Input */}
      <div className="mb-8">
        <SearchBar defaultValue={query} />
      </div>

      {/* Results */}
      {query.length >= 2 && articles ? (
        <>
          <p className="mb-6 text-sm text-[#a1a1aa]">
            &ldquo;{query}&rdquo; の検索結果:{" "}
            <span className="text-[#fafafa]">{articles.pagination.total}件</span>
          </p>

          {articles.data.length > 0 ? (
            <div className="space-y-4">
              {articles.data.map((article) => (
                <ArticleCard
                  key={article.id}
                  article={article}
                  variant="horizontal"
                />
              ))}
            </div>
          ) : (
            <div className="rounded-lg border border-[#262626] bg-[#141414] p-8 text-center">
              <p className="text-[#a1a1aa]">
                該当する記事が見つかりませんでした。
              </p>
              <p className="mt-2 text-sm text-[#a1a1aa]">
                別のキーワードで検索してみてください。
              </p>
            </div>
          )}

          {/* Pagination */}
          {articles.pagination.total_pages > 1 && (
            <nav aria-label="ページネーション" className="mt-8 flex items-center justify-center gap-2">
              {page > 1 && (
                <Button variant="outline" size="sm" asChild>
                  <Link href={`/search?q=${encodeURIComponent(query)}&page=${page - 1}`}>
                    &lt; 前へ
                  </Link>
                </Button>
              )}
              {Array.from(
                { length: articles.pagination.total_pages },
                (_, i) => i + 1,
              ).map((p) => (
                <Button
                  key={p}
                  variant={p === page ? "default" : "outline"}
                  size="sm"
                  asChild
                >
                  <Link href={`/search?q=${encodeURIComponent(query)}&page=${p}`}>
                    {p}
                  </Link>
                </Button>
              ))}
              {page < articles.pagination.total_pages && (
                <Button variant="outline" size="sm" asChild>
                  <Link href={`/search?q=${encodeURIComponent(query)}&page=${page + 1}`}>
                    次へ &gt;
                  </Link>
                </Button>
              )}
            </nav>
          )}
        </>
      ) : query ? (
        <p className="text-sm text-[#a1a1aa]">
          検索キーワードは2文字以上入力してください。
        </p>
      ) : (
        <div className="rounded-lg border border-[#262626] bg-[#141414] p-8 text-center">
          <p className="text-lg font-medium">記事を検索</p>
          <p className="mt-2 text-sm text-[#a1a1aa]">
            キーワードを入力して、Claude Code関連の技術記事を探しましょう。
          </p>
        </div>
      )}
    </div>
  );
}
