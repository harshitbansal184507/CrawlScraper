package main

import (
	"fmt"

	"github.com/harshitbansal184507/CrawlScraper/pkg/scraper"
)

func main() {
	s := scraper.New(scraper.DefaultConfig())
	
	result, _ := s.ScrapeURL("https://example.com/")
	
	if result.Status == "success" {
		fmt.Printf("Title: %s\n", result.Data.Title)
		fmt.Printf("Paragraphs: %d\n", len(result.Data.Paragraphs))
	}
}