package http

import (
	"context"
	"encoding/json"
	"github.com/AhegaoHD/WBTL0/internal/model"
	"github.com/AhegaoHD/WBTL0/pkg/logger"
	"net/http"
)

type OrderService interface {
	GetOrderWithDetailsByUid(ctx context.Context, orderUID string) (model.Order, error)
}

type controller struct {
	l            *logger.Logger
	orderService OrderService
}

func newController(l *logger.Logger, orderService OrderService) *controller {
	return &controller{
		l:            l,
		orderService: orderService,
	}
}

func (c *controller) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("id")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	order, err := c.orderService.GetOrderWithDetailsByUid(r.Context(), orderID)

	if err != nil {
		http.Error(w, `{"data":"Order not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
