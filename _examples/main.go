package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"gorabbitgo"
)

func main() {

	c := &gorabbitgo.Conn{
		Url: "amqp://guest:guest@localhost:5672/",
	}

	_, err := c.Open()
	if err != nil {
		fmt.Println(err)
	}

	pub := gorabbitgo.NewPublisher(c)

	msg := amqp.Publishing{Body: []byte("hello")}
	err = pub.Publish("test", "some key", msg)
	if err != nil {
		fmt.Println(err)
	}
}
