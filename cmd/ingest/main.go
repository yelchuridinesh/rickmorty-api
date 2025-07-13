package main

import (
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/cache"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/client"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/db"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/tracing"
	"context"
	"fmt"
	"log"
)

func main() {
	shutdown := tracing.InitTracer("rickmorty-api")
	defer shutdown(context.Background())
	cache.InitRedis()

	if err := db.InitPostgres(); err != nil {
		log.Fatalf("Postgres init failed: %v", err)
	}

	//check cache
	cached, err := cache.GetCachedCharacters()
	if err != nil {
		log.Fatalf("Redis error: %v", err)
	}

	if cached != nil {
		fmt.Println("Loaded characters from Redis cache")
		for _, c := range cached {
			fmt.Printf("- %s (ID: %d) from %s\n", c.Name, c.ID, c.Origin.Name)
		}
		return
	}

	// Cache miss → Fetch from API
	characters, err := client.GetCharsWithFilters()
	if err != nil {
		log.Fatalf("API error: %v", err)
	}

	if err := cache.SetCachedCharacters(characters); err != nil {
		log.Printf("Failed to cache results: %v", err)
	}
	err = db.InsertCharacters(characters)
	if err != nil {
		log.Printf("Failed to insert characters into DB: %v", err)
	}

	fmt.Printf("Fetched %d characters from API and cached them\n", len(characters))
}
