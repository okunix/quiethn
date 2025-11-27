package hackernews

import (
	"context"
	"fmt"

	"github.com/okunix/quiethn/cache"
)

type HackerNewsClientWithCache struct {
	parent HackerNewsClient
	cache  *cache.Cache
}

func NewHackerNewsClientWithCache(parent HackerNewsClient) HackerNewsClient {
	return &HackerNewsClientWithCache{
		parent: parent,
		cache:  cache.NewCache(),
	}
}

func (h *HackerNewsClientWithCache) GetItemById(
	ctx context.Context,
	id uint32,
) (*NewsItem, error) {
	key := fmt.Sprintf("NewsItem:%d", id)
	newsItemFromCache, err := h.cache.Get(key)
	if err == nil {
		return newsItemFromCache.(*NewsItem), nil
	}

	newsItem, err := h.parent.GetItemById(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		h.cache.Put(key, newsItem, 30)
	}()

	return newsItem, nil
}

func (h *HackerNewsClientWithCache) GetTopStories(
	ctx context.Context,
	limit int,
) ([]*NewsItem, error) {
	key := "TopStories"
	storiesFromCache, err := h.cache.Get(key)
	if err == nil {
		return storiesFromCache.([]*NewsItem), nil
	}

	stories, err := h.parent.GetTopStories(ctx, limit)
	if err != nil {
		return nil, err
	}

	go func() {
		h.cache.Put(key, stories, 30)
	}()

	return stories, nil
}
