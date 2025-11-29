package hackernews

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"time"
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

	newsItems := make([]*NewsItem, 0, limit)
	for _, id := range ids {
		item, err := client.GetItemById(ctx, id)
		if err != nil {
			slog.Error("failed to get item by id", "id", id, "error", err.Error())
			continue
		}
		if item.Type != "story" {
			continue
		}
		newsItems = append(newsItems, item)
		if len(newsItems) == limit {
			break
		}
	}

	// sorting
	slices.SortFunc(newsItems, func(a *NewsItem, b *NewsItem) int {
		if a.Time.Sub(b.Time) > 0 {
			return -1
		} else if a.Time.Sub(b.Time) < 0 {
			return 1
		}
		return 0
	})

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
	if newsItemResponse.Deleted {
		return nil, fmt.Errorf("item %d is deleted", id)
	}
	newsItem, err := NewNewsItem(
		newsItemResponse.Id,
		newsItemResponse.Title,
		newsItemResponse.URL,
		newsItemResponse.Type,
		time.Unix(newsItemResponse.Time, 0),
	)
	return &newsItem, err
}
