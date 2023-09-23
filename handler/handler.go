package handler

import (
	"WB_Tech_level_0/app"
	"WB_Tech_level_0/model"
	"WB_Tech_level_0/repository"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
)

func GetOrderHandler(w http.ResponseWriter, r *http.Request, myApp *app.App) {
	orderID := r.URL.Query().Get("id")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	order, found := myApp.OrderCache.Get(orderID)

	if !found {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func SubscribeToNATS(myApp *app.App) {
	myApp.NC.Subscribe("subject", func(msg *nats.Msg) {
		var newOrder model.Order
		err := json.Unmarshal(msg.Data, &newOrder)
		if err != nil {
			log.Fatal(err)
		}
		err = repository.SendOrder(newOrder, myApp.DB)
		if err != nil {
			log.Fatal(err)
		}
		repository.SaveOrderToCache(myApp.OrderCache, &newOrder)
	})
}
