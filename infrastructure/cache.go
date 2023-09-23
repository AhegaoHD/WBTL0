package infrastructure

import (
	"WB_Tech_level_0/model"
	"database/sql"
	"sync"
)

type OrderCache struct {
	mu    sync.RWMutex
	cache map[string]*model.Order
}

func InitCache(db *sql.DB) *OrderCache {
	orderCache := &OrderCache{
		cache: make(map[string]*model.Order),
	}
	return orderCache
}

func (oc *OrderCache) Get(key string) (*model.Order, bool) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	value, ok := oc.cache[key]
	return value, ok
}

func (oc *OrderCache) Set(key string, value *model.Order) {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	oc.cache[key] = value
}
