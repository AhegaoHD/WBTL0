package service

import (
	"context"

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

type OrderService struct {
	orderRepo OrderRepo
}

func NewOrderService(
	or OrderRepo,
) *OrderService {

	return &OrderService{
		orderRepo: or,
	}
}

func (s *OrderService) GetOrdersWithDetails(ctx context.Context) ([]model.Order, error) {
	orders, err := s.orderRepo.GetOrdersWithDetails(ctx)
	if err != nil {
		return nil, err
	}
	return orders, err
}

func (s *OrderService) GetOrderWithDetailsByUid(ctx context.Context, orderUID string) (model.Order, error) {
	order, err := s.orderRepo.GetOrderWithDetailsByUid(ctx, orderUID)
	if err != nil {
		return model.Order{}, err
	}
	return order, err
}
