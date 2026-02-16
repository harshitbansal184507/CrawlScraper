package scraper

type scraper struct {
	httpClient HTTPClient
	parser     Parser
	config     *Config
}

//entry point 
func New(cfg *Config) Scraper {

	//validation
	if err := ValidateConfig(cfg); err != nil {
		cfg = DefaultConfig()
	}
	
	// Create HTTP client
	httpClient := NewHTTPClient(cfg)
	
	// Create parser
	parser := NewParser(cfg)
	
	return &scraper{
		httpClient: httpClient,
		parser:     parser,
		config:     cfg,
	}
}

// Scrape scrapes multiple URLs concurrently
func (s *scraper) Scrape(req *ScrapeRequest) (*ScrapeResponse, error) {
	// Validate request
	if err := ValidateRequest(req); err != nil {
		return nil, err
	}
	
	pool := newscrapeEngine(s.httpClient, s.parser, s.config)
	
	results := pool.process(req.URLs)
	
	// response
	return &ScrapeResponse{
		Results: results,
	}, nil
}

// for scraping a single URL
func (s *scraper) ScrapeURL(url string) (*ScrapeResult, error) {
	if err := ValidateURL(url); err != nil {
		return nil, NewInvalidURLError(url, err)
	}
	
	req := &ScrapeRequest{
		URLs: map[string]string{
			"default": url,
		},
	}
	
	resp, err := s.Scrape(req)
	if err != nil {
		return nil, err
	}
	
	return resp.Results["default"], nil
}