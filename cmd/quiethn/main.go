package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/okunix/quiethn"
	"github.com/okunix/quiethn/config"
	"github.com/okunix/quiethn/hackernews"
	"github.com/okunix/quiethn/router"
	"github.com/redis/go-redis/v9"
)

func main() {
	serverPort := config.ServerPort
	serverHost := config.ServerHost

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

	hnRepo := hackernews.NewHackerNewsRepo(rdb)

	router := router.NewRouter(hnRepo, quiethn.StaticFS)

	slog.Info("server is running", "host", serverHost, "port", serverPort)

	err = http.ListenAndServe(net.JoinHostPort(serverHost, serverPort), router)
	if err != nil {
		slog.Error("server fail", "error", err.Error())
	}
}
