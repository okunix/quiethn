package main

import (
	"embed"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/okunix/quietHN/hackernews"
	"github.com/okunix/quietHN/router"
)

//go:embed static/*
var staticFS embed.FS

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

	router := router.NewRouter(hnCache, staticFS)

	slog.Info("server is running", "host", serverHost, "port", serverPort)
	http.ListenAndServe(net.JoinHostPort(serverHost, serverPort), router)
}
