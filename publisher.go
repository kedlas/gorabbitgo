package gorabbitgo

import "github.com/streadway/amqp"

type Publisher struct {
	cl *Client
}

func NewPublisher(conn *Conn) *Publisher {
	return &Publisher{cl: NewClient(conn)}
}

func (p *Publisher) Publish(exchange string, key string, msg amqp.Publishing) error {
	ch := make(chan *amqp.Channel, 1)
	p.cl.GetChannel(ch)
	channel := <-ch

	return channel.Publish(exchange, key, false, false, msg)
}
