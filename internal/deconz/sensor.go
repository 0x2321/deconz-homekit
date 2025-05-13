// Package deconz provides interfaces and types for interacting with the deCONZ REST API.
package deconz

import "deconz-homekit/internal/client"

// Sensor represents a sensor device in the deCONZ ecosystem.
// This struct contains all the properties and state information for a sensor,
// including its configuration, identification, and current readings.
// Sensors can be of various types including motion sensors, temperature sensors,
// open/close sensors, water leak sensors, etc.
type Sensor struct {
	// Config contains the configuration parameters for this sensor
	// This may include settings like sensitivity, reporting intervals, etc.
	Config ObjectMap `json:"config"`

	// Endpoint is the Zigbee endpoint number for this sensor
	Endpoint int `json:"endpoint"`

	// ETag is used for caching and resource versioning
	ETag string `json:"etag"`

	// LastSeen is the timestamp when the sensor was last seen by the gateway
	LastSeen string `json:"lastseen"`

	// Manufacturer is the name of the device manufacturer
	Manufacturer string `json:"manufacturername"`

	// ModelId is the model identifier of the sensor
	ModelId string `json:"modelid"`

	// Name is the user-assigned name of the sensor
	Name string `json:"name"`

	// State contains the current state/readings of the sensor
	// The contents vary depending on the sensor type (e.g., presence, temperature, etc.)
	State ObjectMap `json:"state"`

	// SwVersion is the firmware version running on the sensor
	SwVersion string `json:"swversion"`

	// Type is the type of the sensor (e.g., "ZHAPresence", "ZHATemperature", "ZHAOpenClose")
	Type string `json:"type"`

	// UniqueId is the unique identifier for this sensor
	UniqueId string `json:"uniqueid"`
}

// GetSensor retrieves detailed information about a specific sensor from the deCONZ gateway.
//
// Parameters:
//   - id: The identifier of the sensor to retrieve
//
// Returns:
//   - *Sensor: A pointer to the retrieved Sensor structure
//   - error: Any error encountered during the API request
func (ac *ApiClient) GetSensor(id string) (*Sensor, error) {
	return client.Get[Sensor](ac.buildUrl("/sensors/" + id))
}
