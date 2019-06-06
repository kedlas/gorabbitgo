package gorabbitgo

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Client struct {
	conn    *Conn
	channel *amqp.Channel
	prepare Prepare
}

type ChanFuture struct {
	ch  *amqp.Channel
	err error
}

type Prepare func(ch *amqp.Channel) error

func NewClient(conn *Conn, prepare Prepare) *Client {
	return &Client{conn: conn, prepare: prepare}
}

func (cl *Client) isOK() bool {
	return cl.conn.isOK() && cl.channel != nil
}

func (cl *Client) GetChannel(ch chan *ChanFuture) {
	if cl.isOK() {
		ch <- &ChanFuture{ch: cl.channel, err: nil}
		return
	}

	cl.channel = nil

	// request the valid connection
	connCh := make(chan *ConnFuture, 1)
	cl.conn.GetConnection(connCh)
	cf := <-connCh

	if cf.err != nil {
		ch <- &ChanFuture{ch: nil, err: errors.Wrap(cf.err, "Cannot create channel on invalid connection")}
		return
	}

	channel, err := cf.conn.Channel()
	if err != nil {
		ch <- &ChanFuture{ch: nil, err: errors.Wrap(cf.err, "Cannot create channel")}
		return
	}

	cl.channel = channel

	err = cl.prepare(cl.channel)
	if err != nil {
		ch <- &ChanFuture{ch: nil, err: errors.Wrap(err, "Cannot prepare channel")}
	}

	ch <- &ChanFuture{ch: cl.channel, err: nil}
}
