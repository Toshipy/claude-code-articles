"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ArticleCard } from "@/components/ArticleCard";
import { Bookmark, LogIn } from "lucide-react";

// In a real app, auth state would come from a context/provider.
// For now, show the unauthenticated state.
export default function BookmarksPage() {
  const isLoggedIn = false;
  const bookmarks: never[] = [];

  if (!isLoggedIn) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-16 sm:px-6 text-center">
        <Bookmark className="mx-auto mb-4 h-12 w-12 text-[#a1a1aa]" />
        <h1 className="text-2xl font-bold">ブックマーク</h1>
        <p className="mt-3 text-[#a1a1aa]">
          ブックマーク機能を使うにはログインが必要です
        </p>
        <div className="mt-6 flex justify-center gap-3">
          <Button>
            <LogIn className="mr-1 h-4 w-4" />
            ログイン
          </Button>
          <Button variant="outline" asChild>
            <Link href="/">トップに戻る</Link>
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-3xl px-4 py-8 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold">
          ブックマーク
          <span className="ml-2 text-lg text-[#a1a1aa]">
            ({bookmarks.length}件)
          </span>
        </h1>
        <p className="mt-2 text-[#a1a1aa]">保存した記事の一覧です</p>
      </div>

      {bookmarks.length > 0 ? (
        <div className="space-y-4">
          {/* Bookmark items would render here */}
        </div>
      ) : (
        <div className="rounded-lg border border-[#262626] bg-[#141414] p-8 text-center">
          <p className="text-[#a1a1aa]">
            まだブックマークした記事がありません。
          </p>
          <Button variant="outline" className="mt-4" asChild>
            <Link href="/">記事を探す</Link>
          </Button>
        </div>
      )}
    </div>
  );
}
