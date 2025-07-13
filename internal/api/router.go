package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GetDataRouter() *gin.Engine {
	r := gin.Default()
	r.Use(RateLimitMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/characters", GetCharactersHandler)
	r.GET("/healthcheck", GetHealthCheckHandler)
	return r
}
