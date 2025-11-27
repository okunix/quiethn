package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/okunix/quiethn"
	"github.com/okunix/quiethn/hackernews"
	"github.com/okunix/quiethn/router"
)

func GetenvWithDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func main() {
	serverPort := GetenvWithDefault("HN_SERVER_PORT", "80")
	serverHost := GetenvWithDefault("HN_SERVER_HOST", "0.0.0.0")

	hn := hackernews.NewHackerNewsClient("https://hacker-news.firebaseio.com")
	hnCache := hackernews.NewHackerNewsClientWithCache(hn)

	router := router.NewRouter(hnCache, quiethn.StaticFS)

	slog.Info("server is running", "host", serverHost, "port", serverPort)

	err := http.ListenAndServe(net.JoinHostPort(serverHost, serverPort), router)
	if err != nil {
		slog.Error("server fail", "error", err.Error())
	}
}
