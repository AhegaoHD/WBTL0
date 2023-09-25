package cache

import (
	"context"
	"sync"

	"github.com/AhegaoHD/WBTL0/internal/model"
)

type OrderRepo interface {
	GetOrdersWithDetails(ctx context.Context) ([]model.Order, error)
}

type OrderCache struct {
	mu    sync.RWMutex
	cache map[string]*model.Order
}

func NewOrderCache(orderRepo OrderRepo) *OrderCache {
	orderCache := &OrderCache{
		cache: make(map[string]*model.Order),
	}
	orders, _ := orderRepo.GetOrdersWithDetails(context.Background())
	for _, value := range orders {
		orderCache.Set(&value)
	}
	return orderCache
}

func (oc *OrderCache) Get(key string) (*model.Order, bool) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	value, ok := oc.cache[key]
	return value, ok
}

func (oc *OrderCache) Set(value *model.Order) {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	oc.cache[value.Order_uid] = value
}
