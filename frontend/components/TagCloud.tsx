import type { Tag } from "@/lib/types";
import { TagBadge } from "./TagBadge";

type TagCloudProps = {
  tags: Tag[];
};

export function TagCloud({ tags }: TagCloudProps) {
  return (
    <div className="flex flex-wrap gap-2">
      {tags.map((tag) => (
        <TagBadge
          key={tag.id}
          tag={{ slug: tag.slug, name: tag.name, count: tag.article_count }}
          size="md"
        />
      ))}
    </div>
  );
}
