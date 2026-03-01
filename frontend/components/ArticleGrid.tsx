import type { ArticleSummary } from "@/lib/types";
import { ArticleCard } from "./ArticleCard";
import { cn } from "@/lib/utils";

type ArticleGridProps = {
  articles: ArticleSummary[];
  variant?: "default" | "compact";
  columns?: 2 | 3;
};

export function ArticleGrid({ articles, variant = "default", columns = 3 }: ArticleGridProps) {
  if (variant === "compact") {
    return (
      <div className="space-y-1">
        {articles.map((article) => (
          <ArticleCard key={article.id} article={article} variant="compact" />
        ))}
      </div>
    );
  }

  return (
    <div
      className={cn(
        "grid gap-6",
        "grid-cols-1 md:grid-cols-2",
        columns === 3 && "lg:grid-cols-3",
      )}
    >
      {articles.map((article) => (
        <ArticleCard key={article.id} article={article} />
      ))}
    </div>
  );
}
