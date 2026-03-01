import type { Metadata } from "next";
import { Header } from "@/components/Header";
import "./globals.css";

export const metadata: Metadata = {
  title: "ClaudeCode Articles - Claude Codeの技術記事プラットフォーム",
  description:
    "Claude Codeに関する最新技術記事を集約。開発者が効率的に情報収集できる専用プラットフォーム。",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja" className="dark">
      <body className="min-h-screen bg-[#0a0a0a] text-[#fafafa] antialiased">
        <a
          href="#main-content"
          className="sr-only focus:not-sr-only focus:absolute focus:z-[100] focus:bg-amber-500 focus:text-black focus:px-4 focus:py-2"
        >
          メインコンテンツへスキップ
        </a>
        <Header />
        <main id="main-content">{children}</main>
        <footer className="border-t border-[#262626] py-8 text-center text-sm text-[#a1a1aa]">
          <div className="mx-auto max-w-7xl px-4">
            &copy; 2026 ClaudeCode Articles
          </div>
        </footer>
      </body>
    </html>
  );
}
