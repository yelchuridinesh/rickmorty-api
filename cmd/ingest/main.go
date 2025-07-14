// package main

// import (
// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/cache"
// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/client"
// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/db"
// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/tracing"
// 	"context"
// 	"fmt"
// 	"log"
// )

// func main() {
// 	shutdown := tracing.InitTracer("rickmorty-api")
// 	defer shutdown(context.Background())
// 	cache.InitRedis()

// 	if err := db.InitPostgres(); err != nil {
// 		log.Fatalf("Postgres init failed: %v", err)
// 	}

// 	//check cache
// 	cached, err := cache.GetCachedCharacters()
// 	if err != nil {
// 		log.Fatalf("Redis error: %v", err)
// 	}

// 	if cached != nil {
// 		fmt.Println("Loaded characters from Redis cache")
// 		for _, c := range cached {
// 			fmt.Printf("- %s (ID: %d) from %s\n", c.Name, c.ID, c.Origin.Name)
// 		}
// 		return
// 	}

// 	// Cache miss → Fetch from API
// 	characters, err := client.GetCharsWithFilters()
// 	if err != nil {
// 		log.Fatalf("API error: %v", err)
// 	}

// 	if err := cache.SetCachedCharacters(characters); err != nil {
// 		log.Printf("Failed to cache results: %v", err)
// 	}
// 	err = db.InsertCharacters(characters)
// 	if err != nil {
// 		log.Printf("Failed to insert characters into DB: %v", err)
// 	}

// 	fmt.Printf("Fetched %d characters from API and cached them\n", len(characters))
// }

// ingest/main.go
package main

import (
	"context"
	"fmt"
	"log"

	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/cache"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/client"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/db"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/tracing"

	"go.opentelemetry.io/otel"
)

func main() {
	// 1) Initialize tracer
	shutdown := tracing.InitTracer("rickmorty-api")
	defer shutdown(context.Background())

	// 2) Create a root context & tracer handle
	ctx := context.Background()
	tr := otel.Tracer("rickmorty-api")

	// 3) Instrument Redis initialization
	ctx, span := tr.Start(ctx, "InitRedis")
	cache.InitRedis()
	span.End()

	// 4) Instrument Postgres initialization
	ctx, span = tr.Start(ctx, "InitPostgres")
	if err := db.InitPostgres(); err != nil {
		span.RecordError(err)
		span.End()
		log.Fatalf("Postgres init failed: %v", err)
	}
	span.End()

	// 5) Check cache
	ctx, span = tr.Start(ctx, "CheckCache")
	cached, err := cache.GetCachedCharacters()
	if err != nil {
		span.RecordError(err)
		span.End()
		log.Fatalf("Redis error: %v", err)
	}
	span.End()

	if cached != nil {
		ctx, span = tr.Start(ctx, "LoadFromCache")
		fmt.Println("Loaded characters from Redis cache")
		for _, c := range cached {
			fmt.Printf("- %s (ID: %d) from %s\n", c.Name, c.ID, c.Origin.Name)
		}
		span.End()
		return
	}

	// 6) Cache miss → Fetch from API
	ctx, span = tr.Start(ctx, "FetchFromAPI")
	characters, err := client.GetCharsWithFilters()
	if err != nil {
		span.RecordError(err)
		span.End()
		log.Fatalf("API error: %v", err)
	}
	span.End()

	// 7) Cache the results
	ctx, span = tr.Start(ctx, "CacheResults")
	if err := cache.SetCachedCharacters(characters); err != nil {
		span.RecordError(err)
		log.Printf("Failed to cache results: %v", err)
	}
	span.End()

	// 8) Insert into DB
	ctx, span = tr.Start(ctx, "InsertToDB")
	if err := db.InsertCharacters(characters); err != nil {
		span.RecordError(err)
		log.Printf("Failed to insert characters into DB: %v", err)
	}
	span.End()

	// 9) Final log
	ctx, span = tr.Start(ctx, "Done")
	fmt.Printf("Fetched %d characters from API and cached them\n", len(characters))
	span.End()
}
