import type { ArticleSummary, Tag } from "@/lib/types";
import { ArticleGrid } from "./ArticleGrid";
import { TagCloud } from "./TagCloud";
import { Separator } from "@/components/ui/separator";

type SidebarProps = {
  popularArticles?: ArticleSummary[];
  tags?: Tag[];
};

export function Sidebar({ popularArticles, tags }: SidebarProps) {
  return (
    <aside className="hidden lg:block w-80 shrink-0 sticky top-20 self-start space-y-6">
      {popularArticles && popularArticles.length > 0 && (
        <div>
          <h2 className="mb-3 text-sm font-semibold uppercase tracking-wider text-[#a1a1aa]">
            人気記事
          </h2>
          <ArticleGrid articles={popularArticles} variant="compact" />
        </div>
      )}

      {popularArticles && tags && <Separator />}

      {tags && tags.length > 0 && (
        <div>
          <h2 className="mb-3 text-sm font-semibold uppercase tracking-wider text-[#a1a1aa]">
            タグ
          </h2>
          <TagCloud tags={tags} />
        </div>
      )}
    </aside>
  );
}
