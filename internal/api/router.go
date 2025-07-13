package api

import (
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GetDataRouter() *gin.Engine {
	r := gin.Default()
	r.Use(RateLimitMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/healthcheck", metrics.InstrumentHandler("HealthCheck", GetHealthCheckHandler))
	r.GET("/characters", metrics.InstrumentHandler("GetCharacters", GetCharactersHandler))
	return r
}
