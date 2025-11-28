package hackernews

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type HackerNewsRepo interface {
	GetTopStories(ctx context.Context) ([]*NewsItem, error)
	GetItemById(ctx context.Context, id uint32) (*NewsItem, error)
	AddItem(ctx context.Context, item NewsItem) (*NewsItem, error)
	DeleteItemById(ctx context.Context, id uint32) (*NewsItem, error)
	StoreTopStories(ctx context.Context, stories []*NewsItem) ([]*NewsItem, error)
}

type HackerNewsRepoRedis struct {
	rdb *redis.Client
}

func NewHackerNewsRepo(rdb *redis.Client) HackerNewsRepo {
	return &HackerNewsRepoRedis{
		rdb: rdb,
	}
}

var (
	itemKey = func(id uint32) string {
		return fmt.Sprintf("newsItem:%d", id)
	}
	topStoriesKey = "topStories"
)

func (h *HackerNewsRepoRedis) AddItem(ctx context.Context, item NewsItem) (*NewsItem, error) {
	key := itemKey(item.Id)
	valueBytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	value := string(valueBytes)
	err = h.rdb.Set(ctx, key, value, 0).Err()
	return &item, err
}

func (h *HackerNewsRepoRedis) DeleteItemById(ctx context.Context, id uint32) (*NewsItem, error) {
	key := itemKey(id)
	item, err := h.GetItemById(ctx, id)
	if err != nil {
		return nil, err
	}
	err = h.rdb.Del(ctx, key).Err()
	return item, err
}

func (h *HackerNewsRepoRedis) GetItemById(ctx context.Context, id uint32) (*NewsItem, error) {
	key := itemKey(id)
	valueJson, err := h.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var value NewsItem
	err = json.Unmarshal([]byte(valueJson), &value)
	return &value, err
}

func (h *HackerNewsRepoRedis) GetTopStories(ctx context.Context) ([]*NewsItem, error) {
	stories := []*NewsItem{}
	key := topStoriesKey
	valueJson, err := h.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return stories, nil
		}
		return nil, err
	}
	err = json.Unmarshal([]byte(valueJson), &stories)
	return stories, err
}

// for now i use simple Set with json encoded string but further i want move to SortedSet and HashMaps
func (h *HackerNewsRepoRedis) StoreTopStories(
	ctx context.Context,
	stories []*NewsItem,
) ([]*NewsItem, error) {
	key := topStoriesKey
	valueBytes, err := json.Marshal(stories)
	if err != nil {
		return nil, err
	}
	value := string(valueBytes)
	err = h.rdb.Set(ctx, key, value, 0).Err()
	return stories, err
}
