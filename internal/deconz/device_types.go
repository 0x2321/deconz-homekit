// Package deconz provides interfaces and types for interacting with the deCONZ REST API.
// It handles communication with the deCONZ gateway to control and monitor Zigbee devices
// such as lights, sensors, and switches through the HomeKit bridge.
package deconz

// Reference to the deCONZ REST API constants:
// https://github.com/dresden-elektronik/deconz-rest-plugin/blob/stable/devices/generic/constants.json

// DeviceType represents a ZHA (Zigbee Home Automation) device type.
// These types are used to identify and categorize different Zigbee devices
// in the deCONZ ecosystem.
type DeviceType string

const (
	// AirPurifierDevice represents a ZHA air purifier device.
	// This device type is used for air purifiers that can be controlled via Zigbee.
	AirPurifierDevice DeviceType = "ZHAAirPurifier"

	// AirQualityDevice represents a ZHA air quality sensor.
	// These sensors monitor air quality metrics such as particulate matter or VOCs.
	AirQualityDevice DeviceType = "ZHAAirQuality"

	// AlarmDevice represents a ZHA alarm sensor.
	// These devices detect and report alarm conditions.
	AlarmDevice DeviceType = "ZHAAlarm"

	// AncillaryControlDevice represents a ZHA ancillary control device.
	// These are supplementary control devices in a Zigbee network.
	AncillaryControlDevice DeviceType = "ZHAAncillaryControl"

	// BatteryDevice represents a ZHA battery sensor.
	// These sensors report battery levels of Zigbee devices.
	BatteryDevice DeviceType = "ZHABattery"

	// CarbonDioxideDevice represents a ZHA carbon dioxide sensor.
	// These sensors detect and report CO2 levels in the environment.
	CarbonDioxideDevice DeviceType = "ZHACarbonDioxide"

	// CarbonMonoxideDevice represents a ZHA carbon monoxide sensor.
	// These sensors detect and report dangerous CO levels for safety purposes.
	CarbonMonoxideDevice DeviceType = "ZHACarbonMonoxide"

	// ColorLightDevice represents a ZHA color light.
	// These lights support color control (typically RGB).
	ColorLightDevice DeviceType = "Color light"

	// ColorTemperatureLightDevice represents a ZHA color temperature light.
	// These lights support adjustable white color temperature (warm to cool).
	ColorTemperatureLightDevice DeviceType = "Color temperature light"

	// ConsumptionDevice represents a ZHA consumption sensor.
	// These sensors measure and report energy or resource consumption.
	ConsumptionDevice DeviceType = "ZHAConsumption"

	// DimmableLightDevice represents a ZHA dimmable light.
	// These lights support brightness adjustment but not color control.
	DimmableLightDevice DeviceType = "Dimmable light"

	// DimmablePlugInUnitDevice represents a ZHA dimmable plug-in unit.
	// These are plug-in modules that can control brightness of connected lights.
	DimmablePlugInUnitDevice DeviceType = "Dimmable plug-in unit"

	// DimmerSwitchDevice represents a ZHA dimmer switch.
	// These are wall switches or remotes that control light brightness.
	DimmerSwitchDevice DeviceType = "Dimmer switch"

	// DoorLockDevice represents a ZHA door lock sensor.
	// These devices detect and report the locked/unlocked state of doors.
	DoorLockDevice DeviceType = "ZHADoorLock"

	// DoorLockControllerDevice represents a ZHA door lock controller.
	// These devices can control door locks remotely.
	DoorLockControllerDevice DeviceType = "Door lock controller"

	// DoorLockSensorDevice represents a ZHA door lock sensor.
	// Alternative type for devices that detect door lock states.
	DoorLockSensorDevice DeviceType = "Door Lock"

	// ExtendedColorLightDevice represents a ZHA extended color light.
	// These lights support advanced color control features beyond basic RGB.
	ExtendedColorLightDevice DeviceType = "Extended color light"

	// FireSensorDevice represents a ZHA fire sensor.
	// These sensors detect and report fire or smoke conditions.
	FireSensorDevice DeviceType = "ZHAFire"

	// HumiditySensorDevice represents a ZHA humidity sensor.
	// These sensors measure and report relative humidity levels.
	HumiditySensorDevice DeviceType = "ZHAHumidity"

	// LevelControlSwitchDevice represents a ZHA level control switch.
	// These switches can control variable levels (like dimming).
	LevelControlSwitchDevice DeviceType = "Level control switch"

	// LightLevelSensorDevice represents a ZHA light level sensor.
	// These sensors measure and report ambient light levels.
	LightLevelSensorDevice DeviceType = "ZHALightLevel"

	// MoistureSensorDevice represents a ZHA moisture sensor.
	// These sensors detect and report moisture levels.
	MoistureSensorDevice DeviceType = "ZHAMoisture"

	// OnOffLightDevice represents a ZHA on/off light.
	// These are basic lights that can only be turned on or off.
	OnOffLightDevice DeviceType = "On/Off light"

	// OnOffLightSwitchDevice represents a ZHA on/off light switch.
	// These switches control on/off state of connected lights.
	OnOffLightSwitchDevice DeviceType = "On/Off light switch"

	// OnOffOutputDevice represents a ZHA on/off output.
	// These are general purpose on/off control devices.
	OnOffOutputDevice DeviceType = "On/Off output"

	// OnOffPlugInUnitDevice represents a ZHA on/off plug-in unit.
	// These are plug-in modules that can turn connected devices on/off.
	OnOffPlugInUnitDevice DeviceType = "On/Off plug-in unit"

	// OnOffSwitchDevice represents a ZHA on/off switch.
	// These are general purpose on/off switches.
	OnOffSwitchDevice DeviceType = "On/Off switch"

	// OpenCloseSensorDevice represents a ZHA open/close sensor.
	// These sensors detect and report if doors/windows are open or closed.
	OpenCloseSensorDevice DeviceType = "ZHAOpenClose"

	// ParticulateMatterDevice represents a ZHA particulate matter sensor.
	// These sensors measure and report particulate matter in the air.
	ParticulateMatterDevice DeviceType = "ZHAParticulateMatter"

	// PresenceSensorDevice represents a ZHA presence sensor.
	// These sensors detect and report motion or presence in an area.
	PresenceSensorDevice DeviceType = "ZHAPresence"

	// PressureDevice represents a ZHA pressure sensor.
	// These sensors measure and report atmospheric pressure.
	PressureDevice DeviceType = "ZHAPressure"

	// RangeExtenderDevice represents a ZHA range extender.
	// These devices extend the range of the Zigbee network.
	RangeExtenderDevice DeviceType = "Range extender"

	// RelativeRotaryDevice represents a ZHA relative rotary control.
	// These are rotary controllers that report relative movement.
	RelativeRotaryDevice DeviceType = "ZHARelativeRotary"

	// SmartPlugDevice represents a ZHA smart plug.
	// These are intelligent power outlets that can be controlled remotely.
	SmartPlugDevice DeviceType = "Smart plug"

	// SpectralDevice represents a ZHA spectral sensor.
	// These sensors measure and report spectral characteristics.
	SpectralDevice DeviceType = "ZHASpectral"

	// SwitchDevice represents a ZHA switch.
	// These are general purpose switches or buttons.
	SwitchDevice DeviceType = "ZHASwitch"

	// TemperatureDevice represents a ZHA temperature sensor.
	// These sensors measure and report ambient temperature.
	TemperatureDevice DeviceType = "ZHATemperature"

	// ThermostatDevice represents a ZHA thermostat.
	// These devices control heating and cooling systems.
	ThermostatDevice DeviceType = "ZHAThermostat"

	// TimeDevice represents a ZHA time device.
	// These devices provide time-related functionality.
	TimeDevice DeviceType = "ZHATime"

	// VibrationDevice represents a ZHA vibration sensor.
	// These sensors detect and report vibration or movement.
	VibrationDevice DeviceType = "ZHAVibration"

	// WarningDevice represents a ZHA warning device.
	// These devices provide alerts or warnings (like sirens).
	WarningDevice DeviceType = "Warning device"

	// WaterDevice represents a ZHA water leak sensor.
	// These sensors detect and report water leaks.
	WaterDevice DeviceType = "ZHAWater"

	// WindowCoveringDevice represents a ZHA window covering device.
	// These devices control blinds, shades, or curtains.
	WindowCoveringDevice DeviceType = "Window covering device"
)
