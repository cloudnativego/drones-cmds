package integrations_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/cloudnativego/drones-cmds/service"
	dronescommon "github.com/cloudnativego/drones-common"
	"github.com/streadway/amqp"
)

var (
	appEnv, _ = cfenv.Current()
	server    = NewServer(appEnv)

	telemetry1     = []byte("{\"drone_id\": \"abc1234\", \"battery\": 80, \"uptime\": 3200, \"core_temp\": 20 }")
	telemetry2     = []byte("{\"drone_id\": \"drone2\", \"battery\": 40, \"uptime\": 1200, \"core_temp\": 10 }")
	doneTelemetry  = make(chan error)
	telemetryCount = 0

	alert1     = []byte("{\"drone_id\": \"abc1234\", \"fault_code\": 1, \"description\": \"super fail\"}")
	alert2     = []byte("{\"drone_id\": \"drone2\", \"fault_code\": 2, \"description\": \"overheating\"}")
	doneAlert  = make(chan error)
	alertCount = 0

	position1     = []byte("{\"drone_id\": \"abc1234\", \"latitude\": 31.01, \"longitude\": 72.5, \"altitude\": 3500.12, \"current_speed\": 15.12, \"heading_cardinal\": 0}")
	position2     = []byte("{\"drone_id\": \"drone2\", \"latitude\": 31.01, \"longitude\": 72.5, \"altitude\": 3500.12, \"current_speed\": 15.12, \"heading_cardinal\": 0}")
	position3     = []byte("{\"drone_id\": \"drone3\", \"latitude\": 31.01, \"longitude\": 72.5, \"altitude\": 3500.12, \"current_speed\": 15.12, \"heading_cardinal\": 0}")
	donePosition  = make(chan error)
	positionCount = 0
)

func TestIntegration(t *testing.T) {
	fmt.Println("== Integration Test Scenario ==")

	telemetryReply, err := submitTelemetry(t, telemetry1)
	if err != nil {
		t.Errorf("Failed to submit telemetry: %s\n", err)
		return
	}
	if telemetryReply.DroneID != "abc1234" {
		t.Errorf("Failed to get a matching reply from the command server when submitting telemetry comand: %+v\n", telemetryReply)
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

	alertReply, err := submitAlert(t, alert1)
	if err != nil {
		t.Errorf("Failed to submit an alert command: %s\n", err.Error())
		return
	}
	if alertReply.DroneID != "abc1234" {
		t.Errorf("Failed to get matching reply from submitting alert command: %+v\n", alertReply)
		return
	}

	alertReply2, err := submitAlert(t, alert2)
	if err != nil {
		t.Errorf("Failed to submit 2nd alert: %s\n", err.Error())
		return
	}
	if alertReply2.DroneID != "drone2" {
		t.Errorf("Expecting matching reply from submitting 2nd alert command, got %+v\n", alertReply2)
		return
	}

	positionReply, _ := submitPosition(t, position1)
	positionReply2, _ := submitPosition(t, position2)
	positionReply3, _ := submitPosition(t, position3)
	if positionReply.DroneID != "abc1234" {
		t.Errorf("Got wrong reply from position submit: %+v\n", positionReply)
		return
	}
	if positionReply2.DroneID != "drone2" {
		t.Errorf("Got wrong reply from position submit, expected drone2, got %+vs\n", positionReply2)
		return
	}
	if positionReply3.DroneID != "drone3" {
		t.Errorf("Got wrong reply from position submit, expected drone3, got %+v\n", positionReply3)
		return
	}

	consumeRabbit(t)
	<-doneTelemetry
	<-doneAlert
	<-donePosition
	if telemetryCount != 2 {
		t.Errorf("Didn't dequeue 2 telemetry events.\n")
	}
	if alertCount != 2 {
		t.Errorf("Didn't dequeue 2 alerts.\n")
	}
	if positionCount != 3 {
		t.Error("Didn't dequeue 3 positions.\n")
	}
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

	alertQ, err := ch.QueueDeclare(
		"alerts", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	positionsQ, err := ch.QueueDeclare(
		"positions", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	telemetryIn, err := ch.Consume(
		telemetryQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	alertsIn, err := ch.Consume(
		alertQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	positionsIn, err := ch.Consume(
		positionsQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	go func() {
		for d := range telemetryIn {
			reactTelemetry(d)
		}
		doneTelemetry <- nil
		for a := range alertsIn {
			reactAlert(a)
		}
		doneAlert <- nil
		for p := range positionsIn {
			reactPosition(p)
		}
		donePosition <- nil
	}()
}

func reactTelemetry(telemetryRaw amqp.Delivery) {
	var event dronescommon.TelemetryUpdatedEvent
	err := json.Unmarshal(telemetryRaw.Body, &event)
	if err == nil {
		fmt.Printf("Telemetry Received: %+v\n", event)
		telemetryCount++
	} else {
		fmt.Printf("Failed to de-serialize raw telemetry from queue, %v\n", err)
	}
	return
}

func reactPosition(positionRaw amqp.Delivery) {
	var event dronescommon.PositionChangedEvent
	err := json.Unmarshal(positionRaw.Body, &event)
	if err == nil {
		fmt.Printf("Position received: %+v\n", event)
		positionCount++
	} else {
		fmt.Printf("Failed to de-serialize raw alert from queue, %v\n", err)
	}
}

func reactAlert(alertRaw amqp.Delivery) {
	var event dronescommon.AlertSignalledEvent
	err := json.Unmarshal(alertRaw.Body, &event)
	if err == nil {
		fmt.Printf("Alert received: %+v\n", event)
		alertCount++
	} else {
		fmt.Printf("Failed to de-serialize raw alert from queue, %v\n", err)
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

func submitAlert(t *testing.T, body []byte) (reply dronescommon.AlertSignalledEvent, err error) {
	rawReply, err := submitCommand(t, "/api/cmds/alerts", body)
	var alertReply dronescommon.AlertSignalledEvent
	if err == nil {
		err = json.Unmarshal(rawReply, &alertReply)
	}
	if err != nil {
		t.Errorf("Failed to submit alert: %s\n", err)
		return
	}
	reply = alertReply
	return
}

func submitPosition(t *testing.T, body []byte) (reply dronescommon.PositionChangedEvent, err error) {
	rawReply, err := submitCommand(t, "/api/cmds/positions", body)
	var positionReply dronescommon.PositionChangedEvent
	if err == nil {
		err = json.Unmarshal(rawReply, &positionReply)
	}
	if err != nil {
		t.Errorf("Failed to submit position: %s\n", err)
		return
	}
	reply = positionReply
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
