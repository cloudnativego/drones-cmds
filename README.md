[![wercker status](https://app.wercker.com/status/0e2174b4304cf85753b912cc4ca0aafe/m/master "wercker status")](https://app.wercker.com/project/bykey/0e2174b4304cf85753b912cc4ca0aafe)

# Drone Army Sample - Command Processing Service
This is part of the Drone Army sample suite, this service is responsible for processing incoming commands
and converting them into events for dispatch into queues.

## RESTful Endpoints
The following is a list of the REST endpoints exposed by this service.

* **POST** to */api/cmds/telemetry* - Submit a new telemetry update command, adds a telemetry changed event to a queue.
* **POST** to */api/cmds/alerts* - Submit a new alert command, adds an alert signaled event to a queue.
* **POST** to */api/cmds/position* - Submit a new position command, adds a position changed event to a queue.
