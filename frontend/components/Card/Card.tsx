"use client";

import { useState } from "react";

// Use Next's same-origin proxy so browser requests never require CORS.
const API_BASE_URL = "/api";

interface ShortenResponse {
  short_code: string;
  short_url: string;
  original_url: string;
}

const Card = () => {
  const [url, setUrl] = useState("");
  const [result, setResult] = useState<ShortenResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url.trim()) return;

    setLoading(true);
    setError(null);
    setResult(null);
    setCopied(false);

    try {
      const response = await fetch(`${API_BASE_URL}/shorten`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url: url.trim() }),
      });

      const responseBody = await response.text();
      let data: ShortenResponse | { error?: string } | null = null;

      if (responseBody) {
        try {
          data = JSON.parse(responseBody);
        } catch {
          // A proxy or server error can be HTML/text instead of the API's JSON.
        }
      }

      if (!response.ok) {
        const apiError = data && "error" in data ? data.error : undefined;
        throw new Error(
          apiError ||
            `Request failed (${response.status}). Is the Go backend running on port 8080?`,
        );
      }

      if (!data || !("short_code" in data)) {
        throw new Error("The backend returned an invalid response.");
      }

      setResult(data);
      setUrl("");
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Something went wrong";
      setError(message || "Something went wrong. Is the backend running?");
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = async () => {
    if (!result) return;
    try {
      await navigator.clipboard.writeText(result.short_url);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      setError("Failed to copy to clipboard");
    }
  };

  return (
    <div className="w-full max-w-2xl mt-8 sm:mt-12 p-5 sm:p-8 bg-white rounded-2xl shadow-xl border border-neutral-200">
      <h2 className="text-xl font-semibold text-center text-neutral-800">
        Paste the URL to be shortened
      </h2>

      <form
        onSubmit={handleSubmit}
        className="mt-6 flex flex-col sm:flex-row gap-3"
      >
        <input
          type="url"
          required
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://example.com/very-long-url"
          className="w-full flex-1 px-5 py-4 text-neutral-700 border border-neutral-300 rounded-xl outline-none focus:ring-2 focus:ring-neutral-900 focus:border-transparent transition-all"
        />

        <button
          type="submit"
          disabled={loading}
          className="w-full sm:w-auto px-6 py-4 bg-neutral-900 text-white font-medium 
            rounded-xl hover:bg-neutral-800 transition-colors cursor-pointer disabled:opacity-60 disabled:cursor-not-allowed"
        >
          {loading ? "Shortening..." : "Shorten URL"}
        </button>
      </form>

      {/* Error State */}
      {error && (
        <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-xl text-sm text-red-600 text-center">
          {error}
        </div>
      )}

      {/* Success Result Box */}
      {result ? (
        <div className="mt-6 p-4 bg-neutral-50 rounded-xl border border-neutral-200 flex flex-col sm:flex-row items-center justify-between gap-3">
          <div className="w-full sm:w-auto truncate text-center sm:text-left">
            <span className="block text-xs font-semibold text-neutral-400 uppercase tracking-wider">
              Shortened Link
            </span>
            <a
              href={result.short_url}
              target="_blank"
              rel="noreferrer"
              className="text-neutral-900 font-medium hover:underline break-all"
            >
              {result.short_url}
            </a>
          </div>

          <button
            onClick={handleCopy}
            className={`w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg transition-all ${
              copied
                ? "bg-emerald-600 text-white"
                : "bg-white border border-neutral-300 text-neutral-700 hover:bg-neutral-100"
            }`}
          >
            {copied ? "Copied!" : "Copy"}
          </button>
        </div>
      ) : (
        !error && (
          <p className="mt-4 text-sm text-neutral-400 text-center">
            Your shortened URL will appear here.
          </p>
        )
      )}
    </div>
  );
};

export default Card;
