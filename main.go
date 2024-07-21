package main

import (
	"fmt"
	"sync"
	"time"
)

// CacheItem represent a single cache item
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

// Cache repesents the in-memory cache
type Cache struct {
	Items map[string]CacheItem
	mutex sync.RWMutex
}

// NewCache create a new Cache instance
func NewCache() *Cache {
	return &Cache{
		Items: make(map[string]CacheItem),
	}
}

// Set add an item to the cache with an optional expiration time (in seconds)
func (c *Cache) Set(key string, value interface{}, duration int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Unix() + duration
	}
	c.Items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	item, exist := c.Items[key]
	if !exist {
		return nil, exist
	}
	exp := item.Expiration
	if exp > 0 && time.Now().Unix() > exp {
		delete(c.Items, key)
		return nil, false
	}
	return item.Value, true
}

func main() {
	cache := NewCache()
	cache.Set("key1", 10, 20)
	for {
		fmt.Println(cache.Get("key1"))
		time.Sleep(1 * time.Second)
	}

}
