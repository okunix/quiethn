package hackernews

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
)

type HackerNewsClient interface {
	GetTopStories(ctx context.Context, limit int) ([]*NewsItem, error)
	GetItemById(ctx context.Context, id uint32) (*NewsItem, error)
}

type HackerNewsClientImpl struct {
	baseURL string
}

func NewHackerNewsClient(baseURL string) HackerNewsClient {
	return &HackerNewsClientImpl{
		baseURL: baseURL,
	}
}

var (
	getTopStoriesPath = "/v0/topstories.json"
	getItemByIdPath   = func(id uint32) string {
		return fmt.Sprintf("/v0/item/%d.json", id)
	}
)

func (client *HackerNewsClientImpl) GetTopStories(
	ctx context.Context,
	limit int,
) ([]*NewsItem, error) {
	if limit > 500 {
		return nil, errors.New("limit it too high, max 500")
	}

	getTopStoriesURL := client.baseURL + getTopStoriesPath
	req, err := http.NewRequestWithContext(ctx, "GET", getTopStoriesURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var ids []uint32
	err = json.NewDecoder(resp.Body).Decode(&ids)
	if err != nil {
		return nil, err
	}

	var mutex sync.Mutex
	newsItems := make([]*NewsItem, 0, limit)
	getStoriesContext, cancel := context.WithCancel(ctx)
	defer cancel()
	sem := make(chan struct{}, limit)
	for _, id := range ids {
		sem <- struct{}{}

		go func(ctx context.Context, itemId uint32) {
			defer func() { <-sem }()

			item, err := client.GetItemById(ctx, itemId)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				slog.Error("failed to get item by id", "id", id, "error", err.Error())
				return
			}
			if item.Type != "story" {
				return
			}

			mutex.Lock()
			newsItems = append(newsItems, item)
			mutex.Unlock()

		}(getStoriesContext, id)

		if len(newsItems) == limit {
			cancel()
			break
		}
	}

	return newsItems, nil
}

func (client *HackerNewsClientImpl) GetItemById(ctx context.Context, id uint32) (*NewsItem, error) {
	getItemByIdURL := client.baseURL + getItemByIdPath(id)
	req, err := http.NewRequestWithContext(ctx, "GET", getItemByIdURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var newsItemResponse NewsItemResponse
	err = json.NewDecoder(resp.Body).Decode(&newsItemResponse)
	if err != nil {
		return nil, err
	}
	newsItem, err := NewNewsItem(
		newsItemResponse.Id,
		newsItemResponse.Name,
		newsItemResponse.URL,
		newsItemResponse.Type,
	)
	return &newsItem, err
}
