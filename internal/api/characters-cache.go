package api

import (
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/cache"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/model"
	"log"
	"sort"
)

func getCharactersFromCache(limit, offset int, sortBy string) ([]model.Character, bool) {
	cached, err := cache.GetCachedCharacters()
	if err != nil || cached == nil {
		log.Println("Cache MISS: falling back to DB")
		return nil, false
	}
	log.Println("Cache HIT: serving characters from Redis")
	// Sort
	if sortBy == "name" {
		sort.Slice(cached, func(i, j int) bool {
			return cached[i].Name < cached[j].Name
		})
	} else { // default is "id"
		sort.Slice(cached, func(i, j int) bool {
			return cached[i].ID < cached[j].ID
		})
	}

	// Apply pagination
	if offset >= len(cached) {
		return []model.Character{}, true
	}

	end := offset + limit
	if end > len(cached) {
		end = len(cached)
	}

	return cached[offset:end], true
}
