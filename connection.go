package gorabbitgo

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type Conn struct {
	Url string
	conn *amqp.Connection
}

func (c *Conn) Open() (*amqp.Connection, error) {
	if c.isOK() {
		return c.conn, nil
	}

	c.conn = nil

	conn, err := amqp.Dial(c.Url)
	if err != nil {
		// todo - get rid of logging
		fmt.Printf("Failed to Open connection to rabbitMQ broker. %v", err)
		return conn, err
	}

	c.conn = conn

	return conn, nil
}

func (c *Conn) Close() {
	if c.conn == nil {
		return
	}

	err := c.conn.Close()
	if err != nil {
		// todo - get rid of logging
		fmt.Printf("Failed to Close connection to rabbitMQ broker. %v", err)
	}
}

func (c *Conn) isOK() bool {
	return c.conn != nil && !c.conn.IsClosed()
}

func (c *Conn) GetConnection(ch chan *amqp.Connection) {
	if c.isOK() {
		ch <- c.conn
		return
	}

	_, err := c.Open()
	if err == nil && c.isOK() {
		c.GetConnection(ch)
		return
	}

	// todo - wait in background instead of blocking sleep
	time.Sleep(time.Second * 5)
	c.GetConnection(ch)
}
