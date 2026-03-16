package main

import (
	"log"

	"github.com/harshitbansal184507/CrawlScraper/internal/config"
	"github.com/harshitbansal184507/CrawlScraper/internal/database"
)

func main() {
		cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Printf("Server ready on port %s", cfg.Server.Port)
	
	select {}
}