package scraper

import (
	"fmt"
	"time"
)

type ScraperError struct {
	Type       ErrorType
	Message    string
	Err        error
	URL        string
	StatusCode int
	Timestamp  time.Time
	
}

func (e *ScraperError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (underlying: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}


func (e *ScraperError) Unwrap() error {
	return e.Err
}


func NewScraperError(errType ErrorType, message string, err error) *ScraperError {
	return &ScraperError{
		Type:      errType,
		Message:   message,
		Err:       err,
		Timestamp: time.Now(),
	}
}


func NewTimeoutError(url string, err error) *ScraperError {
	return &ScraperError{
		Type:      ErrorTypeTimeout,
		Message:   fmt.Sprintf("request timeout for %s", url),
		Err:       err,
		URL:       url,
		Timestamp: time.Now(),
	}
}

func NewHTTPError(url string, statusCode int, message string) *ScraperError {
	return &ScraperError{
		Type:       ErrorTypeHTTP,
		Message:    fmt.Sprintf("HTTP %d: %s", statusCode, message),
		URL:        url,
		StatusCode: statusCode,
		Timestamp:  time.Now(),
	}
}

func NewNetworkError(url string, err error) *ScraperError {
	return &ScraperError{
		Type:      ErrorTypeNetwork,
		Message:   fmt.Sprintf("network error for %s", url),
		Err:       err,
		URL:       url,
		Timestamp: time.Now(),
	}
}

func NewParseError(url string, err error) *ScraperError {
	return &ScraperError{
		Type:      ErrorTypeParse,
		Message:   fmt.Sprintf("failed to parse HTML from %s", url),
		Err:       err,
		URL:       url,
		Timestamp: time.Now(),
	}
}

func NewInvalidURLError(url string, err error) *ScraperError {
	return &ScraperError{
		Type:      ErrorTypeInvalidURL,
		Message:   fmt.Sprintf("invalid URL: %s", url),
		Err:       err,
		URL:       url,
		Timestamp: time.Now(),
	}
}

func (e *ScraperError) ToScrapeError() *ScrapeError {
	retryPossible := e.Type == ErrorTypeTimeout || 
		e.Type == ErrorTypeNetwork || 
		(e.Type == ErrorTypeHTTP && e.StatusCode >= 500)
	
	details := ""
	if e.Err != nil {
		details = e.Err.Error()
	}
	
	return &ScrapeError{
		Type:          e.Type,
		Message:       e.Message,
		HTTPStatus:    e.StatusCode,
		Timestamp:     e.Timestamp,
		RetryPossible: retryPossible,
		Details:       details,
	}
}
