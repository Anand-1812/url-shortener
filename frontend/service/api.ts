import {
  ShortenRequest,
  ShortenResponse,
  ApiErrorResponse,
} from "@/types/types";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

export async function shortenURL(
  originalURL: string,
): Promise<ShortenResponse> {
  const payload: ShortenRequest = { url: originalURL };

  const response = await fetch(`${API_BASE_URL}/shorten`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    let errorMessage = "failed to short url";
    try {
      const errorData: ApiErrorResponse = await response.json();
      errorMessage = errorData.error || errorMessage;
    } catch {
      const text = await response.text();
      if (text) errorMessage = text;
    }

    throw new Error(errorMessage);
  }

  return response.json();
}
