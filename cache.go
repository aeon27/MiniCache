package minicache

import (
	"MiniCache/lru"
	"sync"
)

type cache struct {
	mu         *sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add()
