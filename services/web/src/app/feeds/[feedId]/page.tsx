import Link from "next/link";
import { listArticles } from "@/app/lib/api";

export const dynamic = "force-dynamic";

export default async function FeedArticlesPage({
  params,
}: {
  params: Promise<{ feedId: string }>;
}) {
  const { feedId } = await params;
  let articles: Awaited<ReturnType<typeof listArticles>> = [];
  let error: string | null = null;

  try {
    articles = await listArticles(feedId);
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load articles";
  }

  return (
    <div>
      <div className="mb-6">
        <Link
          href="/feeds"
          className="text-blue-600 hover:text-blue-800 text-sm font-medium"
        >
          &larr; Back to Feeds
        </Link>
      </div>

      <h1 className="text-2xl font-bold mb-6">Articles</h1>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 rounded-lg p-4 mb-6">
          {error}
        </div>
      )}

      {articles.length === 0 && !error ? (
        <p className="text-gray-500 text-center py-12">
          No articles found. Try syncing your feeds.
        </p>
      ) : (
        <div className="space-y-3">
          {articles.map((article) => (
            <a
              key={article.id}
              href={article.url}
              target="_blank"
              rel="noopener noreferrer"
              className="block bg-white border border-gray-200 rounded-lg p-4 hover:border-blue-300 hover:shadow-sm transition-all"
            >
              <h2 className="font-semibold">{article.title || "Untitled"}</h2>
              {article.publishedAt && (
                <p className="text-sm text-gray-500 mt-1">
                  {new Date(article.publishedAt).toLocaleDateString("ja-JP", {
                    year: "numeric",
                    month: "long",
                    day: "numeric",
                  })}
                </p>
              )}
            </a>
          ))}
        </div>
      )}
    </div>
  );
}
