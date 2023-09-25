package cache

import (
	"context"
	"sync"

	"github.com/AhegaoHD/WBTL0/internal/model"
)

type OrderRepo interface {
	GetOrderWithDetailsByUid(ctx context.Context, orderUID string) (model.Order, error)
	GetOrdersWithDetails(ctx context.Context) ([]model.Order, error)
	GetAllOrders(ctx context.Context) ([]model.Order, error)
	GetDelivery(ctx context.Context, orderId string) (model.Delivery, error)
	GetPayment(ctx context.Context, orderId string) (model.Payment, error)
	GetItems(ctx context.Context, orderId string) ([]model.Item, error)
	SetOrdersWithDetails(ctx context.Context, order *model.Order) error
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
