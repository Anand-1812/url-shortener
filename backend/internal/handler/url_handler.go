// Package handler for api endpoints
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"backend/internal/repository"
	"backend/internal/service"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{service: service}
}

type ShortenURL struct {
	URL string `json:"url"`
}

type ShortenURLResponce struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenURL
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL != "" {
		http.Error(w, `{"error": "invalid body or missing url"}`, http.StatusBadRequest)
		return
	}

	record, err := h.service.Shorten(r.Context(), req.URL)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ShortenURLResponce{
		ShortCode:   record.ShortURL,
		ShortURL:    "http://" + r.Host + "/" + record.ShortURL,
		OriginalURL: record.OriginalURL,
	})
}

func (h URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.NotFound(w, r)
		return
	}

	record, err := h.service.Resolve(r.Context(), code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.NotFound(w, r)
			return
		}

		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, record.OriginalURL, http.StatusFound)
}
