// Package accessoryManager provides functionality for creating and managing HomeKit accessories
// that represent deCONZ devices.
package accessoryManager

import (
	"deconz-homekit/internal/deconz"
	deviceConfiguration "deconz-homekit/internal/device_configuration"
	"fmt"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
	"maps"
	"slices"
	"strconv"
)

// SwitchDevice represents a multi-button switch or remote control in HomeKit.
// It implements the DeviceService interface and provides functionality for
// handling button presses from Zigbee remotes and switches.
// Unlike other device types, a SwitchDevice can contain multiple services,
// one for each button on the physical device.
type SwitchDevice struct {
	// device is a reference to the parent Device
	device *Device

	// services is a map of button IDs to HomeKit stateless programmable switch services
	services map[string]*service.StatelessProgrammableSwitch

	// configs is a map of button IDs to button configurations
	// These configurations define how deCONZ button events map to HomeKit button events
	configs map[string]deviceConfiguration.ButtonConfiguration
}

// S returns the underlying HomeKit service.
// This method implements the DeviceService interface.
// For SwitchDevice, this returns nil because it doesn't have a single service,
// but rather multiple services (one per button) that are added directly to the accessory.
//
// Returns:
//   - *service.S: Always nil for SwitchDevice
func (sensor *SwitchDevice) S() *service.S {
	return nil
}

// UpdateState updates the switch's state based on updates from the deCONZ gateway.
// This method implements the DeviceService interface.
// It processes button events and triggers the appropriate HomeKit event.
//
// Parameters:
//   - state: The updated state object from deCONZ
//   - _: The updated config object from deCONZ (not used for switches)
func (sensor *SwitchDevice) UpdateState(state deconz.StateObject, _ deconz.StateObject) {
	// Process button events from the deCONZ gateway
	if state != nil && state.Has("buttonevent") {
		// Get the button event code from the state
		event := fmt.Sprintf("%d", state.ValueToInt("buttonevent"))

		// Split the event code into device ID (button number) and event ID (press type)
		deviceId, eventId := deviceConfiguration.SplitEventId(event)
		sensor.device.log.Infof("button %s got event %s", deviceId, eventId)

		// Map the deCONZ event to a HomeKit event based on the button configuration
		switch sensor.configs[deviceId].EventMap[event] {
		case deviceConfiguration.ButtonSinglePress:
			_ = sensor.services[deviceId].ProgrammableSwitchEvent.SetValue(characteristic.ProgrammableSwitchEventSinglePress)
		case deviceConfiguration.ButtonDoublePress:
			_ = sensor.services[deviceId].ProgrammableSwitchEvent.SetValue(characteristic.ProgrammableSwitchEventDoublePress)
		case deviceConfiguration.ButtonLongPress:
			_ = sensor.services[deviceId].ProgrammableSwitchEvent.SetValue(characteristic.ProgrammableSwitchEventLongPress)
		}
	}
}

// addButton adds a button service to the switch device.
// Each button on a physical remote control or switch is represented as a separate
// stateless programmable switch service in HomeKit.
//
// Parameters:
//   - config: The button configuration defining the button's behavior
func (sensor *SwitchDevice) addButton(config deviceConfiguration.ButtonConfiguration) {
	// Get a sample event ID to determine the button number
	someEventId := slices.Collect(maps.Keys(config.EventMap))[0]
	buttonNumber, _ := deviceConfiguration.SplitEventId(someEventId)

	// Set the service label index (button number) for the HomeKit service
	buttonIndex, _ := strconv.Atoi(buttonNumber)
	indexCharacteristic := characteristic.NewServiceLabelIndex()
	_ = indexCharacteristic.SetValue(buttonIndex)

	// Determine which button press types (single, double, long) this button supports
	enabledButtonStates := []int{}

	// Helper function to add a button state if it's not already in the list
	appendButtonState := func(id int) {
		if !slices.Contains(enabledButtonStates, id) {
			enabledButtonStates = append(enabledButtonStates, id)
		}
	}

	// Check each event in the button's configuration to see what press types it supports
	for _, event := range slices.Collect(maps.Values(config.EventMap)) {
		switch event {
		case deviceConfiguration.ButtonSinglePress:
			appendButtonState(characteristic.ProgrammableSwitchEventSinglePress)
		case deviceConfiguration.ButtonDoublePress:
			appendButtonState(characteristic.ProgrammableSwitchEventDoublePress)
		case deviceConfiguration.ButtonLongPress:
			appendButtonState(characteristic.ProgrammableSwitchEventLongPress)
		}
	}

	// Create a new HomeKit stateless programmable switch service for this button
	newButton := service.NewStatelessProgrammableSwitch()

	// Set the valid values for the programmable switch event characteristic
	// This tells HomeKit which press types this button supports
	newButton.ProgrammableSwitchEvent.C.ValidVals = enabledButtonStates

	// Add the service label index characteristic to the service
	newButton.AddC(indexCharacteristic.C)

	// Store the button service and configuration
	sensor.services[buttonNumber] = newButton
	sensor.configs[buttonNumber] = config

	// Add the button service directly to the accessory
	sensor.device.Accessory.AddS(newButton.S)
}

// NewSwitch creates a new switch device service.
// This is used for remote controls and wall switches with one or more buttons.
//
// Parameters:
//   - config: A pointer to the deCONZ subdevice configuration
//
// Returns:
//   - error: An error if the service could not be created
func (device *Device) NewSwitch(config *deconz.Subdevice) error {
	sensor := new(SwitchDevice)
	sensor.device = device
	sensor.services = make(map[string]*service.StatelessProgrammableSwitch)
	sensor.configs = make(map[string]deviceConfiguration.ButtonConfiguration)

	// Get detailed information about the sensor from the deCONZ gateway
	sensorInfo, err := device.client.GetSensor(config.UniqueId)
	if err != nil {
		return err
	}

	// Load device configurations from the devices directory
	// These configurations define how different button events map to HomeKit events
	deviceConfigs, err := deviceConfiguration.LoadFromDirectory("./devices")
	if err != nil {
		return fmt.Errorf("error loading device configurations: %v", err)
	}

	// Find the configuration for this specific device model
	deviceConfig, ok := deviceConfigs[sensorInfo.ModelId]
	if !ok {
		return fmt.Errorf("could not find device %s", sensorInfo.ModelId)
	}

	// Add a service for each button defined in the device configuration
	for _, buttonConfig := range deviceConfig.Buttons {
		sensor.addButton(buttonConfig)
	}

	// Initialize the switch state
	sensor.UpdateState(nil, config.Config)

	// Register the service with the device
	device.Services[config.UniqueId] = sensor
	return nil
}
