package scraper

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// httpClient implements the HTTPClient interface
// It wraps Go's standard http.Client with our configuration
type httpClient struct {
	client *http.Client // The underlying Go HTTP client
	config *Config      // Our scraper configuration
}

// NewHTTPClient creates a new HTTP client with the given configuration
// This is a factory function that sets up the client properly
func NewHTTPClient(cfg *Config) HTTPClient {
	// Validate config first
	if err := ValidateConfig(cfg); err != nil {
		// If config is invalid, return client with default config
		// In production, you might want to panic or return error
		cfg = DefaultConfig()
	}
	
	// Create custom transport with our timeout settings
	transport := &http.Transport{
		// How long to wait to establish a connection
		DialContext: (&net.Dialer{
			Timeout: cfg.ConnectionTimeout,
		}).DialContext,
		
		// Maximum time to wait for response headers
		ResponseHeaderTimeout: cfg.ReadTimeout,
		
		// Compression settings - ensure gzip is handled automatically
		DisableCompression: false, // Allow automatic decompression
		
		// Connection pool settings
		MaxIdleConns:        100,              // Max idle connections total
		MaxIdleConnsPerHost: 10,               // Max idle per host
		IdleConnTimeout:     90 * time.Second, // How long to keep idle connections
	}
	
	// Create the HTTP client
	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.TotalTimeout, // Overall request timeout
		
		// Custom redirect policy
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// via contains all previous requests in redirect chain
			if len(via) >= cfg.MaxRedirects {
				return fmt.Errorf("stopped after %d redirects", cfg.MaxRedirects)
			}
			
			// If FollowRedirects is false, stop after first redirect
			if !cfg.FollowRedirects && len(via) > 0 {
				return http.ErrUseLastResponse // Stop following redirects
			}
			
			return nil // Continue following redirects
		},
	}
	
	return &httpClient{
		client: client,
		config: cfg,
	}
}

// Get fetches a URL and returns our custom HTTPResponse
// This implements the HTTPClient interface
func (c *httpClient) Get(url string) (*HTTPResponse, error) {
	// Record start time for measuring response time
	startTime := time.Now()
	
	// Validate URL before making request
	if err := ValidateURL(url); err != nil {
		return nil, NewInvalidURLError(url, err)
	}
	
	// Create HTTP request with context for timeout control
	ctx, cancel := context.WithTimeout(context.Background(), c.config.TotalTimeout)
	defer cancel() // Always cancel context when done
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, NewNetworkError(url, err)
	}
	
	// Set User-Agent header
	req.Header.Set("User-Agent", c.config.UserAgent)
	
	// Add common headers that browsers send
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	// Note: Do NOT set Accept-Encoding manually - Go's Transport handles this automatically
	// when DisableCompression is false
	
	// Make the actual HTTP request
	resp, err := c.client.Do(req)
	if err != nil {
		// Check if it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return nil, NewTimeoutError(url, err)
		}
		// Otherwise it's a network error
		return nil, NewNetworkError(url, err)
	}
	defer resp.Body.Close() // Always close response body
	
	// Calculate response time
	responseTime := time.Since(startTime)
	
	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		// Read a bit of the body for error context
		bodySnippet := make([]byte, 200)
		n, _ := resp.Body.Read(bodySnippet)
		
		statusText := http.StatusText(resp.StatusCode)
		if n > 0 {
			statusText = fmt.Sprintf("%s: %s", statusText, string(bodySnippet[:n]))
		}
		
		return nil, NewHTTPError(url, resp.StatusCode, statusText)
	}
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewNetworkError(url, fmt.Errorf("failed to read response body: %w", err))
	}
	
	// Determine final URL (after redirects)
	finalURL := url
	if resp.Request != nil && resp.Request.URL != nil {
		finalURL = resp.Request.URL.String()
	}
	
	// Get content type
	contentType := resp.Header.Get("Content-Type")
	
	// Get content length
	contentLength := resp.ContentLength
	if contentLength < 0 {
		// If not set in headers, use actual body length
		contentLength = int64(len(body))
	}
	
	// Build our custom HTTPResponse
	return &HTTPResponse{
		StatusCode:    resp.StatusCode,
		Body:          body,
		Headers:       resp.Header,
		FinalURL:      finalURL,
		ContentType:   contentType,
		ContentLength: contentLength,
		ResponseTime:  responseTime,
	}, nil
}