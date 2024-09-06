package cache

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *Cache
)

type Cache struct {
	Client *redis.Client
}

func New() *Cache {
	if redisClient != nil {
		return redisClient
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redisClient = &Cache{Client: rdb}
	return redisClient
}
