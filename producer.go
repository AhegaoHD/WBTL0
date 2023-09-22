package main

import (
	"WB_Tech_level_0/pg"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"os"
	"strconv"
	"time"
)

func main() {
	nc, err := nats.Connect("127.0.0.1:4222")
	if err != nil {
		panic(err)
	}
	defer nc.Close()
	order := getExampleOrder()
	id := order.Order_uid
	subject := "subject"

	for i := 1; ; i++ {
		order.Order_uid = id + strconv.Itoa(i)
		orderM, _ := json.Marshal(order)
		err := nc.Publish(subject, orderM)
		fmt.Printf("send %s \n", order.Order_uid)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(2 * time.Second)
	}

}
func getExampleOrder() pg.Order {
	var ord pg.Order
	exampleFile, err := os.Open("example.json")
	jsonParser := json.NewDecoder(exampleFile)
	if err = jsonParser.Decode(&ord); err != nil {
		fmt.Println("parsing config file", err.Error())
	}
	return ord
}
