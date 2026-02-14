"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { syncFeeds } from "@/app/lib/api";

export function SyncButton() {
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  async function handleSync() {
    setLoading(true);
    try {
      await syncFeeds();
      router.refresh();
    } catch {
      // Silently handle - user can retry
    } finally {
      setLoading(false);
    }
  }

  return (
    <button
      onClick={handleSync}
      disabled={loading}
      className="bg-gray-100 text-gray-700 px-4 py-2 rounded-lg text-sm font-medium hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
    >
      {loading ? "Syncing..." : "Sync All"}
    </button>
  );
}
