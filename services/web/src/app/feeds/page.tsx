import Link from "next/link";
import { listFeeds } from "@/app/lib/api";
import { AddFeedForm } from "./add-feed-form";
import { SyncButton } from "./sync-button";

export const dynamic = "force-dynamic";

export default async function FeedsPage() {
  let feeds: Awaited<ReturnType<typeof listFeeds>> = [];
  let error: string | null = null;

  try {
    feeds = await listFeeds();
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load feeds";
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Feeds</h1>
        <SyncButton />
      </div>

      <AddFeedForm />

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 rounded-lg p-4 mb-6">
          {error}
        </div>
      )}

      {feeds.length === 0 && !error ? (
        <p className="text-gray-500 text-center py-12">
          No feeds registered yet. Add one above.
        </p>
      ) : (
        <div className="space-y-3">
          {feeds.map((feed) => (
            <Link
              key={feed.id}
              href={`/feeds/${feed.id}`}
              className="block bg-white border border-gray-200 rounded-lg p-4 hover:border-blue-300 hover:shadow-sm transition-all"
            >
              <h2 className="font-semibold text-lg">
                {feed.title || "Untitled"}
              </h2>
              <p className="text-sm text-gray-500 mt-1 truncate">{feed.url}</p>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
