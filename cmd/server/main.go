package main

import (
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/api"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/cache"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/db"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/tracing"
	"context"
	"log"
)

func main() {
	shutdown := tracing.InitTracer("rickmorty-api")
	defer shutdown(context.Background())
	if err := db.InitPostgres(); err != nil {
		log.Fatal("DB init failed:", err)
	}
	cache.InitRedis()

	router := api.GetDataRouter()
	log.Println("✅ Server running at http://localhost:8080")
	router.Run(":8080")
}
