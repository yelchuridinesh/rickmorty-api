// package main

// import (
// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/api"
// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/cache"
// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/db"
// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/tracing"
// 	"context"
// 	"log"
// )

// func main() {
// 	shutdown := tracing.InitTracer("rickmorty-api")
// 	defer shutdown(context.Background())
// 	if err := db.InitPostgres(); err != nil {
// 		log.Fatal("DB init failed:", err)
// 	}
// 	cache.InitRedis()

// 	router := api.GetDataRouter()
// 	log.Println("✅ Server running at http://localhost:8080")
// 	router.Run(":8080")
// }

package main

import (
	"context"
	"log"

	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/api"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/cache"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/db"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/tracing"

	"go.opentelemetry.io/otel"
)

func main() {
	// initialize tracing
	shutdown := tracing.InitTracer("rickmorty-api")
	defer shutdown(context.Background())

	// create a tracer and root context
	ctx := context.Background()
	tr := otel.Tracer("rickmorty-api")

	// instrument Postgres init
	ctx, span := tr.Start(ctx, "InitPostgres")
	if err := db.InitPostgres(); err != nil {
		span.RecordError(err)
		span.End()
		log.Fatal("DB init failed:", err)
	}
	span.End()

	// instrument Redis init
	ctx, span = tr.Start(ctx, "InitRedis")
	cache.InitRedis()
	span.End()

	// instrument router creation
	ctx, span = tr.Start(ctx, "SetupRouter")
	router := api.GetDataRouter()
	span.End()

	log.Println("✅ Server running at http://localhost:8080")

	// instrument server start (this span will end when the server shuts down)
	ctx, span = tr.Start(ctx, "StartServer")
	router.Run(":8080")
	span.End()
}
