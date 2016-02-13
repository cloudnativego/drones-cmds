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

	dispatcher := fakes.NewFakeQueueDispatcher()
	// TODO - if we detect real bound Rabbit, use the AMQP dispatcher
	//dispatcher := NewAMQPDispatcher(nil)
	initRoutes(mx, formatter, dispatcher)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, dispatcher queueDispatcher) {
	mx.HandleFunc("/api/cmds/telemetry", addTelemetryHandler(formatter, dispatcher)).Methods("POST")
	mx.HandleFunc("/api/cmds/alerts", addAlertHandler(formatter, dispatcher)).Methods("POST")
	mx.HandleFunc("/api/cmds/positions", addPositionHandler(formatter, dispatcher)).Methods("POST")
}
