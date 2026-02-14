import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Foresee - RSS Reader",
  description: "A simple RSS feed reader",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body className="bg-gray-50 text-gray-900 min-h-screen">
        <header className="bg-white border-b border-gray-200">
          <div className="max-w-4xl mx-auto px-4 py-4">
            <a href="/feeds" className="text-xl font-bold text-blue-600">
              Foresee
            </a>
          </div>
        </header>
        <main className="max-w-4xl mx-auto px-4 py-8">{children}</main>
      </body>
    </html>
  );
}
