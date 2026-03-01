"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, useEffect, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Search,
  Bookmark,
  Menu,
  X,
  Terminal,
} from "lucide-react";

export function Header() {
  const router = useRouter();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [searchOpen, setSearchOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 10);
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();
        setSearchOpen(true);
      }
      if (e.key === "Escape") {
        setSearchOpen(false);
      }
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, []);

  const handleSearch = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      const trimmed = searchQuery.trim();
      if (trimmed.length >= 2) {
        router.push(`/search?q=${encodeURIComponent(trimmed)}`);
        setSearchOpen(false);
        setSearchQuery("");
      }
    },
    [searchQuery, router],
  );

  return (
    <>
      <header
        className={`sticky top-0 z-50 w-full border-b border-[#262626] transition-all duration-200 ${
          scrolled
            ? "bg-[#0a0a0a]/80 backdrop-blur-md"
            : "bg-[#0a0a0a]"
        }`}
      >
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2 text-[#fafafa]">
            <Terminal className="h-5 w-5 text-amber-500" />
            <span className="text-lg font-bold">
              ClaudeCode <span className="hidden sm:inline">Articles</span>
            </span>
            <span className="animate-pulse text-amber-500">_</span>
          </Link>

          {/* Desktop Search Trigger */}
          <button
            onClick={() => setSearchOpen(true)}
            className="hidden md:flex items-center gap-2 rounded-md border border-[#262626] bg-[#141414] px-3 py-1.5 text-sm text-[#a1a1aa] hover:border-[#404040] transition-colors"
          >
            <Search className="h-4 w-4" />
            <span>検索...</span>
            <kbd className="ml-4 rounded border border-[#262626] bg-[#1a1a1a] px-1.5 py-0.5 text-xs">
              ⌘K
            </kbd>
          </button>

          {/* Desktop Right Actions */}
          <div className="hidden md:flex items-center gap-2">
            <Button variant="ghost" size="icon" asChild>
              <Link href="/bookmarks" aria-label="ブックマーク">
                <Bookmark className="h-5 w-5" />
              </Link>
            </Button>
            <Button variant="outline" size="sm">
              ログイン
            </Button>
          </div>

          {/* Mobile Menu Toggle */}
          <button
            className="md:hidden text-[#fafafa]"
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            aria-label="メニュー"
          >
            {mobileMenuOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
          </button>
        </div>

        {/* Mobile Menu */}
        {mobileMenuOpen && (
          <div className="border-t border-[#262626] bg-[#0a0a0a] px-4 py-4 md:hidden">
            <nav className="flex flex-col gap-3">
              <button
                onClick={() => {
                  setMobileMenuOpen(false);
                  setSearchOpen(true);
                }}
                className="flex items-center gap-2 rounded-md border border-[#262626] bg-[#141414] px-3 py-2 text-sm text-[#a1a1aa]"
              >
                <Search className="h-4 w-4" />
                記事を検索...
              </button>
              <Link
                href="/bookmarks"
                className="flex items-center gap-2 px-3 py-2 text-sm text-[#fafafa] hover:text-amber-500"
                onClick={() => setMobileMenuOpen(false)}
              >
                <Bookmark className="h-4 w-4" />
                ブックマーク
              </Link>
              <Button variant="outline" size="sm" className="w-full">
                ログイン
              </Button>
            </nav>
          </div>
        )}
      </header>

      {/* Search Modal */}
      {searchOpen && (
        <div
          className="fixed inset-0 z-[60] flex items-start justify-center bg-black/60 backdrop-blur-sm pt-[20vh]"
          onClick={() => setSearchOpen(false)}
        >
          <div
            className="w-full max-w-lg mx-4 rounded-lg border border-[#262626] bg-[#141414] p-4 shadow-2xl"
            onClick={(e) => e.stopPropagation()}
          >
            <form onSubmit={handleSearch}>
              <div className="relative">
                <Search className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-[#a1a1aa]" />
                <Input
                  autoFocus
                  type="search"
                  placeholder="記事を検索..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="h-12 pl-11 text-base"
                  aria-label="記事を検索"
                />
              </div>
            </form>
            <p className="mt-2 text-xs text-[#a1a1aa]">
              Enterで検索 ・ Escで閉じる
            </p>
          </div>
        </div>
      )}
    </>
  );
}
