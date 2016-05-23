package integrations_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/cloudnativego/drones-cmds/service"
	"github.com/streadway/amqp"
)

type fakeMessage struct {
	a string
	b string
}

func TestAmqpDispatcherSubmitsToQueue(t *testing.T) {
	rabbitHost := os.Getenv("RABBITMQ_PORT_5672_TCP_ADDR")
	rabbitURL := "amqp://guest:guest@192.168.99.100:5672/" // default fallback if not using wercker env vars

	if len(rabbitHost) > 0 {
		rabbitURL = fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitHost)
	}
	fmt.Printf("\nUsing URL (%s) for Rabbit.\n", rabbitURL)

	conn, err := amqp.Dial(rabbitURL)
	failOnError(t, err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(t, err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(t, err, "Failed to declare a queue")
	dispatcher := NewAMQPDispatcher(ch, q.Name, true) // require delivery ack for testing
	fmt.Println("About to dispatch message to queue...")
	err = dispatcher.DispatchMessage(fakeMessage{a: "hello", b: "world"})
	failOnError(t, err, "Failed to dispatch message on channel/queue")
	fmt.Println("'spatched.")
}

func failOnError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Errorf("%s: %s", msg, err)
	}
}
