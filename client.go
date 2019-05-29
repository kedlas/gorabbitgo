package gorabbitgo

import (
	"fmt"
	"github.com/streadway/amqp"
)

type Client struct {
	conn *Conn
	channel *amqp.Channel
}

func NewClient(conn *Conn) *Client {
	return &Client{conn:conn}
}

func (cl *Client) isOK() bool {
	return cl.conn.isOK() && cl.channel != nil
}

func (cl *Client) GetChannel(ch chan *amqp.Channel) {
	if cl.isOK() {
		ch <- cl.channel
		return
	}

	cl.channel = nil

	// request the valid connection
	connCh := make(chan *amqp.Connection, 1)
	cl.conn.GetConnection(connCh)
	conn := <-connCh

	channel, err := conn.Channel()
	if err != nil {
		fmt.Printf("Cannot create channel %v", err)
		cl.GetChannel(ch)
		return
	}

	cl.channel = channel
	ch <- channel
}
