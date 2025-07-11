package api

import "github.com/gin-gonic/gin"

func GetDataRouter() *gin.Engine {
	r := gin.Default()
	r.Use(RateLimitMiddleware())
	r.GET("/characters", GetCharactersHandler)
	r.GET("/healthcheck", GetHealthCheckHandler)
	return r
}
