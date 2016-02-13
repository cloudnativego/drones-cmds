package main

import (
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

	dispatcher := &amqpDispatcher{} // TODO: replace with real values
	initRoutes(mx, formatter, dispatcher)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, dispatcher queueDispatcher) {
	mx.HandleFunc("/api/cmds/telemetry", addTelemetryHandler(formatter, dispatcher)).Methods("POST")
	mx.HandleFunc("/api/cmds/alerts", addAlertHandler(formatter)).Methods("POST")
	mx.HandleFunc("/api/cmds/positions", addPositionHandler(formatter)).Methods("POST")
}
