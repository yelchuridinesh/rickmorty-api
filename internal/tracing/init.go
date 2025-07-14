// internal/tracing/init.go
package tracing

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func InitTracer(serviceName string) func(context.Context) {
	// Determine the collector endpoint:

	endpoint := "http://jaeger-collector.observability.svc.cluster.local:14268/api/traces"
	log.Printf("[OTEL] Using Jaeger endpoint: %s", endpoint)

	// Create the exporter:
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		log.Fatalf("[OTEL][ERROR] failed to initialize Jaeger exporter: %v", err)
	}

	// Build the TraceProvider:
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	log.Printf("[OTEL] TracerProvider initialized for service %q", serviceName)

	// Return a shutdown function that will flush and log any error:
	return func(ctx context.Context) {
		// give up if it takes longer than 5s to shut down
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("[OTEL][ERROR] error shutting down tracer provider: %v", err)
		} else {
			log.Print("[OTEL] TracerProvider shut down cleanly")
		}
	}
}
