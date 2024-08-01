package util

import (
	"fmt"
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
	val       any //TODO change this to an interface
}

func NewCache(duration time.Duration) *Cache {
	CreatedCache := &Cache{
		Entries:  make(map[string]cacheEntry),
		Duration: duration,
	}
	go CreatedCache.reapLoop()
	return CreatedCache
}

func (c *Cache) Add(key string, val any) {
	c.mu.Lock() //TODO: Understand the concurrent elements of this program more
	defer c.mu.Unlock()
	c.Entries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.Entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() { //TODO: Make sure I understand why this doesn't block the program; go?
	ticker := time.NewTicker(time.Duration(c.Duration) / 2)
	defer ticker.Stop()

	for range ticker.C { //we can write for range like this for channels where we don't need the value
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.Entries {
			if now.Sub(entry.createdAt) > time.Duration(c.Duration) {
				fmt.Println("Deleting old data...")
				delete(c.Entries, key)
			}
		}
		c.mu.Unlock()
	}
}
