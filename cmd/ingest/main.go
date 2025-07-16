package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/cache"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/client"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/db"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/tracing"

	"go.opentelemetry.io/otel"
)

func main() {
	// 1) Initialize tracer once for the whole program
	shutdown := tracing.InitTracer("rickmorty-api")
	defer func() {
		// give up if shutdown takes longer than 5s
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdown(ctx)
	}()

	// 2) Create a root context & tracer handle
	ctx := context.Background()
	tracer := otel.Tracer("rickmorty-api")

	// 3) Instrument Redis initialization
	ctx, span := tracer.Start(ctx, "InitRedis")
	cache.InitRedis()
	span.End()

	// 4) Instrument Postgres initialization
	ctx, span = tracer.Start(ctx, "InitPostgres")
	if err := db.InitPostgres(); err != nil {
		span.RecordError(err)
		span.End()
		log.Fatalf("Postgres init failed: %v", err)
	}
	span.End()

	// 5) Check cache
	ctx, span = tracer.Start(ctx, "CheckCache")
	cached, err := cache.GetCachedCharacters()
	if err != nil {
		span.RecordError(err)
		span.End()
		log.Fatalf("Redis error: %v", err)
	}
	span.End()

	if cached != nil {
		// 6a) Cache hit
		ctx, span = tracer.Start(ctx, "LoadFromCache")
		fmt.Println("Loaded characters from Redis cache")
		for _, c := range cached {
			fmt.Printf("- %s (ID: %d) from %s\n", c.Name, c.ID, c.Origin.Name)
		}
		span.End()
		return
	}

	// 6b) Cache miss → Fetch from external API
	ctx, span = tracer.Start(ctx, "FetchFromAPI")
	characters, err := client.GetCharsWithFilters()
	if err != nil {
		span.RecordError(err)
		span.End()
		log.Fatalf("API error: %v", err)
	}
	span.End()

	// 7) Cache the results in Redis
	ctx, span = tracer.Start(ctx, "CacheResults")
	if err := cache.SetCachedCharacters(characters); err != nil {
		span.RecordError(err)
		log.Printf("Failed to cache results: %v", err)
	}
	span.End()

	// 8) Insert characters into Postgres
	ctx, span = tracer.Start(ctx, "InsertToDB")
	if err := db.InsertCharacters(characters); err != nil {
		span.RecordError(err)
		log.Printf("Failed to insert characters into DB: %v", err)
	}
	span.End()

	// 9) Final log
	ctx, span = tracer.Start(ctx, "Done")
	fmt.Printf("Fetched %d characters from API and cached them\n", len(characters))
	span.End()
}
