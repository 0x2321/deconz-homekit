// Package deconz provides interfaces and types for interacting with the deCONZ REST API.
package deconz

import (
	"deconz-homekit/internal/client"
	"fmt"
)

// Value represents a device state or configuration value with its last update timestamp.
// This structure is used to track when a particular value was last changed.
type Value struct {
	// LastUpdated is the ISO 8601 timestamp when the value was last updated
	LastUpdated string `json:"lastupdated"`

	// Value is the actual value, which can be of various types (boolean, number, string, etc.)
	Value interface{} `json:"value"`
}

// Subdevice represents a logical component of a physical Zigbee device.
// Many Zigbee devices contain multiple functional components (e.g., a switch might have
// buttons, battery status, etc.), each represented as a subdevice.
type Subdevice struct {
	// Type identifies the functional category of this subdevice (e.g., light, sensor)
	Type DeviceType `json:"type"`

	// UniqueId is the unique identifier for this subdevice
	UniqueId string `json:"uniqueid"`

	// Config contains configuration parameters for this subdevice
	Config ExtendedObjectMap `json:"config"`

	// State contains the current state values for this subdevice
	State ExtendedObjectMap `json:"state"`
}

// Device represents a physical Zigbee device in the deCONZ ecosystem.
// A device can contain multiple subdevices representing different functional aspects.
type Device struct {
	// UniqueId is the unique identifier for this device
	UniqueId string `json:"uniqueid"`

	// Manufacturer is the name of the device manufacturer
	Manufacturer string `json:"manufacturername"`

	// Model is the model identifier of the device
	Model string `json:"modelid"`

	// Name is the user-assigned name of the device
	Name string `json:"name"`

	// Product is the product identifier of the device
	Product string `json:"productid"`

	// SwVersion is the firmware version running on the device
	SwVersion string `json:"swversion"`

	// Subdevices is a list of functional components within this device
	Subdevices []Subdevice `json:"subdevices"`
}

// ListDevices retrieves a list of all device unique identifiers from the deCONZ gateway.
//
// Returns:
//   - *[]string: A pointer to a slice of device unique identifiers
//   - error: Any error encountered during the API request
func (ac *ApiClient) ListDevices() (*[]string, error) {
	return client.Get[[]string](ac.buildUrl("/devices"))
}

// GetDevice retrieves detailed information about a specific device from the deCONZ gateway.
//
// Parameters:
//   - uniqueId: The unique identifier of the device to retrieve
//
// Returns:
//   - *Device: A pointer to the retrieved Device structure
//   - error: Any error encountered during the API request
func (ac *ApiClient) GetDevice(uniqueId string) (*Device, error) {
	return client.Get[Device](ac.buildUrl("/devices/" + uniqueId))
}

// GetAllDevices retrieves detailed information about all devices from the deCONZ gateway.
// This method first gets a list of all device IDs, then queries each device individually.
//
// Returns:
//   - []*Device: A slice of pointers to Device structures
//   - error: Any error encountered during the API requests
func (ac *ApiClient) GetAllDevices() ([]*Device, error) {
	allDevices := []*Device{}

	// Get list of all device IDs from the gateway
	devicesList, err := ac.ListDevices()
	if err != nil {
		return nil, err
	}

	// Query each device individually to get detailed information
	for _, deviceId := range *devicesList {
		device, err := ac.GetDevice(deviceId)
		if err != nil {
			// Log the error but continue with other devices
			fmt.Println(err)
			continue
		}
		allDevices = append(allDevices, device)
	}

	return allDevices, nil
}
