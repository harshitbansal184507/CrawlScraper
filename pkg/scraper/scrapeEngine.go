package scraper

import (
	"sync"
	"time"
)

// job represents a single URL scraping job
type job struct {
	key string // "/blogs" or " homepage"
	url string 
}

// result represents the result of scraping a single URL
type result struct {
	key    string       
	result *ScrapeResult 
}

// manages concurrent scraping of multiple URLs
type scrapeEngine struct {
	httpClient HTTPClient        // fetching pages
	parser     Parser            //extracting data
	config     *Config         
	jobs       chan job          // Channel for distributing jobs to workers
	results    chan result       // Channel for collecting results
	wg         sync.WaitGroup    // WaitGroup to track worker completion
}

//  creates a new worker pool
func newscrapeEngine(client HTTPClient, parser Parser, cfg *Config) *scrapeEngine {
	return &scrapeEngine{
		httpClient: client,
		parser:     parser,
		config:     cfg,
		jobs:       make(chan job, cfg.MaxConcurrency),  
		results:    make(chan result, cfg.MaxConcurrency),
	}
}

// process processes a batch of URLs concurrently
func (wp *scrapeEngine) process(urls map[string]string) map[string]*ScrapeResult {
	// Start workers
	wp.startWorkers()
	
	// Send jobs to workers
	go wp.sendJobs(urls)
	
	// Collect results
	return wp.collectResults(len(urls))
}

// start configured number of worker goroutines
func (wp *scrapeEngine) startWorkers() {
	// Create MaxConcurrency workers
	for i := 0; i < wp.config.MaxConcurrency; i++ {
		wp.wg.Add(1) // Increment WaitGroup counter
		
		// Start worker goroutine
		go wp.worker(i)
	}
}

// worker is a goroutine that processes jobs from the jobs channel
func (wp *scrapeEngine) worker(id int) {
	defer wp.wg.Done()
	
	for job := range wp.jobs {
		result := wp.scrapeURL(job.key, job.url)
		
		wp.results <- result
	}
}

// sendJobs sends all jobs to the jobs channel
func (wp *scrapeEngine) sendJobs(urls map[string]string) {
	// Send each URL as a job
	for key, url := range urls {
		wp.jobs <- job{
			key: key,
			url: url,
		}
	}
	
	// Close jobs channel to signal no more jobs
	close(wp.jobs)
	
	// Wait for all workers to finish
	wp.wg.Wait()
	
	// Close results channel to signal no more results
	close(wp.results)
}

// collectResults collects all results from the results channel
func (wp *scrapeEngine) collectResults(expected int) map[string]*ScrapeResult {
	results := make(map[string]*ScrapeResult, expected)
	
	// Collect all results
	for result := range wp.results {
		results[result.key] = result.result
	}
	
	return results
}

// scrapeURL scrapes a single URL and returns the result
func (wp *scrapeEngine) scrapeURL(key, url string) result {
	// Record start time
	startTime := time.Now()
	
	//  Fetch the page
	resp, err := wp.httpClient.Get(url)
	if err != nil {
		// HTTP request failed
		scraperErr, ok := err.(*ScraperError)
		if !ok {
			// Wrap error if it's not already a ScraperError
			scraperErr = NewNetworkError(url, err)
		}
		
		return result{
			key: key,
			result: &ScrapeResult{
				Status: "failed",
				Error:  scraperErr.ToScrapeError(),
				Metadata: &ScrapeMetadata{
					FinalURL:     url,
					ResponseTime: time.Since(startTime),
					Timestamp:    time.Now(),
				},
			},
		}
	}
	
	// Parsing
	data, err := wp.parser.Parse(resp.Body)
	if err != nil {
		// Parsing failed
		parseErr := NewParseError(resp.FinalURL, err)
		
		return result{
			key: key,
			result: &ScrapeResult{
				Status: "failed",
				Error:  parseErr.ToScrapeError(),
				Metadata: &ScrapeMetadata{
					FinalURL:      resp.FinalURL,
					ResponseTime:  resp.ResponseTime,
					Timestamp:     time.Now(),
					ContentType:   resp.ContentType,
					ContentLength: resp.ContentLength,
				},
			},
		}
	}
	//successful
	return result{
		key: key,
		result: &ScrapeResult{
			Status: "success",
			Data:   data,
			Metadata: &ScrapeMetadata{
				FinalURL:      resp.FinalURL,
				ResponseTime:  resp.ResponseTime,
				Timestamp:     time.Now(),
				ContentType:   resp.ContentType,
				ContentLength: resp.ContentLength,
			},
		},
	}
}

