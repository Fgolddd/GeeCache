package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	mtx        sync.Mutex //the cache mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mtx.Lock()
	defer c.mtx.Unlock() //release the lock when the function exits

	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok //return a value of type ByteView
	}

	return
}
