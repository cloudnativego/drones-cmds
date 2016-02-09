package main

import (
	"net/http"

	"github.com/unrolled/render"
)

func addTelemetryHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}

func addAlertHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}

func addPositionHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}
