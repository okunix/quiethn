package database

import (
	"context"
	"strconv"
	"sync"

	"github.com/okunix/quiethn/config"
	"github.com/redis/go-redis/v9"
)

var (
	rdbInit sync.Once
	rdb     *redis.Client
)

func NewRedisClient() *redis.Client {
	rdbInit.Do(func() {
		redisDB, err := strconv.Atoi(config.RedisDB)
		if err != nil {
			panic(err)
		}
		rdb := redis.NewClient(&redis.Options{
			Addr:     config.RedisAddr,
			Password: config.RedisPassword,
			DB:       redisDB,
		})
		if err := rdb.Ping(context.TODO()).Err(); err != nil {
			panic(err)
		}
	})
	return rdb
}
