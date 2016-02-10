package main

import (
	"net/http"
	"time"

	dronescommon "github.com/cloudnativego/drones-common"
	"github.com/unrolled/render"
)

func addTelemetryHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		evt := dronescommon.TelemetryUpdatedEvent{
			DroneID:          "foo",
			RemainingBattery: 1,
			Uptime:           1,
			CoreTemp:         1,
			ReceivedOn:       time.Now().UnixNano(),
		}
		formatter.JSON(w, http.StatusCreated, evt)
	}
}

func addAlertHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusCreated, "tbd")
	}
}

func addPositionHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusCreated, "tbd")
	}
}
