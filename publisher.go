package gorabbitgo

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Publisher struct {
	cl *Client
}

func NewPublisher(conn *Conn, prepare Prepare) *Publisher {
	return &Publisher{cl: NewClient(conn, prepare)}
}

func (p *Publisher) Publish(exchange string, key string, msg amqp.Publishing) error {
	ch := make(chan *ChanFuture, 1)
	p.cl.GetChannel(ch)
	chf := <-ch

	if chf.err != nil {
		return errors.Wrap(chf.err, "Failed to publish")
	}

	return chf.ch.Publish(exchange, key, false, false, msg)
}
