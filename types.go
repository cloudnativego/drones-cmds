package main

type telemetryCommand struct {
	DroneID          string `json:"drone_id"`
	RemainingBattery int    `json:"battery"`
	Uptime           int    `json:"uptime"`
	CoreTemp         int    `json:"core_temp"`
}

type alertCommand struct {
	DroneID     string `json:"drone_id"`
	FaultCode   int    `json:"fault_code"`
	Description string `json:"description"`
}

type positionCommand struct {
	DroneID         string  `json:"drone_id"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	Altitude        float32 `json:"altitude"`
	CurrentSpeed    float32 `json:"current_speed"`
	HeadingCardinal int     `json:"heading_cardinal"`
}

type queueDispatcher interface {
	DispatchMessage(message interface{}) (err error)
}

func (telemetry telemetryCommand) isValid() (valid bool) {
	valid = true
	if len(telemetry.DroneID) == 0 {
		valid = false
	}
	if telemetry.Uptime == 0 {
		valid = false
	}
	return valid
}

func (alert alertCommand) isValid() (valid bool) {
	valid = true
	if len(alert.DroneID) == 0 {
		valid = false
	}
	if len(alert.Description) == 0 {
		valid = false
	}
	return
}

func (position positionCommand) isValid() (valid bool) {
	valid = true
	if len(position.DroneID) == 0 {
		valid = false
	}
	if position.Longitude < 0 || position.Latitude < 0 || position.HeadingCardinal < 0 {
		valid = false
	}
	// There's little other validation we can add to a position report without
	// assigning 0 as a "bad data" signal, which we know is a terrible idea for
	// any IoT data gathering app.
	return
}
