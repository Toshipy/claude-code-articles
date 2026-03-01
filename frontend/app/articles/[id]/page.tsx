import Image from "next/image";
import Link from "next/link";
import { notFound } from "next/navigation";
import { api } from "@/lib/api";
import { TagBadge } from "@/components/TagBadge";
import { ArticleCard } from "@/components/ArticleCard";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { formatDate, estimateReadingTime } from "@/lib/utils";
import {
  Bookmark,
  ExternalLink,
  ChevronRight,
  Home,
} from "lucide-react";

type Props = {
  params: Promise<{ id: string }>;
};

export default async function ArticleDetailPage({ params }: Props) {
  const { id } = await params;

  const res = await api.getArticle(id).catch(() => null);
  if (!res) notFound();

  const article = res.data;
  const readingTime = estimateReadingTime(article.content ?? "");

  return (
    <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      {/* Breadcrumb */}
      <nav aria-label="パンくずリスト" className="mb-6 flex items-center gap-1 text-sm text-[#a1a1aa]">
        <Link href="/" className="hover:text-amber-500 transition-colors">
          <Home className="h-4 w-4" />
        </Link>
        <ChevronRight className="h-3 w-3" />
        {article.tags[0] && (
          <>
            <Link
              href={`/tags/${article.tags[0].slug}`}
              className="hover:text-amber-500 transition-colors"
            >
              {article.tags[0].name}
            </Link>
            <ChevronRight className="h-3 w-3" />
          </>
        )}
        <span className="truncate text-[#fafafa]">{article.title}</span>
      </nav>

      <div className="flex gap-8">
        {/* Main Content */}
        <article className="min-w-0 flex-1">
          {/* Thumbnail */}
          {article.thumbnail_url && (
            <div className="relative mb-6 aspect-video w-full overflow-hidden rounded-lg">
              <Image
                src={article.thumbnail_url}
                alt={article.title}
                fill
                className="object-cover"
                priority
                sizes="(max-width: 1024px) 100vw, 800px"
              />
            </div>
          )}

          {/* Tags */}
          <div className="mb-4 flex flex-wrap gap-2">
            {article.tags.map((tag) => (
              <TagBadge key={tag.id} tag={tag} size="md" />
            ))}
          </div>

          {/* Title */}
          <h1 className="text-2xl font-bold leading-tight sm:text-3xl">
            {article.title}
          </h1>

          {/* Meta */}
          <div className="mt-3 flex items-center gap-3 text-sm text-[#a1a1aa]">
            <div className="flex items-center gap-1.5">
              <Image
                src={article.source.icon_url}
                alt={article.source.name}
                width={16}
                height={16}
                className="rounded"
              />
              <span>{article.source.name}</span>
            </div>
            <span>{formatDate(article.published_at)}</span>
            <span>読了 {readingTime}分</span>
          </div>

          <Separator className="my-6" />

          {/* Article Body */}
          <div className="prose prose-invert max-w-none prose-headings:text-[#fafafa] prose-p:text-[#d4d4d8] prose-a:text-amber-500 prose-a:no-underline hover:prose-a:underline prose-code:bg-[#0d1117] prose-code:rounded prose-code:px-1 prose-pre:bg-[#0d1117] prose-pre:border prose-pre:border-[#262626]">
            <p className="text-[#a1a1aa] leading-relaxed whitespace-pre-wrap">
              {article.content ?? article.summary}
            </p>
          </div>

          <Separator className="my-6" />

          {/* Actions */}
          <div className="flex items-center gap-3">
            <Button variant="outline" size="sm">
              <Bookmark className="mr-1 h-4 w-4" />
              ブックマーク
            </Button>
            <Button variant="outline" size="sm" asChild>
              <a
                href={article.url}
                target="_blank"
                rel="noopener noreferrer"
              >
                <ExternalLink className="mr-1 h-4 w-4" />
                元記事を読む
              </a>
            </Button>
          </div>

          {/* Related Articles (mobile) */}
          {article.related_articles && article.related_articles.length > 0 && (
            <section className="mt-12 lg:hidden">
              <h2 className="mb-4 text-lg font-bold">関連記事</h2>
              <div className="grid gap-4 grid-cols-1 sm:grid-cols-2">
                {article.related_articles.map((related) => (
                  <Link
                    key={related.id}
                    href={`/articles/${related.id}`}
                    className="group flex items-center gap-3 rounded-lg border border-[#262626] bg-[#141414] p-3 transition-colors hover:border-amber-500/50"
                  >
                    {related.thumbnail_url && (
                      <div className="relative h-16 w-24 shrink-0 overflow-hidden rounded">
                        <Image
                          src={related.thumbnail_url}
                          alt={related.title}
                          fill
                          className="object-cover"
                          sizes="96px"
                        />
                      </div>
                    )}
                    <span className="text-sm font-medium group-hover:text-amber-500 transition-colors line-clamp-2">
                      {related.title}
                    </span>
                  </Link>
                ))}
              </div>
            </section>
          )}
        </article>

        {/* Desktop Sidebar */}
        <aside className="hidden lg:block w-72 shrink-0 sticky top-20 self-start space-y-6">
          {/* Related Articles */}
          {article.related_articles && article.related_articles.length > 0 && (
            <div>
              <h2 className="mb-3 text-sm font-semibold uppercase tracking-wider text-[#a1a1aa]">
                関連記事
              </h2>
              <div className="space-y-3">
                {article.related_articles.map((related) => (
                  <Link
                    key={related.id}
                    href={`/articles/${related.id}`}
                    className="group flex items-center gap-3 rounded-md p-2 transition-colors hover:bg-[#1a1a1a]"
                  >
                    {related.thumbnail_url && (
                      <div className="relative h-12 w-16 shrink-0 overflow-hidden rounded">
                        <Image
                          src={related.thumbnail_url}
                          alt={related.title}
                          fill
                          className="object-cover"
                          sizes="64px"
                        />
                      </div>
                    )}
                    <span className="text-sm group-hover:text-amber-500 transition-colors line-clamp-2">
                      {related.title}
                    </span>
                  </Link>
                ))}
              </div>
            </div>
          )}
        </aside>
      </div>
    </div>
  );
}
