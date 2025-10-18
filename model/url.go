package model

// URLMapping represents a mapping between original and shortened URLs
type URLMapping struct {
	OriginalURL string
	ShortURL    string
}

// ShortenRequest represents a request to shorten a URL
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse represents a response with the shortened URL
type ShortenResponse struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// MetricsResponse represents the response for the metrics API
type MetricsResponse struct {
	TopDomains []DomainCount `json:"top_domains"`
}

// DomainCount represents a domain and its count
type DomainCount struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}