// Package accessoryManager provides functionality for creating and managing HomeKit accessories
// that represent deCONZ devices. It handles the conversion between deCONZ device types and
// HomeKit accessory types, as well as processing real-time updates from the deCONZ gateway.
package accessoryManager

import (
	"deconz-homekit/internal/deconz"
	"github.com/brutella/hap/accessory"
	"maps"
	"slices"
)

// AccessoryManager manages all HomeKit accessories and their services.
// It maintains mappings between deCONZ devices and HomeKit accessories,
// and handles real-time updates from the deCONZ gateway.
type AccessoryManager struct {
	// Devices is a map of deCONZ device unique IDs to Device objects
	Devices map[string]*Device

	// Services is a map of deCONZ device unique IDs to DeviceService interfaces
	// This provides quick access to services for processing updates
	Services map[string]DeviceService
}

// NewAccessoryManager creates a new AccessoryManager and initializes it with devices
// from the deCONZ gateway.
//
// Parameters:
//   - client: A pointer to the deCONZ API client for communication with the gateway
//   - devices: A slice of deCONZ devices to be converted to HomeKit accessories
//
// Returns:
//   - *AccessoryManager: A pointer to the initialized AccessoryManager
func NewAccessoryManager(client *deconz.ApiClient, devices []*deconz.Device) *AccessoryManager {
	am := new(AccessoryManager)
	am.Devices = make(map[string]*Device)
	am.Services = make(map[string]DeviceService)

	// Create HomeKit devices for each deCONZ device
	for _, config := range devices {
		device, err := NewDevice(client, config)
		if err != nil {
			// Skip devices that cannot be converted to HomeKit accessories
			continue
		}
		am.Devices[config.UniqueId] = device
	}

	// Collect all services from all devices for quick lookup during updates
	for _, device := range am.Devices {
		maps.Copy(am.Services, device.Services)
	}

	return am
}

// GetAccessories returns all HomeKit accessories managed by this AccessoryManager.
// This is used when setting up the HomeKit server.
//
// Returns:
//   - []*accessory.A: A slice of pointers to HomeKit accessories
func (am *AccessoryManager) GetAccessories() []*accessory.A {
	accessories := []*accessory.A{}

	// Collect all accessories from all devices
	for _, device := range am.Devices {
		accessories = append(accessories, device.Accessory)
	}

	return accessories
}

// ProcessUpdate processes a real-time update message from the deCONZ gateway.
// It updates the state of the corresponding HomeKit accessory service.
//
// Parameters:
//   - msg: A pointer to the message containing the update information
func (am *AccessoryManager) ProcessUpdate(msg *deconz.Messsage) {
	// Only process updates for lights and sensors
	if !slices.Contains([]deconz.RessourceType{deconz.LightsRessource, deconz.SensorsRessource}, msg.RessourceType) {
		// Ignore messages for other resource types
		return
	}

	// Only process state change events
	if msg.EventType != deconz.ChangedEvent {
		// For other event types (added, deleted, scene-called), a restart would be needed
		// to properly handle the changes in the device configuration
		return
	}

	// Find the service corresponding to the device and update its state
	id := *msg.UniqueID
	if service := am.Services[id]; service != nil {
		if msg.State != nil {
			service.UpdateState(msg.State)
		}
		if msg.Config != nil {
			service.UpdateConfig(msg.Config)
		}
	}
}
