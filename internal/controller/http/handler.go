package http

//import (
//	"encoding/json"
//	"github.com/AhegaoHD/WBTL0/app"
//	"github.com/AhegaoHD/WBTL0/internal/model"
//	"github.com/AhegaoHD/WBTL0/repository"
//	"github.com/nats-io/nats.go"
//	"log"
//)

//func SubscribeToNATS(myApp *app.App) {
//	myApp.NC.Subscribe("subject", func(msg *nats.Msg) {
//		var newOrder model.Order
//		err := json.Unmarshal(msg.Data, &newOrder)
//		if err != nil {
//			log.Fatal(err)
//		}
//		err = repository.SendOrder(newOrder, myApp.DB)
//		if err != nil {
//			log.Fatal(err)
//		}
//		repository.SaveOrderToCache(myApp.OrderCache, &newOrder)
//	})
//}
