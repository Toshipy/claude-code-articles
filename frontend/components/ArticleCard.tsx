import Image from "next/image";
import Link from "next/link";
import type { Tag } from "@/lib/types";
import { Card } from "@/components/ui/card";
import { TagBadge } from "./TagBadge";
import { formatDate } from "@/lib/utils";
import { cn } from "@/lib/utils";
import { Bookmark } from "lucide-react";

type ArticleCardProps = {
  article: {
    id: string;
    title: string;
    summary: string;
    thumbnail_url: string;
    tags: Tag[];
    published_at: string;
    source: { name: string; icon_url: string };
  };
  variant?: "default" | "compact" | "horizontal";
};

export function ArticleCard({ article, variant = "default" }: ArticleCardProps) {
  if (variant === "compact") {
    return (
      <Link href={`/articles/${article.id}`} className="group block">
        <div className="flex items-start gap-3 rounded-md p-2 transition-colors hover:bg-[#1a1a1a]">
          <div className="min-w-0 flex-1">
            <h3 className="line-clamp-2 text-sm font-medium text-[#fafafa] group-hover:text-amber-500 transition-colors">
              {article.title}
            </h3>
            <div className="mt-1 flex items-center gap-2 text-xs text-[#a1a1aa]">
              <span>{article.source.name}</span>
              <span>{formatDate(article.published_at)}</span>
            </div>
          </div>
        </div>
      </Link>
    );
  }

  if (variant === "horizontal") {
    return (
      <Link href={`/articles/${article.id}`} className="group block">
        <Card className="overflow-hidden transition-all duration-150 hover:border-amber-500/50 hover:-translate-y-0.5">
          <div className="flex">
            <div className="relative h-32 w-48 shrink-0">
              <Image
                src={article.thumbnail_url}
                alt={article.title}
                fill
                className="object-cover"
                sizes="192px"
              />
            </div>
            <div className="flex flex-1 flex-col justify-between p-4">
              <div>
                <h3 className="line-clamp-2 font-semibold text-[#fafafa] group-hover:text-amber-500 transition-colors">
                  {article.title}
                </h3>
                <p className="mt-1 line-clamp-2 text-sm text-[#a1a1aa]">
                  {article.summary}
                </p>
              </div>
              <div className="mt-2 flex items-center gap-2">
                <div className="flex flex-wrap gap-1">
                  {article.tags.slice(0, 3).map((tag) => (
                    <TagBadge key={tag.id} tag={tag} size="sm" />
                  ))}
                </div>
                <span className="ml-auto text-xs text-[#a1a1aa]">
                  {formatDate(article.published_at)}
                </span>
              </div>
            </div>
          </div>
        </Card>
      </Link>
    );
  }

  // default variant
  return (
    <Link href={`/articles/${article.id}`} className="group block">
      <Card className={cn(
        "overflow-hidden transition-all duration-150",
        "hover:border-amber-500/50 hover:-translate-y-0.5",
      )}>
        <div className="relative aspect-video w-full overflow-hidden">
          <Image
            src={article.thumbnail_url}
            alt={article.title}
            fill
            className="object-cover transition-transform duration-300 group-hover:scale-105"
            sizes="(max-width: 768px) 100vw, (max-width: 1024px) 50vw, 33vw"
          />
        </div>
        <div className="p-4">
          <h3 className="line-clamp-2 text-base font-semibold leading-snug text-[#fafafa] group-hover:text-amber-500 transition-colors">
            {article.title}
          </h3>
          <p className="mt-2 line-clamp-2 text-sm text-[#a1a1aa]">
            {article.summary}
          </p>
          <div className="mt-3 flex flex-wrap gap-1">
            {article.tags.slice(0, 3).map((tag) => (
              <TagBadge key={tag.id} tag={tag} size="sm" />
            ))}
          </div>
          <div className="mt-3 flex items-center justify-between text-xs text-[#a1a1aa]">
            <div className="flex items-center gap-2">
              <Image
                src={article.source.icon_url}
                alt={article.source.name}
                width={16}
                height={16}
                className="rounded"
              />
              <span>{article.source.name}</span>
              <span>{formatDate(article.published_at)}</span>
            </div>
            <Bookmark className="h-4 w-4 opacity-0 transition-opacity group-hover:opacity-100" />
          </div>
        </div>
      </Card>
    </Link>
  );
}
