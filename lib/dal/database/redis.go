package database

import (
	"flaver/lib/dal/database/tools"

	"github.com/go-redis/redis"
)

func GetRedisClient() *redis.Client {
	return tools.GetRedisClient()
}