package scraper

import "time"

type ScrapeRequest struct {
	URLs map[string]string `json:"urls"` 
}
type ScrapeResponse struct {
	Results map[string]*ScrapeResult `json:"results"` // key -> url and value -> scrape result for particular url 

}

type ScrapeResult struct {
	Status   string         `json:"status"`            
	Data     *ScrapedData   `json:"data,omitempty"`     
	Error    *ScrapeError   `json:"error,omitempty"`   
	Metadata *ScrapeMetadata `json:"metadata,omitempty"` 
}

type ScrapedData struct {
	// Metadata
	Title           string            `json:"title,omitempty"`
	MetaDescription string            `json:"meta_description,omitempty"`
	MetaKeywords    string            `json:"meta_keywords,omitempty"`
	CanonicalURL    string            `json:"canonical_url,omitempty"`
	OGTags          map[string]string `json:"og_tags,omitempty"` 
	
	// Content
	Headings   map[string][]string `json:"headings,omitempty"`  // h tags
	Paragraphs []string            `json:"paragraphs,omitempty"` // p tags 
	TextContent string             `json:"text_content,omitempty"` 
	
	// Structure
	Links  []Link  `json:"links,omitempty"`  // a tags
	Images []Image `json:"images,omitempty"` // img tags
	
	RawHTML string `json:"raw_html,omitempty"` // full html 
}

type Link struct {
	URL  string `json:"url"`
	Text string `json:"text"`
	Rel  string `json:"rel,omitempty"`  // relation between the current page and the linked page
}

type Image struct {
	Src string `json:"src"`
	Alt string `json:"alt"` // This attribute provides a descriptive text for the image. This text is displayed in place of the image if it fails to load
}

type ScrapeError struct {
	Type          ErrorType `json:"type"`                     
	Message       string    `json:"message"`                 
	HTTPStatus    int       `json:"http_status,omitempty"` 
	Timestamp     time.Time `json:"timestamp"`                
	Details       string    `json:"details,omitempty"`      
}

type ErrorType string
//constants for error types 
const (
    ErrorTypeTimeout    ErrorType = "timeout"
    ErrorTypeHTTP       ErrorType = "http_error"
    ErrorTypeNetwork    ErrorType = "network_error"
    ErrorTypeParse      ErrorType = "parse_error"
    ErrorTypeInvalidURL ErrorType = "invalid_url"
    ErrorTypeUnknown    ErrorType = "unknown"
)

type ScrapeMetadata struct {
	FinalURL     string        `json:"final_url"`     
	ResponseTime time.Duration `json:"response_time"`  
	Timestamp    time.Time     `json:"timestamp"`     
	ContentType  string        `json:"content_type,omitempty"`
	ContentLength int64        `json:"content_length,omitempty"`
}


type Config struct {
	// Concurrency controls how many URLs are scraped simultaneously
	MaxConcurrency int
	
	ConnectionTimeout time.Duration // time to establish connection
	ReadTimeout       time.Duration // time to read response
	TotalTimeout      time.Duration // overall timeout per request
	
	UserAgent      string
	FollowRedirects bool
	MaxRedirects   int
	
	IncludeRawHTML bool
	MaxTextLength  int 
}

func DefaultConfig() *Config {
	return &Config{
		MaxConcurrency:    10,
		ConnectionTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		TotalTimeout:      45 * time.Second,
		UserAgent:         "CrawlScraper/0.1.0",
		FollowRedirects:   true,
		MaxRedirects:      10,
		IncludeRawHTML:    false,
		MaxTextLength:     0, 
	}
}

// main functionality interface for the scraper


type Scraper interface {
	Scrape(req *ScrapeRequest) (*ScrapeResponse, error) // for multiple urls 
	
	ScrapeURL(url string) (*ScrapeResult, error) // for single url
}

// HTTPClient is an interface for making HTTP requests
type HTTPClient interface {
	Get(url string) (*HTTPResponse, error)
}

type HTTPResponse struct {
	StatusCode    int
	Body          []byte
	Headers       map[string][]string
	FinalURL      string 
	ContentType   string
	ContentLength int64
	ResponseTime  time.Duration
}

type Parser interface {
	Parse(html []byte) (*ScrapedData, error)
}
