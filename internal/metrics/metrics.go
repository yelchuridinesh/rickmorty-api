// internal/metrics/metrics.go
package metrics

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	// // "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
)

// Histogram of HTTP handler durations, labeled by handler name and HTTP status.
var RequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "api_request_duration_seconds",
		Help: "Histogram of HTTP request latencies.",
		// You can also set buckets here if you want non-default ranges.
	},
	[]string{"handler", "status"},
)

// Counters for cache hits/misses
var (
	CacheHits   = promauto.NewCounter(prometheus.CounterOpts{Name: "api_cache_hits_total", Help: "Total number of cache hits."})
	CacheMisses = promauto.NewCounter(prometheus.CounterOpts{Name: "api_cache_misses_total", Help: "Total number of cache misses."})
)

// Business metric: how many Character objects have been returned
var CharactersProcessed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "api_characters_processed_total",
	Help: "Total number of characters returned to clients.",
})

// InstrumentHandler wraps any gin.HandlerFunc, measures its duration and status code.
func InstrumentHandler(name string, h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			status := fmt.Sprint(c.Writer.Status())
			RequestDuration.WithLabelValues(name, status).Observe(v)
		}))
		defer timer.ObserveDuration()

		h(c)
	}
}
