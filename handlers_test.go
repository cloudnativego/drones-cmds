package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func MakeTestServer() *negroni.Negroni {
	server := negroni.New() // don't need all the middleware here or logging.
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	server.UseHandler(mx)
	return server
}

func TestAddValidTelemetryCreatesCommand(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	server := MakeTestServer()
	recorder = httptest.NewRecorder()
	body := []byte("{\"drone_id\":\"drone666\", \"battery\": 72, \"uptime\": 6941, \"core_temp\": 21 }")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/telemetry", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new telemetry item to return 201, got %d", recorder.Code)
	}

	var telemetryResponse dronescommon.TelemetryUpdatedEvent
	payload := recorder.Body.Bytes()
	err := json.Unmarshal(payload, &telemetryResponse)
	if err != nil {
		t.Errorf("Could not unmarshal payload into newMatchResponse object")
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

	server := MakeTestServer()
	recorder = httptest.NewRecorder()
	body := []byte("{\"foo\":\"bar\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/telemetry", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected creation of invalid/unparseable new telemetry item to return bad request, got %d", recorder.Code)
	}
}

func TestAddValidPositionCreatesCommand(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	server := MakeTestServer()
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

	server := MakeTestServer()
	recorder = httptest.NewRecorder()
	body := []byte("{\"foo\":\"bar\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/alerts", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new alert item to return 201, got %d", recorder.Code)
	}
}
