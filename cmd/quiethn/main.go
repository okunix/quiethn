package main

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/okunix/quiethn"
	"github.com/okunix/quiethn/config"
	"github.com/okunix/quiethn/hackernews"
	"github.com/okunix/quiethn/router"
)

func main() {
	serverPort := config.ServerPort
	serverHost := config.ServerHost

	hn := hackernews.NewHackerNewsClient("https://hacker-news.firebaseio.com")
	hnCache := hackernews.NewHackerNewsClientWithCache(hn)

	router := router.NewRouter(hnCache, quiethn.StaticFS)

	slog.Info("server is running", "host", serverHost, "port", serverPort)

	err := http.ListenAndServe(net.JoinHostPort(serverHost, serverPort), router)
	if err != nil {
		slog.Error("server fail", "error", err.Error())
	}
}
