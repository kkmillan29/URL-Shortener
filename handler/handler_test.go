package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/komalkantmillan/urlshortener/model"
	"github.com/komalkantmillan/urlshortener/service"
)

func TestShortenURLHandler(t *testing.T) {
	shortenerService := service.NewURLShortenerService()
	handler := NewHandler(shortenerService)

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "Valid URL",
			requestBody:    map[string]string{"url": "https://www.example.com"},
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "Empty URL",
			requestBody:    map[string]string{"url": ""},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  false,
		},
		{
			name:           "Invalid JSON",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			if tt.requestBody != nil {
				reqBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			} else {
				reqBody = []byte("{invalid json")
			}

			req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.ShortenURL(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.checkResponse {
				var response model.ShortenResponse
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response.OriginalURL != tt.requestBody["url"] {
					t.Errorf("Handler returned wrong original URL: got %v want %v", response.OriginalURL, tt.requestBody["url"])
				}

				if response.ShortURL == "" {
					t.Errorf("Handler returned empty short URL")
				}
			}
		})
	}
}

func TestRedirectURLHandler(t *testing.T) {
	shortenerService := service.NewURLShortenerService()
	handler := NewHandler(shortenerService)

	// First, create a shortened URL
	originalURL := "https://www.example.com"
	shortCode, err := shortenerService.ShortenURL(originalURL)
	if err != nil {
		t.Fatalf("Failed to shorten URL: %v", err)
	}

	tests := []struct {
		name           string
		shortURL       string
		expectedStatus int
		expectedURL    string
	}{
		{
			name:           "Valid Short URL",
			shortURL:       shortCode,
			expectedStatus: http.StatusMovedPermanently,
			expectedURL:    originalURL,
		},
		{
			name:           "Invalid Short URL",
			shortURL:       "nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedURL:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/"+tt.shortURL, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/{shortURL}", handler.RedirectURL).Methods("GET")
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedURL != "" {
				location := rr.Header().Get("Location")
				if location != tt.expectedURL {
					t.Errorf("Handler redirected to wrong URL: got %v want %v", location, tt.expectedURL)
				}
			}
		})
	}
}

func TestGetMetricsHandler(t *testing.T) {
	shortenerService := service.NewURLShortenerService()
	handler := NewHandler(shortenerService)

	// Add some URLs with different domains
	urls := []string{
		"https://www.example.com/page1",
		"https://www.example.com/page2",
		"https://www.test.com/page1",
	}

	for _, url := range urls {
		_, err := shortenerService.ShortenURL(url)
		if err != nil {
			t.Fatalf("Failed to shorten URL: %v", err)
		}
	}

	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.GetMetrics(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response model.MetricsResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.TopDomains) != 2 {
		t.Errorf("Handler returned wrong number of domains: got %v want %v", len(response.TopDomains), 2)
	}

	if response.TopDomains[0].Domain != "www.example.com" || response.TopDomains[0].Count != 2 {
		t.Errorf("Handler returned wrong top domain: got %v with count %v, want www.example.com with count 2",
			response.TopDomains[0].Domain, response.TopDomains[0].Count)
	}
}