package queue_subscribe

import "github.com/nats-io/stan.go"

type NatsRepository struct {
	conn stan.Conn
}

func NewNatsRepository(conn stan.Conn) *NatsRepository {
	return &NatsRepository{conn: conn}
}

func (r *NatsRepository) Subscribe(subject, queueGroup string, handler stan.MsgHandler, options ...stan.SubscriptionOption) (stan.Subscription, error) {
	return r.conn.QueueSubscribe(subject, queueGroup, handler, options...)
}
