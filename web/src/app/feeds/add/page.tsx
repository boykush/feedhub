"use client";

import { useState } from "react";
import { addFeed } from "../../lib/api";

export default function AddFeedPage() {
	const [url, setUrl] = useState("");
	const [loading, setLoading] = useState(false);
	const [result, setResult] = useState<{
		feedId: string;
		title: string;
	} | null>(null);
	const [error, setError] = useState<string | null>(null);

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setLoading(true);
		setResult(null);
		setError(null);

		try {
			const data = await addFeed(url);
			setResult(data);
			setUrl("");
		} catch (err) {
			setError(err instanceof Error ? err.message : "An error occurred");
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className="max-w-2xl mx-auto">
			<h1 className="text-2xl font-bold mb-6">Add Feed</h1>

			<form
				onSubmit={handleSubmit}
				className="bg-white rounded-lg border border-gray-200 p-6"
			>
				<div className="mb-4">
					<label
						htmlFor="feed-url"
						className="block text-sm font-medium text-gray-700 mb-2"
					>
						Feed URL
					</label>
					<input
						id="feed-url"
						type="url"
						value={url}
						onChange={(e) => setUrl(e.target.value)}
						placeholder="https://example.com/feed.xml"
						required
						className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					/>
				</div>
				<button
					type="submit"
					disabled={loading || !url}
					className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
				>
					{loading ? "Adding..." : "Add Feed"}
				</button>
			</form>

			{result && (
				<div className="mt-4 bg-green-50 border border-green-200 rounded-lg p-4">
					<p className="font-medium text-green-800">Feed added successfully</p>
					<p className="text-sm text-green-700 mt-1">Title: {result.title}</p>
					<p className="text-sm text-green-700 font-mono">
						ID: {result.feedId}
					</p>
				</div>
			)}

			{error && (
				<div className="mt-4 bg-red-50 border border-red-200 rounded-lg p-4">
					<p className="font-medium text-red-800">Failed to add feed</p>
					<p className="text-sm text-red-700 mt-1">{error}</p>
				</div>
			)}
		</div>
	);
}
