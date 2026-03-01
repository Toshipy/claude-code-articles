"use client";

import { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { Input } from "@/components/ui/input";
import { Search, X } from "lucide-react";

type SearchBarProps = {
  defaultValue?: string;
};

export function SearchBar({ defaultValue = "" }: SearchBarProps) {
  const [query, setQuery] = useState(defaultValue);
  const router = useRouter();

  const handleSubmit = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      const trimmed = query.trim();
      if (trimmed.length >= 2) {
        router.push(`/search?q=${encodeURIComponent(trimmed)}`);
      }
    },
    [query, router],
  );

  const handleClear = useCallback(() => {
    setQuery("");
  }, []);

  return (
    <form onSubmit={handleSubmit} className="relative w-full">
      <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-[#a1a1aa]" />
      <Input
        type="search"
        placeholder="記事を検索..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        className="pl-10 pr-10"
        aria-label="記事を検索"
      />
      {query && (
        <button
          type="button"
          onClick={handleClear}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-[#a1a1aa] hover:text-[#fafafa]"
          aria-label="検索をクリア"
        >
          <X className="h-4 w-4" />
        </button>
      )}
    </form>
  );
}
