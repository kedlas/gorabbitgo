package gorabbitgo

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"time"
)

type Conn struct {
	url              string
	conn             *amqp.Connection
	reconnectCounter int
	reconnectLimit   int
	logger           Logger
}

type ConnFuture struct {
	conn *amqp.Connection
	err  error
}

func NewConn(url string) *Conn {
	return &Conn{url: url, logger: &NullLogger{}, reconnectLimit: 10}
}

func (c *Conn) Open() (*amqp.Connection, error) {
	if c.isOK() {
		return c.conn, nil
	}

	c.logger.Print("Opening connection to rabbitMQ...")

	c.conn = nil
	conn, err := amqp.Dial(c.url)
	if err != nil {
		c.logger.Print(fmt.Sprintf("Failed to Open connection to rabbitMQ broker. %v", err))

		c.reconnectCounter = c.reconnectCounter + 1

		return conn, err
	}

	c.logger.Print("Connection to rabbitMQ opened")

	c.reconnectCounter = 0
	c.conn = conn

	return conn, nil
}

func (c *Conn) Close() {
	if c.conn == nil {
		return
	}

	c.logger.Print("Closing connection to rabbitMQ")
	err := c.conn.Close()
	if err != nil {
		c.logger.Print(fmt.Sprintf("Failed to close connection. %v", err))
	}
}

func (c *Conn) isOK() bool {
	return c.conn != nil && !c.conn.IsClosed()
}

func (c *Conn) GetConnection(ch chan *ConnFuture) {
	if !c.isOK() {
		c.reconnect(ch)
		return
	}

	ch <- &ConnFuture{conn: c.conn, err: nil}
}

func (c *Conn) reconnect(ch chan *ConnFuture) {
	_, err := c.Open()

	// connection established
	if err == nil && c.isOK() {
		ch <- &ConnFuture{conn: c.conn, err: nil}
		return
	}

	// max reconnect count limit reached
	if c.reconnectCounter > c.reconnectLimit {
		ch <- &ConnFuture{conn: c.conn, err: errors.New("amqp max reconnect count limit reached")}
	}

	rTimer := time.NewTimer(time.Second * 5)
	<-rTimer.C

	c.reconnect(ch)
}

func (c *Conn) SetLogger(logger Logger) {
	c.logger = logger
}
