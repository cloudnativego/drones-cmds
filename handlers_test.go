package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudnativego/drones-cmds/fakes"
	dronescommon "github.com/cloudnativego/drones-common"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func MakeTestServer(dispatcher queueDispatcher) *negroni.Negroni {
	server := negroni.New() // don't need all the middleware here or logging.
	mx := mux.NewRouter()
	initRoutes(mx, formatter, dispatcher)
	server.UseHandler(mx)
	return server
}

func TestAddValidTelemetryCreatesCommand(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	dispatcher := fakes.NewFakeQueueDispatcher()

	server := MakeTestServer(dispatcher)
	recorder = httptest.NewRecorder()
	body := []byte("{\"drone_id\":\"drone666\", \"battery\": 72, \"uptime\": 6941, \"core_temp\": 21 }")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/telemetry", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new telemetry item to return 201, got %d", recorder.Code)
	}
	if dispatcher.DispatchCount != 1 {
		t.Errorf("Expected queue dispatch count of 1, got %d", dispatcher.DispatchCount)
	}
	if len(dispatcher.Messages["telemetry"]) != 1 {
		t.Errorf("Expected telemetry message count of 1, got %d", len(dispatcher.Messages["telemetry"]))
	}

	var telemetryResponse dronescommon.TelemetryUpdatedEvent
	payload := recorder.Body.Bytes()
	err := json.Unmarshal(payload, &telemetryResponse)
	if err != nil {
		t.Errorf("Could not unmarshal payload into telemetry response object")
	}
	if telemetryResponse.DroneID != "drone666" {
		t.Errorf("Expected drone ID of 'drone666' got %s", telemetryResponse.DroneID)
	}
	if telemetryResponse.Uptime != 6941 {
		t.Errorf("Expected drone uptime of 6941, got %d", telemetryResponse.Uptime)
	}
}

func TestAddInvalidTelemetryReturnsBadRequest(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	dispatcher := fakes.NewFakeQueueDispatcher()
	server := MakeTestServer(dispatcher)
	recorder = httptest.NewRecorder()
	body := []byte("{\"foo\":\"bar\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/telemetry", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected creation of invalid/unparseable new telemetry item to return bad request, got %d", recorder.Code)
	}
	if dispatcher.DispatchCount != 0 {
		t.Errorf("Expected dispatcher to dispatch 0 messages, got %d", dispatcher.DispatchCount)
	}
}

func TestAddValidPositionCreatesCommand(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	dispatcher := fakes.NewFakeQueueDispatcher()
	server := MakeTestServer(dispatcher)
	recorder = httptest.NewRecorder()
	body := []byte("{\"foo\":\"bar\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/positions", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new position item to return 201, got %d", recorder.Code)
	}
}

func TestAddValidAlertCreatesCommand(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	dispatcher := fakes.NewFakeQueueDispatcher()
	server := MakeTestServer(dispatcher)
	recorder = httptest.NewRecorder()

	body := []byte("{\"drone_id\":\"alertingdrone4\", \"fault_code\": 12, \"description\": \"all the things are failing\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/alerts", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new alert item to return 201, got %d/%s", recorder.Code, string(recorder.Body.Bytes()))
	}
	if dispatcher.DispatchCount != 1 {
		t.Errorf("Expected queue dispatch count of 1, got %d", dispatcher.DispatchCount)
	}
	if len(dispatcher.Messages["alert"]) != 1 {
		t.Errorf("Expected alert message count of 1, got %d", len(dispatcher.Messages["alert"]))
	}

	var alertResponse dronescommon.AlertSignalledEvent
	payload := recorder.Body.Bytes()
	err := json.Unmarshal(payload, &alertResponse)
	if err != nil {
		t.Errorf("Could not unmarshal payload into alertResponse object")
	}
	if alertResponse.DroneID != "alertingdrone4" {
		t.Errorf("Expected drone ID of 'alertingdrone4' got %s", alertResponse.DroneID)
	}
	if alertResponse.FaultCode != 12 {
		t.Errorf("Expected drone fault code of 12, got %d", alertResponse.FaultCode)
	}
}
