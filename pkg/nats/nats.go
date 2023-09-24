package nats

import (
	"github.com/nats-io/stan.go"
)

type Nats struct {
	Conn stan.Conn
}

func New(clusterID, clientID string, URL *string) (*Nats, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(*URL))
	if err != nil {
		return nil, err
	}

	return &Nats{Conn: sc}, nil
}
