package infrastructure

import (
	"github.com/nats-io/nats.go"
)

func InitNats() *nats.Conn {
	nc, err := nats.Connect("127.0.0.1:4222")
	if err != nil {
		panic(err)
	}
	return nc
}
