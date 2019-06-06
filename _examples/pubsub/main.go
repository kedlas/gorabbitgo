package main

import (
"fmt"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
"gorabbitgo"
	"time"
)

func main() {
	var sent int
	var rcvd int


	logger := &gorabbitgo.StdOutLogger{}

	// Create connection and connect to rabbitMQ broker
	conn := gorabbitgo.NewConn("amqp://guest:guest@localhost:5672/")
	conn.SetLogger(logger)
	_, err := conn.Open()
	if err != nil {
		logger.Print(err)
	}

	// create the publisher with function that will be called whenever new underlying channel is created
	preparePub := func(ch *amqp.Channel) error {
		err := ch.ExchangeDeclare("gorabbitgo-ex", "fanout", true, false, false, false, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to prepare publisher's exchange")
		}
		return nil
	}
	pub := gorabbitgo.NewPublisher(conn, preparePub)

	// creates the consumer that will declare queue and binds it to exchange that publisher publishes to and starts consumption
	// note that exchange is defined in publisher and consumer prepare functions too, thus it must have the same settings
	prepareCons := func(ch *amqp.Channel) error {
		err := ch.ExchangeDeclare("gorabbitgo-ex", "fanout", true, false, false, false, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to prepare consumer's exchange")
		}
		_, err = ch.QueueDeclare("gorabbitgo-q", true, false, false, false, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to prepare consumer's exchange")
		}
		err = ch.QueueBind("gorabbitgo-q", "gorgo", "gorabbitgo-ex", true, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to prepare bind for consumer's queue and exchange")
		}
		return nil
	}
	reader := func(msg *amqp.Delivery) {
		content := string(msg.Body)

		logger.Print(fmt.Sprintf("Received '%s'", content))

		err = msg.Ack(false)
		if err != nil {
			logger.Print(fmt.Sprintf("Failed to ack. %v", err))
			return
		}

		rcvd++
		logger.Print(fmt.Sprintf("Acked '%s'", content))
	}
	cons := gorabbitgo.NewConsumer(conn, prepareCons, reader)

	// Start consuming here
	go func() {
		err = cons.Consume("gorabbitgo-q", "gorgo123", true)
		if err != nil {
			logger.Print(err)
		}
	}()

	// Publish some messages now
	for i := 1; i <= 5; i++ {
		msg := amqp.Publishing{Body: []byte(fmt.Sprintf("hello no #%d", i))}
		err = pub.Publish("gorabbitgo-ex", "gorgo", msg)
		if err != nil {
			fmt.Println(err)
			continue
		}

		sent++
	}

	t := time.NewTimer(time.Second * 2)
	<- t.C

	cons.Stop()
	conn.Close()

	logger.Print("=======================")
	logger.Print(fmt.Sprintf("Messages sent: %d", sent))
	logger.Print(fmt.Sprintf("Messages received: %d", rcvd))
}
