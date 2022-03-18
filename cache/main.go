package cache

import (
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
		Addr:     "0.0.0.0:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redisClient = &Cache{Client: rdb}
	return redisClient
}
