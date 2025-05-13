// Package accessoryManager provides functionality for creating and managing HomeKit accessories
// that represent deCONZ devices.
package accessoryManager

import (
	"deconz-homekit/internal/deconz"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
)

// OpenCloseSensor represents a contact sensor in HomeKit.
// It implements the DeviceService interface and provides functionality for
// monitoring the open/closed state of doors, windows, and other contact sensors.
type OpenCloseSensor struct {
	// service is the HomeKit contact sensor service
	service *service.ContactSensor

	// device is a reference to the parent Device
	device *Device

	// lowBatteryCharacteristic is the HomeKit characteristic for low battery status
	// This is optional and only present if the sensor reports battery status
	lowBatteryCharacteristic   *characteristic.StatusLowBattery
	batteryLevelCharacteristic *characteristic.BatteryLevel
}

// S returns the underlying HomeKit service.
// This method implements the DeviceService interface.
//
// Returns:
//   - *service.S: A pointer to the HomeKit service
func (sensor *OpenCloseSensor) S() *service.S {
	return sensor.service.S
}

// UpdateState updates the sensor's state based on updates from the deCONZ gateway.
// This method implements the DeviceService interface.
//
// Parameters:
//   - state: The updated state object from deCONZ
//   - config: The updated config object from deCONZ (not used for open/close sensors)
func (sensor *OpenCloseSensor) UpdateState(state deconz.MapObject) {
	// Update the contact sensor state based on the "open" value from deCONZ
	// In HomeKit, 1 = detected (open), 0 = not detected (closed)
	if state.ValueToBool("open") {
		sensor.device.log.Info("open")
		_ = sensor.service.ContactSensorState.SetValue(1) // Contact detected (open)
	} else {
		sensor.device.log.Info("closed")
		_ = sensor.service.ContactSensorState.SetValue(0) // Contact not detected (closed)
	}

	// Update the low battery characteristic if available
	if state.Has("lowbattery") && sensor.lowBatteryCharacteristic != nil {
		batteryIsLow := state.ValueToBool("lowbattery")
		// Convert boolean to int (0 = normal, 1 = low)
		_ = sensor.lowBatteryCharacteristic.SetValue(boolToInt[batteryIsLow])
	}
}

// UpdateConfig updates the sensor's configuration based on updates from the deCONZ gateway.
// This method implements the DeviceService interface.
//
// Parameters:
//   - config: The updated configuration object from deCONZ
func (sensor *OpenCloseSensor) UpdateConfig(config deconz.MapObject) {
	// Update the battery level characteristic if available
	if config.Has("battery") && sensor.batteryLevelCharacteristic != nil {
		batteryLevel := config.ValueToInt("battery")
		_ = sensor.batteryLevelCharacteristic.SetValue(batteryLevel)
	}
}

// NewOpenCloseSensor creates a new open/close sensor service.
// This is used for door/window contact sensors that report open/closed states.
//
// Parameters:
//   - config: A pointer to the deCONZ subdevice configuration
//
// Returns:
//   - error: An error if the service could not be created
func (device *Device) NewOpenCloseSensor(config *deconz.Subdevice) error {
	sensor := new(OpenCloseSensor)
	sensor.device = device

	// Create a new HomeKit contact sensor service
	sensor.service = service.NewContactSensor()

	// Add the low battery characteristic if the sensor reports battery status
	if config.State.Has("lowbattery") {
		sensor.lowBatteryCharacteristic = characteristic.NewStatusLowBattery()
		sensor.service.AddC(sensor.lowBatteryCharacteristic.C)
	}

	// Add the battery level characteristic if the sensor reports battery config
	if config.Config.Has("battery") {
		sensor.batteryLevelCharacteristic = characteristic.NewBatteryLevel()
		sensor.service.AddC(sensor.batteryLevelCharacteristic.C)
	}

	// Initialize the sensor state from the current deCONZ state
	sensor.UpdateState(config.State)
	sensor.UpdateConfig(config.Config)

	// Register the service with the device
	device.addDeviceService(config.UniqueId, sensor)
	return nil
}
