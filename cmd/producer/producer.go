package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AhegaoHD/WBTL0/internal/model"
	"github.com/nats-io/stan.go"

	"os"
	"strconv"
	"time"
)

func main() {
	sc, err := stan.Connect("test-cluster", "publisher-1", stan.NatsURL("127.0.0.1:4222"))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}
	defer sc.Close()

	order := getExampleOrder()
	subject := "subject"
	id := order.Order_uid
	for i := 1; ; i++ {
		order.Order_uid = id + strconv.Itoa(i)
		order.Items[0].ChrtId = i
		orderM, _ := json.Marshal(order)
		err := sc.Publish(subject, orderM)
		fmt.Printf("send %s \n", order.Order_uid)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(2 * time.Second)
	}

}

func getExampleOrder() model.Order {
	var ord model.Order
	exampleFile, err := os.Open("example.json")
	jsonParser := json.NewDecoder(exampleFile)
	if err = jsonParser.Decode(&ord); err != nil {
		fmt.Println("parsing config file", err.Error())
	}
	return ord
}
