package main

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

// AmqpDispatcher is used as anchor for dispatch messsage method for real
// AMQP channels
type AmqpDispatcher struct {
	channel       queuePublishableChannel
	queueName     string
	mandatorySend bool
}

// NewAMQPDispatcher returns a new AMQP dispatcher wrapped around a single
// publishing channel.
func NewAMQPDispatcher(publishChannel queuePublishableChannel, name string, mandatory bool) *AmqpDispatcher {
	return &AmqpDispatcher{channel: publishChannel, queueName: name, mandatorySend: mandatory}
}

type queuePublishableChannel interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

// DispatchMessage implementation of dispatch message interface method
func (q *AmqpDispatcher) DispatchMessage(message interface{}) (err error) {
	fmt.Printf("Dispatching message to queue %s\n", q.queueName)
	body, err := json.Marshal(message)
	if err == nil {
		err = q.channel.Publish(
			"",              // exchange
			q.queueName,     // routing key
			q.mandatorySend, // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			fmt.Printf("Failed to dispatch message: %s\n", err)
		}
	} else {
		fmt.Printf("Failed to marshal message %v (%s)\n", message, err)
	}
	return
}
