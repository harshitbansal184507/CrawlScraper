package main

import (
	"fmt"
	"time"

	"github.com/harshitbansal184507/CrawlScraper/pkg/scraper"
)


func main() {
	s := scraper.New(scraper.DefaultConfig())
	
	req := &scraper.ScrapeRequest{
		URLs: map[string]string{
			"example":   "https://example.com",
			"httpbin":   "https://httpbin.org/html",
			"hackernews": "https://news.ycombinator.com",
		},
	}


	start := time.Now()
	resp , err := s.Scrape(req)

	if err!=nil {
		fmt.Println("err:", err)
	}
	time_taken := time.Since(start).Milliseconds()

	println("Time Taken for Scraping(in milliseconds) :",time_taken)

	successCount := 0 


	for key, result := range resp.Results {
		if result.Status == "success" {
			successCount++
			fmt.Printf("✅ %s:\n", key)
			fmt.Printf("   Title: %s\n", result.Data.Title)
			fmt.Printf("   Content: %d chars, %d paragraphs, %d links\n",
				len(result.Data.TextContent),
				len(result.Data.Paragraphs),
				len(result.Data.Links))
			fmt.Printf("   Time: %v\n\n", result.Metadata.ResponseTime)
		} else {
			fmt.Printf("❌ %s: %s\n\n", key, result.Error.Message)
		}
	}

}