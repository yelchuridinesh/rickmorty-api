package api

import (
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/cache"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/db"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetHealthCheckHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	dbOk := make(chan bool, 1)

	//go routine for database check
	go func() {
		err := db.Database.PingContext(ctx)
		dbOk <- (err == nil)
	}()

	cacheOk := make(chan bool, 1)

	go func() {
		cacheOk <- cache.IsAlive(ctx)

	}()

	var dbStatus, redisStatus bool

	select {
	case dbStatus = <-dbOk:
	case <-ctx.Done():
		dbStatus = false
	}

	select {
	case redisStatus = <-cacheOk:
	case <-ctx.Done():
		redisStatus = false
	}

	if !dbStatus || !redisStatus {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"db":    dbStatus,
			"cache": redisStatus,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
