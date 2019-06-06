package gorabbitgo

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Consumer struct {
	cl *Client
	read Read
	quit chan bool
}

type Read func(msg *amqp.Delivery)

func NewConsumer(conn *Conn, prepare Prepare, read Read) *Consumer {
	return &Consumer{cl: NewClient(conn, prepare), read: read}
}

func (c *Consumer) Consume(queue string, id string, exclusive bool) error {
	c.quit = make(chan bool, 1)

	ch := make(chan *ChanFuture, 1)
	c.cl.GetChannel(ch)
	chf := <-ch

	if chf.err != nil {
		return errors.Wrap(chf.err, "Failed to start consume")
	}

	msgs, err := chf.ch.Consume(queue, id, false, exclusive, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to register consumer")
	}

	for {
		select {
		case msg, ok := <- msgs:
			if ok {
				c.read(&msg)
			}
		case <- c.quit:
			_ = chf.ch.Cancel(id, false)
			break
		}
	}
}

func (c *Consumer) Stop() {
	c.quit <- true
}
