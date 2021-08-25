package cache

import (
	"log"
	"sync"
	"time"

	"github.cwx.io/mrubiosan/mx51/weather"
)

// Cache is a Source decorator that caches reports.
type Cache struct {
	sync.Mutex
	Decorated weather.Source
	Ttl       time.Duration
	city      map[weather.Location]cacheEntry
}

type cacheEntry struct {
	expiry time.Time
	result weather.Report
}

// New creates a new Cache struct.
func New(source weather.Source, ttl time.Duration) *Cache {
	return &Cache{
		Decorated: source,
		Ttl:       ttl,
		city:      make(map[weather.Location]cacheEntry),
	}
}

func (c *cacheEntry) stale() bool {
	return c.expiry.Before(time.Now())
}

func (c *Cache) Report(loc weather.Location) (weather.Report, error) {
	c.Lock()
	defer c.Unlock()
	entry, ok := c.city[loc]

	if !ok || entry.stale() {
		log.Printf("Cache miss for %s", loc)
		res, err := c.Decorated.Report(loc)
		if err == nil {
			c.city[loc] = cacheEntry{expiry: time.Now().Add(c.Ttl), result: res}
			return res, err
		} else if !ok { // If it failed and was never populated, pass through the result
			return res, err
		}
		log.Printf("Failed to refresh cache for %s: %s. Returning stale copy.", loc, err.Error())
		entry.result.Stale = true
	} else {
		log.Printf("Cache hit for %s", loc)
	}

	return entry.result, nil
}
