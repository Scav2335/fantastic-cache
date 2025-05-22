package fantastic_cache

import (
	"scav.abc/fantastic-cache/strategies"
	"sync"
)

type cache struct {
	mu         sync.RWMutex
	lru        strategies.Cache
	cacheBytes int
}

func (c *cache) add(k string, v ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Lazy Initialization
	if c.lru == nil {
		c.lru = strategies.NewLruCache(c.cacheBytes, nil)
	}
	c.lru.Set(k, v)
}

func (c *cache) get(k string) (v ByteView, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(k); ok {
		return v.(ByteView), true
	}

	return
}
