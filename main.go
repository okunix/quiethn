package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"strconv"

	"github.com/okunix/quietHN/hackernews"
)

//go:embed templates/*.html
var templateFS embed.FS

type Validatable interface {
	Validate(ctx context.Context) (problems map[string]string)
}

func GetenvWithDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func main() {
	topStoriesLimitEnv := GetenvWithDefault("HN_TOP_STORIES_LIMIT", "30")
	topStoriesLimit, err := strconv.Atoi(topStoriesLimitEnv)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	hn := hackernews.NewHackerNewsClient("https://hacker-news.firebaseio.com")
	hnCache := hackernews.NewHackerNewsClientWithCache(hn)
	// first request
	stories, err := hnCache.GetTopStories(ctx, topStoriesLimit)
	if err != nil {
		panic(err)
	}
	for _, i := range stories {
		fmt.Printf("%+v\n", i)
	}
	fmt.Println()
	// second request
	stories, err = hnCache.GetTopStories(ctx, topStoriesLimit)
	if err != nil {
		panic(err)
	}
	for _, i := range stories {
		fmt.Printf("%+v\n", i)
	}
}
