package service

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/komalkantmillan/urlshortener/model"
)

// URLShortenerService handles URL shortening operations
type URLShortenerService struct {
	urlMap       map[string]string // short -> original
	originalMap  map[string]string // original -> short
	domainCounts map[string]int    // domain -> count
	mutex        sync.RWMutex
}

// NewURLShortenerService creates a new URL shortener service
func NewURLShortenerService() *URLShortenerService {
	return &URLShortenerService{
		urlMap:       make(map[string]string),
		originalMap:  make(map[string]string),
		domainCounts: make(map[string]int),
	}
}

// ShortenURL shortens a URL and returns the shortened version
func (s *URLShortenerService) ShortenURL(originalURL string) (string, error) {
	// Validate URL
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	if parsedURL.Scheme == "" {
		originalURL = "http://" + originalURL
		parsedURL, _ = url.Parse(originalURL)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if URL has already been shortened
	if shortURL, exists := s.originalMap[originalURL]; exists {
		return shortURL, nil
	}

	// Generate short URL
	shortURL := s.generateShortURL(originalURL)

	// Store mappings
	s.urlMap[shortURL] = originalURL
	s.originalMap[originalURL] = shortURL

	// Update domain counts for metrics
	domain := parsedURL.Hostname()
	s.domainCounts[domain]++

	return shortURL, nil
}

// GetOriginalURL retrieves the original URL from a shortened URL
func (s *URLShortenerService) GetOriginalURL(shortURL string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	originalURL, exists := s.urlMap[shortURL]
	return originalURL, exists
}

// GetTopDomains returns the top N domains by count
func (s *URLShortenerService) GetTopDomains(n int) []model.DomainCount {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Convert map to slice for sorting
	var domainCounts []model.DomainCount
	for domain, count := range s.domainCounts {
		domainCounts = append(domainCounts, model.DomainCount{
			Domain: domain,
			Count:  count,
		})
	}

	// Sort by count (descending)
	sort.Slice(domainCounts, func(i, j int) bool {
		return domainCounts[i].Count > domainCounts[j].Count
	})

	// Return top N or all if less than N
	if len(domainCounts) > n {
		return domainCounts[:n]
	}
	return domainCounts
}

// generateShortURL creates a short URL from the original URL
func (s *URLShortenerService) generateShortURL(originalURL string) string {
	// Create MD5 hash of the URL
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Take first 8 characters for the short URL
	shortCode := strings.TrimRight(hash[:8], "=")
	return shortCode
}