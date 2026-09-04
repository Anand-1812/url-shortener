export interface ShortenRequest {
  url: string;
}

export interface ShortenResponse {
  shortCode: string;
  shortURL: string;
  originalURL: string;
}

export interface ApiErrorResponse {
  error: string;
}

export interface HealthResponse {
  status: "healthy" | "unhealthy";
  database: "connected" | "disconnected";
}
