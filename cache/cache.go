package cache

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	mutex  sync.Mutex
	memory map[string]cacheItem
}

type cacheItem struct {
	value          any
	expirationTime time.Time
}

func NewCache() *Cache {
	return &Cache{
		memory: make(map[string]cacheItem),
	}
}

var (
	ErrNilKey = errors.New("key is nil")
)

func (c *Cache) Get(key string) (any, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, found := c.memory[key]
	if !found {
		return nil, ErrNilKey
	}

	if time.Now().After(item.expirationTime) {
		delete(c.memory, key)
		return nil, ErrNilKey
	}

	return item.value, nil
}

func (c *Cache) Put(key string, value any, TTL int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.memory[key] = cacheItem{
		value:          value,
		expirationTime: time.Now().Add(time.Duration(TTL) * time.Second),
	}
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.memory, key)
}

func (c *Cache) Flush() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.memory = make(map[string]cacheItem)
}
