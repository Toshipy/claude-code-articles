import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

type TagBadgeProps = {
  tag: {
    slug: string;
    name: string;
    count?: number;
  };
  size?: "sm" | "md";
  interactive?: boolean;
};

export function TagBadge({ tag, size = "sm", interactive = true }: TagBadgeProps) {
  const classes = cn(
    "border-amber-500/30 text-amber-500 hover:bg-amber-500/10 transition-colors cursor-pointer",
    size === "md" && "px-3 py-1 text-sm",
  );

  if (interactive) {
    return (
      <Link href={`/tags/${tag.slug}`}>
        <Badge variant="outline" className={classes}>
          {tag.name}
          {tag.count !== undefined && (
            <span className="ml-1 text-[#a1a1aa]">({tag.count})</span>
          )}
        </Badge>
      </Link>
    );
  }

  return (
    <Badge variant="outline" className={classes}>
      {tag.name}
      {tag.count !== undefined && (
        <span className="ml-1 text-[#a1a1aa]">({tag.count})</span>
      )}
    </Badge>
  );
}
