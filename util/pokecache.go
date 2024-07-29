package util

import (
	"time"
)

type Cache struct {
	Entries map[string]cacheEntry
	//	mu sync.Mutex to be used when i implement concurrency
	Duration int //time in seconds
}
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(duration int) *Cache {
	CreatedCache := &Cache{Duration: duration}
	go CreatedCache.reapLoop() //this is gonna cause concurrent read write
	return CreatedCache
}

func (c *Cache) Add(key string, val []byte) { //try as a method
	c.Entries[key] = cacheEntry{createdAt: time.Now(), val: val} //difference?
}

func (c *Cache) Get(key string) (cacheEntry, bool) { //returns an entry not just the byte array
	val, ok := c.Entries[key]
	if !ok {
		var Zero cacheEntry
		return Zero, false
	}
	return val, true
}

func (c *Cache) reapLoop() { //this might NEED to be concurrent or my programs gonna block right?
	for {
		for key, val := range c.Entries {
			refresh_time := val.createdAt.Add(time.Duration(c.Duration) * time.Second)
			if val.createdAt.After(refresh_time) {
				delete(c.Entries, key)
			}
		}
	}
}

//time.After?
