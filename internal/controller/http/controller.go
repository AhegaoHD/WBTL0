package http

import (
	"github.com/AhegaoHD/WBTL0/pkg/logger"
	"github.com/gorilla/mux"
)

func NewOrderController(l *logger.Logger, s OrderService) *mux.Router {
	controller := newController(l, s)
	router := mux.NewRouter()
	router.HandleFunc("/order", controller.GetOrderHandler).Methods("GET")
	return router
}
