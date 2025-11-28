package main

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/okunix/quiethn"
	"github.com/okunix/quiethn/config"
	"github.com/okunix/quiethn/database"
	"github.com/okunix/quiethn/hackernews"
	"github.com/okunix/quiethn/router"
)

func main() {
	serverPort := config.ServerPort
	serverHost := config.ServerHost

	rdb := database.NewRedisClient()
	hnRepo := hackernews.NewHackerNewsRepo(rdb)

	router := router.NewRouter(hnRepo, quiethn.StaticFS)

	slog.Info("server is running", "host", serverHost, "port", serverPort)

	err := http.ListenAndServe(net.JoinHostPort(serverHost, serverPort), router)
	if err != nil {
		slog.Error("server fail", "error", err.Error())
	}
}
