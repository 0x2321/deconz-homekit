// Package accessoryManager provides functionality for creating and managing HomeKit accessories
// that represent deCONZ devices.
package accessoryManager

import (
	"deconz-homekit/internal/deconz"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
)

// WaterSensor represents a water leak sensor in HomeKit.
// It implements the DeviceService interface and provides functionality for
// monitoring water leaks from compatible sensors.
type WaterSensor struct {
	// device is a reference to the parent Device
	device *Device

	// service is the HomeKit leak sensor service
	service *service.LeakSensor

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
func (sensor *WaterSensor) S() *service.S {
	return sensor.service.S
}

// UpdateState updates the sensor's state based on updates from the deCONZ gateway.
// This method implements the DeviceService interface.
//
// Parameters:
//   - state: The updated state object from deCONZ
//   - config: The updated config object from deCONZ (not used for water sensors)
func (sensor *WaterSensor) UpdateState(state deconz.StateObject, config deconz.StateObject) {
	// Update the leak detection state based on the "water" value from deCONZ
	// In HomeKit, 1 = leak detected, 0 = no leak detected
	v := state.ValueToBool("water")
	_ = sensor.service.LeakDetected.SetValue(boolToInt[v])

	// Log when a leak is detected (only log positive detections to reduce noise)
	if v {
		sensor.device.log.Info("leak detected")
	}

	// Update the low battery characteristic if available
	if state.Has("lowbattery") && sensor.lowBatteryCharacteristic != nil {
		batteryIsLow := state.ValueToBool("lowbattery")
		// Convert boolean to int (0 = normal, 1 = low)
		_ = sensor.lowBatteryCharacteristic.SetValue(boolToInt[batteryIsLow])
	}

	// Update the battery level characteristic if available
	if config.Has("battery") && sensor.batteryLevelCharacteristic != nil {
		batteryLevel := config.ValueToInt("battery")
		_ = sensor.batteryLevelCharacteristic.SetValue(batteryLevel)
	}
}

// NewWaterSensor creates a new water leak sensor service.
// This is used for sensors that detect water leaks.
//
// Parameters:
//   - config: A pointer to the deCONZ subdevice configuration
//
// Returns:
//   - error: An error if the service could not be created
func (device *Device) NewWaterSensor(config *deconz.Subdevice) error {
	sensor := new(WaterSensor)
	sensor.device = device

	// Create a new HomeKit leak sensor service
	sensor.service = service.NewLeakSensor()

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
	sensor.UpdateState(config.State, config.Config)

	// Register the service with the device
	device.addDeviceService(config.UniqueId, sensor)
	return nil
}
