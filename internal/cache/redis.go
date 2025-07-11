package cache

import (
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/model"
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Context     = context.Background()
	redisTTL    = 5 * time.Minute
	redisKey    = "rickmorty:characters:human_alive_earth"
	redisClient *redis.Client
)

func InitRedis() {
	ttlStr := os.Getenv("REDIS_TTL_MINUTES")
	if ttlStr != "" {
		if ttl, err := strconv.Atoi(ttlStr); err == nil {
			redisTTL = time.Duration(ttl) * time.Minute
		}
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
}

func GetCachedCharacters() ([]model.Character, error) {
	data, err := redisClient.Get(Context, redisKey).Result()
	if err == redis.Nil {
		return nil, nil //cache miss
	} else if err != nil {
		return nil, err
	}
	var characters []model.Character
	if err := json.Unmarshal([]byte(data), &characters); err != nil {
		return nil, err
	}
	return characters, nil

}

func SetCachedCharacters(characters []model.Character) error {
	jsonData, err := json.Marshal(characters)
	if err != nil {
		return err
	}
	return redisClient.Set(Context, redisKey, jsonData, redisTTL).Err()
}

func IsAlive(ctx context.Context) bool {
	_, err := redisClient.Ping(Context).Result()
	return err == nil
}
