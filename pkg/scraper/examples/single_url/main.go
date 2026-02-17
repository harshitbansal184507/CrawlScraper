package main

import (
	"fmt"
	"time"

	"github.com/harshitbansal184507/CrawlScraper/pkg/scraper"
)

func main() {
	s := scraper.New(scraper.DefaultConfig())

	start := time.Now()
	result, _ := s.ScrapeURL("https://example.com/")
	time_taken := time.Since(start)

	fmt.Println("Time taken for scraping :", time_taken)
	
	if result.Status == "success" {
		fmt.Printf("Title: %s\n", result.Data.Title)
		fmt.Printf("Paragraphs: %d\n", len(result.Data.Paragraphs))
	}
}