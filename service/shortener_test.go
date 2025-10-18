package service

import (
	"testing"
)

func TestShortenURL(t *testing.T) {
	service := NewURLShortenerService()

	tests := []struct {
		name        string
		url         string
		wantErr     bool
		expectSame  bool
		secondInput string
	}{
		{
			name:    "Valid URL",
			url:     "https://www.example.com",
			wantErr: false,
		},
		{
			name:    "URL without scheme",
			url:     "www.example.com",
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			url:     "://invalid",
			wantErr: true,
		},
		{
			name:        "Same URL returns same short URL",
			url:         "https://www.test.com",
			wantErr:     false,
			expectSame:  true,
			secondInput: "https://www.test.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortURL, err := service.ShortenURL(tt.url)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ShortenURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err == nil {
				if shortURL == "" {
					t.Error("ShortenURL() returned empty short URL")
				}
				
				// Test that the original URL can be retrieved
				originalURL, exists := service.GetOriginalURL(shortURL)
				if !exists {
					t.Error("GetOriginalURL() could not find the short URL")
				}
				
				// Check if the original URL matches (accounting for http:// prefix addition)
				expectedURL := tt.url
				if originalURL != expectedURL && originalURL != "http://"+expectedURL {
					t.Errorf("GetOriginalURL() got = %v, want %v or %v", originalURL, expectedURL, "http://"+expectedURL)
				}
				
				// Test that the same URL returns the same short URL
				if tt.expectSame {
					secondShortURL, _ := service.ShortenURL(tt.secondInput)
					if shortURL != secondShortURL {
						t.Errorf("ShortenURL() did not return the same short URL for the same input. Got %v and %v", shortURL, secondShortURL)
					}
				}
			}
		})
	}
}

func TestGetTopDomains(t *testing.T) {
	service := NewURLShortenerService()
	
	// Add some URLs with different domains
	urls := []string{
		"https://www.example.com/page1",
		"https://www.example.com/page2",
		"https://www.example.com/page3",
		"https://www.test.com/page1",
		"https://www.test.com/page2",
		"https://www.another.com/page1",
	}
	
	for _, url := range urls {
		_, err := service.ShortenURL(url)
		if err != nil {
			t.Errorf("Failed to shorten URL: %v", err)
		}
	}
	
	// Get top domains
	topDomains := service.GetTopDomains(3)
	
	// Check if we got the expected number of domains
	if len(topDomains) != 3 {
		t.Errorf("GetTopDomains() returned %d domains, want %d", len(topDomains), 3)
	}
	
	// Check if domains are in the correct order (by count)
	if len(topDomains) >= 2 {
		if topDomains[0].Domain != "www.example.com" || topDomains[0].Count != 3 {
			t.Errorf("Expected top domain to be www.example.com with count 3, got %s with count %d", 
				topDomains[0].Domain, topDomains[0].Count)
		}
		
		if topDomains[1].Domain != "www.test.com" || topDomains[1].Count != 2 {
			t.Errorf("Expected second domain to be www.test.com with count 2, got %s with count %d", 
				topDomains[1].Domain, topDomains[1].Count)
		}
		
		if len(topDomains) >= 3 {
			if topDomains[2].Domain != "www.another.com" || topDomains[2].Count != 1 {
				t.Errorf("Expected third domain to be www.another.com with count 1, got %s with count %d", 
					topDomains[2].Domain, topDomains[2].Count)
			}
		}
	}
}