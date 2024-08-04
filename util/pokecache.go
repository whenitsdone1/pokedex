package util

import (
	"sync"
	"time"
)

type Cache struct {
	Entries  map[string]cacheEntry
	mu       sync.RWMutex
	Duration time.Duration //How fast do we want our cache to refresh?
}
type cacheEntry struct {
	createdAt time.Time
	val       Parseable
}

func NewCache(duration time.Duration) *Cache {
	CreatedCache := &Cache{
		Entries:  make(map[string]cacheEntry),
		Duration: duration,
	}
	go CreatedCache.reapLoop()
	return CreatedCache
}

func (c *Cache) Add(key string, val Parseable) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Entries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) (Parseable, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.Entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(time.Duration(c.Duration))
	defer ticker.Stop()
	for range ticker.C { //we can write for range like this for channels where we don't need the value
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.Entries {
			if now.Sub(entry.createdAt) > time.Duration(c.Duration) {
				delete(c.Entries, key)
			}
		}
		c.mu.Unlock()
	}
}
