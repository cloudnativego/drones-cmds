package main

import (
	"github.com/cloudnativego/drones-cmds/fakes"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	positionDispatcher := fakes.NewFakeQueueDispatcher()
	telemetryDispatcher := fakes.NewFakeQueueDispatcher()
	alertDispatcher := fakes.NewFakeQueueDispatcher()
	// TODO - if we detect real bound Rabbit, use the AMQP dispatcher

	initRoutes(mx, formatter, telemetryDispatcher, alertDispatcher, positionDispatcher)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, telemetryDispatcher queueDispatcher, alertDispatcher queueDispatcher, positionDispatcher queueDispatcher) {
	mx.HandleFunc("/api/cmds/telemetry", addTelemetryHandler(formatter, telemetryDispatcher)).Methods("POST")
	mx.HandleFunc("/api/cmds/alerts", addAlertHandler(formatter, alertDispatcher)).Methods("POST")
	mx.HandleFunc("/api/cmds/positions", addPositionHandler(formatter, positionDispatcher)).Methods("POST")
}
