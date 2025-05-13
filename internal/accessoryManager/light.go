// Package accessoryManager provides functionality for creating and managing HomeKit accessories
// that represent deCONZ devices.
package accessoryManager

import (
	"deconz-homekit/internal/deconz"
	"deconz-homekit/internal/helper"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
	"time"
)

// Light represents a light device in HomeKit.
// It implements the DeviceService interface and provides functionality for
// controlling lights with various capabilities (on/off, brightness, color temperature).
type Light struct {
	// ID is the unique identifier of the light (from deCONZ)
	ID string

	// On is the HomeKit characteristic for the on/off state
	On *characteristic.On

	// Brightness is the HomeKit characteristic for brightness level
	Brightness *characteristic.Brightness

	// ColorTemperature is the HomeKit characteristic for color temperature
	ColorTemperature *characteristic.ColorTemperature

	// lastChange tracks when the light was last changed by a user command
	// This is used to prevent feedback loops when updating state
	lastChange *time.Time

	// device is a reference to the parent Device
	device *Device

	// service is the HomeKit service for this light
	service *service.S
}

// NewLight creates a new Light service with the specified service type.
// The service type determines whether the light appears as a lightbulb or outlet in HomeKit.
//
// Parameters:
//   - device: A pointer to the parent Device
//   - config: A pointer to the deCONZ subdevice configuration
//   - serviceType: The HomeKit service type (e.g., service.TypeLightbulb, service.TypeOutlet)
//
// Returns:
//   - *Light: A pointer to the initialized Light
func NewLight(device *Device, config *deconz.Subdevice, serviceType string) *Light {
	lightbulb := new(Light)
	lightbulb.ID = config.UniqueId
	lightbulb.device = device

	// Create a new HomeKit service of the specified type
	lightbulb.service = service.New(serviceType)
	device.addDeviceService(config.UniqueId, lightbulb)

	return lightbulb
}

// S returns the underlying HomeKit service.
// This method implements the DeviceService interface.
//
// Returns:
//   - *service.S: A pointer to the HomeKit service
func (light *Light) S() *service.S {
	return light.service
}

// updateChange records the current time as the last change time.
// This is used to ignore state updates from deCONZ for a short period
// after a user-initiated change to prevent feedback loops.
func (light *Light) updateChange() {
	now := time.Now()
	light.lastChange = &now
}

// enableOn adds the On characteristic to the light service.
// This allows the light to be turned on and off through HomeKit.
func (light *Light) enableOn() {
	light.On = characteristic.NewOn()
	// Register the SetOn method to be called when the value is changed through HomeKit
	light.On.OnValueRemoteUpdate(light.SetOn)

	// Add the characteristic to the service
	light.service.AddC(light.On.C)
}

// enableBrightness adds the Brightness characteristic to the light service.
// This allows the light's brightness to be controlled through HomeKit.
func (light *Light) enableBrightness() {
	light.Brightness = characteristic.NewBrightness()
	// Register the SetBrightness method to be called when the value is changed through HomeKit
	light.Brightness.OnValueRemoteUpdate(light.SetBrightness)

	// Add the characteristic to the service
	light.service.AddC(light.Brightness.C)
}

// enableColorTemperature adds the ColorTemperature characteristic to the light service.
// This allows the light's color temperature to be controlled through HomeKit.
func (light *Light) enableColorTemperature() {
	light.ColorTemperature = characteristic.NewColorTemperature()
	// Register the SetColorTemperature method to be called when the value is changed through HomeKit
	light.ColorTemperature.OnValueRemoteUpdate(light.SetColorTemperature)

	// Set the minimum and maximum color temperature values in mireds
	if details, err := light.device.client.GetLight(light.ID); err == nil {
		if ctMin := details.CtMin; ctMin != nil {
			light.ColorTemperature.SetMinValue(*ctMin)
		}
		if ctMax := details.CtMax; ctMax != nil {
			light.ColorTemperature.SetMaxValue(*ctMax)
		}
	}

	// Add the characteristic to the service
	light.service.AddC(light.ColorTemperature.C)
}

// SetOn turns the light on or off.
// This method is called when the On characteristic is changed through HomeKit.
//
// Parameters:
//   - on: A boolean indicating whether to turn the light on (true) or off (false)
func (light *Light) SetOn(on bool) {
	light.device.log.Infof("set %s", onOffStr[on])

	// Send the command to the deCONZ gateway
	if err := light.device.client.SetLightOn(light.ID, on); err != nil {
		light.device.log.Errorf("failed to set light %s: %+v", onOffStr[on], err)
	}
	light.updateChange()
}

// SetBrightness sets the brightness of the light.
// This method is called when the Brightness characteristic is changed through HomeKit.
//
// Parameters:
//   - v: An integer representing the brightness percentage (0-100)
func (light *Light) SetBrightness(v int) {
	light.device.log.Infof("set brightness to %d%%", v)

	// Send the command to the deCONZ gateway
	if err := light.device.client.SetLightBrightness(light.ID, v); err != nil {
		light.device.log.Errorf("failed to set brightness: %+v", err)
	}
	light.updateChange()
}

// SetColorTemperature sets the color temperature of the light.
// This method is called when the ColorTemperature characteristic is changed through HomeKit.
//
// Parameters:
//   - v: An integer representing the color temperature in mireds
func (light *Light) SetColorTemperature(v int) {
	// Convert mireds to Kelvin for logging (mireds = 1,000,000/Kelvin)
	k := 1_000_000.0 / float64(v)
	light.device.log.Infof("set color temperature to %.1f K (%d)", k, v)

	// Send the command to the deCONZ gateway
	if err := light.device.client.SetLightColorTemperature(light.ID, v); err != nil {
		light.device.log.Errorf("failed to set color temperature: %+v", err)
	}
	light.updateChange()
}

// UpdateState updates the light's state based on updates from the deCONZ gateway.
// This method implements the DeviceService interface.
//
// Parameters:
//   - state: The updated state object from deCONZ
//   - _: The updated config object from deCONZ (not used for lights)
func (light *Light) UpdateState(state deconz.MapObject) {
	// Ignore updates for a short period after a user-initiated change
	// to prevent feedback loops
	if light.lastChange != nil {
		ignoreUntil := light.lastChange.Add(time.Second)
		if time.Now().Before(ignoreUntil) {
			return
		}
	}

	// Update the On characteristic if the state contains an "on" value
	if state.Has("on") && light.On != nil {
		light.On.SetValue(state.ValueToBool("on"))
	}

	// Update the Brightness characteristic if the state contains a "bri" value
	if state.Has("bri") && light.Brightness != nil {
		_ = light.Brightness.SetValue(state.ValueToPercent("bri"))
	}

	// Update the ColorTemperature characteristic if the state contains a "ct" value
	if state.Has("ct") && light.ColorTemperature != nil {
		_ = light.ColorTemperature.SetValue(state.ValueToInt("ct"))
	}
}

// UpdateConfig updates the light's configuration based on updates from the deCONZ gateway.
// This method implements the DeviceService interface.
// For lights, this method currently does nothing as lights don't have configuration
// parameters that need to be updated.
//
// Parameters:
//   - config: The updated configuration object from deCONZ (not used for lights)
func (light *Light) UpdateConfig(_ deconz.MapObject) {
	// nothing to do
}

// NewOnOffLight creates a new on/off light service.
// This is used for lights that only support being turned on or off.
//
// Parameters:
//   - config: A pointer to the deCONZ subdevice configuration
//
// Returns:
//   - error: An error if the service could not be created
func (device *Device) NewOnOffLight(config *deconz.Subdevice) error {
	light := NewLight(device, config, service.TypeLightbulb)
	light.enableOn()
	light.UpdateState(config.State)

	return nil
}

// NewDimmableLight creates a new dimmable light service.
// This is used for lights that support being turned on/off and brightness control.
//
// Parameters:
//   - config: A pointer to the deCONZ subdevice configuration
//
// Returns:
//   - error: An error if the service could not be created
func (device *Device) NewDimmableLight(config *deconz.Subdevice) error {
	light := NewLight(device, config, service.TypeLightbulb)
	light.enableOn()
	light.enableBrightness()
	light.UpdateState(config.State)

	return nil
}

// NewColorTemperatureLight creates a new color temperature light service.
// This is used for lights that support being turned on/off, brightness control,
// and color temperature control.
//
// Parameters:
//   - config: A pointer to the deCONZ subdevice configuration
//
// Returns:
//   - error: An error if the service could not be created
func (device *Device) NewColorTemperatureLight(config *deconz.Subdevice) error {
	light := NewLight(device, config, service.TypeLightbulb)
	light.enableOn()
	light.enableBrightness()
	light.enableColorTemperature()
	light.UpdateState(config.State)

	return nil
}

// NewOnOffPlugDevice creates a new on/off plug device service.
// This is used for plug-in units and outlets that can be turned on or off.
//
// Parameters:
//   - config: A pointer to the deCONZ subdevice configuration
//
// Returns:
//   - error: An error if the service could not be created
func (device *Device) NewOnOffPlugDevice(config *deconz.Subdevice) error {
	plug := NewLight(device, config, service.TypeOutlet)
	plug.enableOn()
	plug.UpdateState(config.State)

	return nil
}
