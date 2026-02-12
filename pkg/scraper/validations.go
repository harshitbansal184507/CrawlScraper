package scraper

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	MaxURLsPerRequest = 100
	MaxURLLength = 2048
)

func ValidateRequest(req *ScrapeRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	
	if  len(req.URLs) == 0 {
		return fmt.Errorf("URLs map cannot be empty")
	}
	
	if len(req.URLs) > MaxURLsPerRequest {
		return fmt.Errorf("too many URLs: got %d, max allowed is %d", len(req.URLs), MaxURLsPerRequest)
	}
	
	for key, urlStr := range req.URLs {
		if key == "" {
			return fmt.Errorf("URL key cannot be empty")
		}
		
		if err := ValidateURL(urlStr); err != nil {
			return fmt.Errorf("invalid URL for key '%s': %w", key, err)
		}
	}
	
	return nil
}

func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	
	if len(urlStr) > MaxURLLength {
		return fmt.Errorf("URL too long: %d characters (max: %d)", len(urlStr), MaxURLLength)
	}
	
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("malformed URL: %w", err)
	}
	
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported scheme '%s': only http and https are supported", parsedURL.Scheme)
	}
	
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	
	return nil
}

func NormalizeURL(urlStr string) string {
	urlStr = strings.TrimSpace(urlStr)
	
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}
	
	return urlStr
}

func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}
	
	if cfg.MaxConcurrency < 1 {
		return fmt.Errorf("MaxConcurrency must be at least 1, got %d", cfg.MaxConcurrency)
	}
	
	if cfg.MaxConcurrency > 100 {
		return fmt.Errorf("MaxConcurrency too high: %d (max: 100)", cfg.MaxConcurrency)
	}
	
	if cfg.ConnectionTimeout <= 0 {
		return fmt.Errorf("ConnectionTimeout must be positive, got %v", cfg.ConnectionTimeout)
	}
	
	if cfg.ReadTimeout <= 0 {
		return fmt.Errorf("ReadTimeout must be positive, got %v", cfg.ReadTimeout)
	}
	
	if cfg.TotalTimeout <= 0 {
		return fmt.Errorf("TotalTimeout must be positive, got %v", cfg.TotalTimeout)
	}
	
	if cfg.MaxRedirects < 0 {
		return fmt.Errorf("MaxRedirects cannot be negative, got %d", cfg.MaxRedirects)
	}
	
	if cfg.UserAgent == "" {
		return fmt.Errorf("UserAgent cannot be empty")
	}
	
	return nil
}