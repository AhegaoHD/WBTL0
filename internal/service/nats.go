package service

import (
	"context"
	"encoding/json"
	"github.com/AhegaoHD/WBTL0/internal/model"
	"github.com/AhegaoHD/WBTL0/pkg/logger"
	"github.com/nats-io/stan.go"
)

type NatsRepo interface {
	Subscribe(subject, queueGroup string, handler stan.MsgHandler, options ...stan.SubscriptionOption) (stan.Subscription, error)
}

type OrderCache interface {
	Get(key string) (*model.Order, bool)
	Set(value *model.Order)
}
type NatsService struct {
	l          *logger.Logger
	natsRepo   NatsRepo
	orderCache OrderCache
	orderRepo  OrderRepo
}

func NewNatsService(l *logger.Logger, natsRepo NatsRepo, orderCache OrderCache, orderRepo OrderRepo) *NatsService {
	return &NatsService{l: l, natsRepo: natsRepo, orderCache: orderCache, orderRepo: orderRepo}
}

func (s *NatsService) StartListening(subject, queueGroup string) (stan.Subscription, error) {
	return s.natsRepo.Subscribe(subject, queueGroup, func(m *stan.Msg) {
		err := s.handleMessage(m)
		if err != nil {
			s.l.Error(err)
		}
	})
}

func (s *NatsService) handleMessage(m *stan.Msg) error {
	var newOrder model.Order
	err := json.Unmarshal(m.Data, &newOrder)
	if err != nil {
		return err
	}
	err = s.orderRepo.SetOrdersWithDetails(context.Background(), &newOrder)
	if err != nil {
		return err
	}
	s.orderCache.Set(&newOrder)
	return nil
}
