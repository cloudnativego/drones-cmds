package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	dronescommon "github.com/cloudnativego/drones-common"
	"github.com/unrolled/render"
)

func addTelemetryHandler(formatter *render.Render, dispatcher queueDispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newTelemetryCommand telemetryCommand
		err := json.Unmarshal(payload, &newTelemetryCommand)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse add telemetry command.")
			return
		}
		if !newTelemetryCommand.isValid() {
			formatter.Text(w, http.StatusBadRequest, "Invalid telemetry command.")
			return
		}

		evt := dronescommon.TelemetryUpdatedEvent{
			DroneID:          newTelemetryCommand.DroneID,
			RemainingBattery: newTelemetryCommand.RemainingBattery,
			Uptime:           newTelemetryCommand.Uptime,
			CoreTemp:         newTelemetryCommand.CoreTemp,
			ReceivedOn:       time.Now().UnixNano(),
		}
		dispatcher.DispatchMessage("telemetry", evt)
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
