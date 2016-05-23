package service

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("About to dispatch message...")
		dispatcher.DispatchMessage(evt)
		formatter.JSON(w, http.StatusCreated, evt)
	}
}

func addAlertHandler(formatter *render.Render, dispatcher queueDispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newAlertCommand alertCommand
		err := json.Unmarshal(payload, &newAlertCommand)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse add alert command.")
			return
		}
		if !newAlertCommand.isValid() {
			formatter.Text(w, http.StatusBadRequest, "Invalid alert command.")
			return
		}
		evt := dronescommon.AlertSignalledEvent{
			DroneID:     newAlertCommand.DroneID,
			FaultCode:   newAlertCommand.FaultCode,
			Description: newAlertCommand.Description,
			ReceivedOn:  time.Now().UnixNano(),
		}
		dispatcher.DispatchMessage(evt)
		formatter.JSON(w, http.StatusCreated, evt)
	}
}

func addPositionHandler(formatter *render.Render, dispatcher queueDispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newPositionCommand positionCommand
		err := json.Unmarshal(payload, &newPositionCommand)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse add position command.")
			return
		}
		if !newPositionCommand.isValid() {
			formatter.Text(w, http.StatusBadRequest, "Invalid position command.")
			return
		}
		evt := dronescommon.PositionChangedEvent{
			DroneID:         newPositionCommand.DroneID,
			Longitude:       newPositionCommand.Longitude,
			Latitude:        newPositionCommand.Latitude,
			Altitude:        newPositionCommand.Altitude,
			CurrentSpeed:    newPositionCommand.CurrentSpeed,
			HeadingCardinal: newPositionCommand.HeadingCardinal,
			ReceivedOn:      time.Now().UnixNano(),
		}
		dispatcher.DispatchMessage(evt)
		formatter.JSON(w, http.StatusCreated, evt)
	}
}
