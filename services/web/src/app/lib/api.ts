const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function healthCheck(): Promise<string> {
  const res = await fetch(`${API_BASE_URL}/api/v1/feeds/health`, {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error(`Health check failed: ${res.status}`);
  }
  const data = await res.json();
  return data.status || "UNKNOWN";
}
