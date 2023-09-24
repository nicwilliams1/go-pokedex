package pokedex_api

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createAt time.Time
	val      []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

func (c *Cache) Add(key string, val []byte, ch chan<- bool) {

	// lock with mutex
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{
		createAt: time.Now(),
		val:      val,
	}

	c.entries[key] = entry
	ch <- true
}

func (c *Cache) Get(key string) ([]byte, bool) {

	// lock with mutex
	c.mu.Lock()
	defer c.mu.Unlock()

	// find entry
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		clear(c.entries)
		c.mu.Unlock()
	}

}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		entries: make(map[string]cacheEntry),
	}

	go c.reapLoop(interval)

	return &c
}
