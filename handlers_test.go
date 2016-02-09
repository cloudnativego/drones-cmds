package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

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
	body := []byte("{\"foo\":\"bar\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/telemetry", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new telemetry item to return 201, got %d", recorder.Code)
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
