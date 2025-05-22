package fan_cache

import (
	strategies2 "scav.abc/fantastic-cache/fan-cache/strategies"
	"sync"
)

type cache struct {
	mu         sync.RWMutex
	lru        strategies2.Cache
	cacheBytes int
}

func (c *cache) add(k string, v ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Lazy Initialization
	if c.lru == nil {
		c.lru = strategies2.NewLruCache(c.cacheBytes, nil)
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
