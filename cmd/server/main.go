package main

import (
	"context"
	"log"
	"time"

	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/api"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/cache"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/db"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/tracing"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// 1) Initialize tracer once for the whole app
	shutdown := tracing.InitTracer("rickmorty-api")
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdown(ctx)
	}()

	// 2) Instrument Postgres initialization
	if err := db.InitPostgres(); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}

	// 3) Instrument Redis initialization
	cache.InitRedis()

	// 4) Build your Gin router
	router := api.GetDataRouter()

	// 5) Add the OpenTelemetry Gin middleware so every request is traced
	router.Use(otelgin.Middleware("rickmorty-api"))

	log.Println("✅ Server running at http://localhost:8080")

	// 6) Start the server on :8080 as before
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// package main

// import (
// 	"context"
// 	"log"

// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/api"
// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/cache"
// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/db"
// 	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/tracing"

// 	"go.opentelemetry.io/otel"
// )

// func main() {
// 	// initialize tracing
// 	shutdown := tracing.InitTracer("rickmorty-api")
// 	defer shutdown(context.Background())

// 	// create a tracer and root context
// 	ctx := context.Background()
// 	tr := otel.Tracer("rickmorty-api")

// 	// instrument Postgres init
// 	ctx, span := tr.Start(ctx, "InitPostgres")
// 	if err := db.InitPostgres(); err != nil {
// 		span.RecordError(err)
// 		span.End()
// 		log.Fatal("DB init failed:", err)
// 	}
// 	span.End()

// 	// instrument Redis init
// 	ctx, span = tr.Start(ctx, "InitRedis")
// 	cache.InitRedis()
// 	span.End()

// 	// instrument router creation
// 	ctx, span = tr.Start(ctx, "SetupRouter")
// 	router := api.GetDataRouter()
// 	span.End()

// 	log.Println("✅ Server running at http://localhost:8080")

// 	// instrument server start (this span will end when the server shuts down)
// 	ctx, span = tr.Start(ctx, "StartServer")
// 	router.Run(":8080")
// 	span.End()
// }
