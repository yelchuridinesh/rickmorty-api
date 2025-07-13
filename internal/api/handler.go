package api

import (
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCharactersHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	sortBy := c.DefaultQuery("sort", "id")

	page, err1 := strconv.Atoi(pageStr)
	limit, err2 := strconv.Atoi(limitStr)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page and limit must be integers"})
		return
	}

	if page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be >= 1"})
		return
	}

	if limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be between 1 and 100"})
		return
	}

	validSort := map[string]bool{"id": true, "name": true}
	if !validSort[sortBy] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort field. Allowed: id, name"})
		return
	}

	offset := (page - 1) * limit

	if chars, ok := getCharactersFromCache(limit, offset, sortBy); ok {
		c.JSON(http.StatusOK, chars)
		return
	}

	chars, err := db.GetCharacters(limit, offset, sortBy)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "failed to fetch characters from DB"})
		return
	}

	c.JSON(http.StatusOK, chars)

}
