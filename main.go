package main

import (
	"WB_Tech_level_0/app"
	"WB_Tech_level_0/handler"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	myApp := &app.App{}
	myApp.Init()

	handler.SubscribeToNATS(myApp)

	http.HandleFunc("/getOrder", func(w http.ResponseWriter, r *http.Request) {
		handler.GetOrderHandler(w, r, myApp)
	})

	http.ListenAndServe(":3333", nil)

}
