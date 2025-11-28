package config

import (
	"fmt"
	"os"
)

func GetenvWithDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func MustGetenv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic(fmt.Sprintf("env \"%s\" is not defined", key))
}

var (
	ServerPort = GetenvWithDefault("HN_SERVER_PORT", "80")
	ServerHost = GetenvWithDefault("HN_SERVER_HOST", "0.0.0.0")

	RedisAddr     = MustGetenv("HN_REDIS_ADDR")
	RedisPassword = GetenvWithDefault("HN_REDIS_PASSWORD", "")
	RedisDB       = GetenvWithDefault("HN_REDIS_DB", "0")
)
