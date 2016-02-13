package main

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

// AmqpDispatcher is used as anchor for dispatch messsage method for real
// AMQP channels
type AmqpDispatcher struct {
	channel   queuePublishableChannel
	queueName string
}

// NewAMQPDispatcher returns a new AMQP dispatcher wrapped around a single
// publishing channel.
func NewAMQPDispatcher(publishChannel queuePublishableChannel, name string) *AmqpDispatcher {
	return &AmqpDispatcher{channel: publishChannel, queueName: name}
}

type queuePublishableChannel interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

// DispatchMessage implementation of dispatch message interface method
func (q *AmqpDispatcher) DispatchMessage(message interface{}) (err error) {
	body, err := json.Marshal(message)
	if err != nil {
		err = q.channel.Publish(
			"",          // exchange
			q.queueName, // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
	}
	return
}
