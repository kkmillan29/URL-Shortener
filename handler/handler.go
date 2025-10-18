package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/komalkantmillan/urlshortener/model"
	"github.com/komalkantmillan/urlshortener/service"
)

// Handler handles HTTP requests
type Handler struct {
	shortenerService *service.URLShortenerService
	baseURL          string
}

// NewHandler creates a new handler
func NewHandler(shortenerService *service.URLShortenerService) *Handler {
	return &Handler{
		shortenerService: shortenerService,
		baseURL:          "http://localhost:8080/",
	}
}

// RegisterHandlers registers all API endpoints
func RegisterHandlers(router *mux.Router, shortenerService *service.URLShortenerService) {
	handler := NewHandler(shortenerService)

	router.HandleFunc("/shorten", handler.ShortenURL).Methods("POST")
	router.HandleFunc("/metrics", handler.GetMetrics).Methods("GET")
	router.HandleFunc("/{shortURL}", handler.RedirectURL).Methods("GET")
}

// ShortenURL handles URL shortening requests
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req model.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortCode, err := h.shortenerService.ShortenURL(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := h.baseURL + shortCode
	resp := model.ShortenResponse{
		OriginalURL: req.URL,
		ShortURL:    shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RedirectURL redirects short URLs to their original URLs
func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	originalURL, exists := h.shortenerService.GetOriginalURL(shortURL)
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Ensure URL has scheme
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

// GetMetrics returns the top domains metrics
func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	topDomains := h.shortenerService.GetTopDomains(3)
	
	resp := model.MetricsResponse{
		TopDomains: topDomains,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}