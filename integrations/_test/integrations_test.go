package integrations_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/cloudnativego/drones-cmds/service"
	dronescommon "github.com/cloudnativego/drones-common"
	"github.com/streadway/amqp"
)

var (
	appEnv, _ = cfenv.Current()
	server    = NewServer(appEnv)

	telemetry1 = []byte("{\"drone_id\": \"abc1234\", \"battery\": 80, \"uptime\": 3200, \"core_temp\": 20 }")
	telemetry2 = []byte("{\"drone_id\": \"drone2\", \"battery\": 40, \"uptime\": 1200, \"core_temp\": 10 }")
)

func TestIntegration(t *testing.T) {
	fmt.Println("== Integration Test Scenario ==")

	consumeRabbit(t)

	telemetryReply, err := submitTelemetry(t, telemetry1)
	if err != nil {
		t.Errorf("Failed to submit telemetry: %s\n", err)
		return
	}
	if telemetryReply.DroneID != "abc1234" {
		t.Errorf("Failed to get a matching reply from the command server when submitting telemetry event: %+v\n", telemetryReply)
		return
	}

	telemetryReply2, err := submitTelemetry(t, telemetry2)
	if err != nil {
		t.Errorf("Failed to submit 2nd telemetry: %s\n", err)
		return
	}
	if telemetryReply2.DroneID != "drone2" {
		t.Errorf("Failed to get a matching reply from 2nd telemetry submit: %+v\n", telemetryReply2)
	}

	// TODO don't use a sleep, use channel synchronization
	time.Sleep(2000 * time.Millisecond)
}

/*
 * == Utility Functions
 */

func consumeRabbit(t *testing.T) {
	rabbitHost := os.Getenv("RABBITMQ_PORT_5672_TCP_ADDR")
	rabbitURL := fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitHost)

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		t.Errorf("Failed to dial rabbit: %s", err.Error())
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	defer ch.Close()

	telemetryQ, err := ch.QueueDeclare(
		"telemetry", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	// TODO: set up a real consumer...
	_, err = ch.Consume(
		telemetryQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	/*	go func() {
		for {
			select {
			case telemetryRaw := <-telemetryIn:
				reactTelemetry(telemetryRaw)
			}
		}
	}() */
}

func reactTelemetry(telemetryRaw amqp.Delivery) {
	var event dronescommon.TelemetryUpdatedEvent
	err := json.Unmarshal(telemetryRaw.Body, &event)
	if err == nil {
		fmt.Printf("Telemetry Received: %+v\n", event)
	} else {
		fmt.Printf("Failed to de-serialize raw telemetry from queue, %v\n", err)
	}
	return
}

func submitTelemetry(t *testing.T, body []byte) (reply dronescommon.TelemetryUpdatedEvent, err error) {
	rawReply, err := submitCommand(t, "/api/cmds/telemetry", body)
	var telemetryReply dronescommon.TelemetryUpdatedEvent
	if err == nil {
		err = json.Unmarshal(rawReply, &telemetryReply)
	} else {
		t.Errorf("Failed to submit telemetry: %+v", err)
		return
	}
	if err != nil {
		t.Errorf("Failed to submit telemetry : %s", err.Error())
	}
	reply = telemetryReply
	return
}

func submitCommand(t *testing.T, url string, body []byte) (rawReply []byte, err error) {
	recorder := httptest.NewRecorder()
	telemetryCommandRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return
	}
	server.ServeHTTP(recorder, telemetryCommandRequest)
	if recorder.Code != 201 {
		t.Errorf("Error submitting command to %s, expected 201 got %d", url, recorder.Code)
		err = fmt.Errorf("Error submitting command to %s, expected 201 got %d", url, recorder.Code)
		return
	}
	rawReply = recorder.Body.Bytes()
	fmt.Printf("Command reply: HTTP %d %d bytes\n", recorder.Code, len(rawReply))
	return
}
