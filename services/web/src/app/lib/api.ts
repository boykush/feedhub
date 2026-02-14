const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export interface Feed {
  id: string;
  url: string;
  title: string;
}

export interface Article {
  id: string;
  feedId: string;
  title: string;
  url: string;
  publishedAt: string;
}

export async function listFeeds(): Promise<Feed[]> {
  const res = await fetch(`${API_BASE_URL}/api/v1/feeds`, {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error(`Failed to fetch feeds: ${res.status}`);
  }
  const data = await res.json();
  return data.feeds || [];
}

export async function listArticles(feedId: string): Promise<Article[]> {
  const res = await fetch(
    `${API_BASE_URL}/api/v1/articles?feedId=${encodeURIComponent(feedId)}`,
    { cache: "no-store" }
  );
  if (!res.ok) {
    throw new Error(`Failed to fetch articles: ${res.status}`);
  }
  const data = await res.json();
  return data.articles || [];
}

export async function addFeed(url: string): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/api/v1/collector/feeds`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ url }),
  });
  if (!res.ok) {
    throw new Error(`Failed to add feed: ${res.status}`);
  }
}

export async function syncFeeds(): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/api/v1/collector/sync`, {
    method: "POST",
  });
  if (!res.ok) {
    throw new Error(`Failed to sync feeds: ${res.status}`);
  }
}
