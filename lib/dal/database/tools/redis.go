package tools

import (
	"flaver/globals"
	"sync"

	"github.com/go-redis/redis"
)

func RedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:       globals.GetViper().GetString("redis.address"),
		Password:   globals.GetViper().GetString("redis.password"),
		MaxRetries: globals.GetViper().GetInt("redis.maxretries"),
		DB:         globals.GetViper().GetInt("redis.db"),
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}

var (
	redisClient *redis.Client
	redisMux    sync.Mutex
)

func getRedisClient() *redis.Client {
	if redisClient != nil {
		return redisClient
	}
	redisMux.Lock()
	defer redisMux.Unlock()
	if redisClient != nil {
		return redisClient
	} else if client, err := RedisClient(); err != nil {
		globals.GetLogger().Fatalf(err.Error())
	} else {
		redisClient = client
	}
	return redisClient
}

func GetRedisClient() *redis.Client {
	return getRedisClient()
}