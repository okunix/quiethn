package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/okunix/quiethn/config"
	"github.com/okunix/quiethn/database"
	"github.com/okunix/quiethn/hackernews"
)

func main() {
	baseURL := config.HackerNewsBaseURL
	ctx := context.Background()
	rdb := database.NewRedisClient()
	hnRepo := hackernews.NewHackerNewsRepo(rdb)
	hnClient := hackernews.NewHackerNewsClient(baseURL)
	topStories, err := hnClient.GetTopStories(ctx, 30)
	if err != nil {
		slog.Error("failed to fetch top stories", "error", err.Error())
		os.Exit(1)
	}

	_, err = hnRepo.StoreTopStories(ctx, topStories)
	if err != nil {
		slog.Error("failed to store top stories", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("updated hackernews top stories")
}
