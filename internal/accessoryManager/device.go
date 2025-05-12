// Package accessoryManager provides functionality for creating and managing HomeKit accessories
// that represent deCONZ devices.
package accessoryManager

import (
	"deconz-homekit/internal/deconz"
	"errors"
	"fmt"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/service"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

// DeviceService is an interface that all HomeKit service implementations must satisfy.
// It provides methods for updating the service state based on deCONZ state updates
// and accessing the underlying HomeKit service.
type DeviceService interface {
	// UpdateState updates the service state based on deCONZ state and config updates
	UpdateState(state deconz.StateObject, config deconz.StateObject)

	// S returns the underlying HomeKit service
	S() *service.S
}

// Device represents a physical device in HomeKit, which may contain multiple services.
// It maps a deCONZ device to a HomeKit accessory and manages its services.
type Device struct {
	// ID is the unique identifier of the device (from deCONZ)
	ID string

	// Accessory is the HomeKit accessory representing this device
	Accessory *accessory.A

	// Services is a map of deCONZ device unique IDs to DeviceService interfaces
	Services map[string]DeviceService

	// client is the deCONZ API client for communicating with the gateway
	client *deconz.ApiClient

	// log is the logger for this device
	log *log.Logger
}

// NewDevice creates a new Device from a deCONZ device configuration.
// It initializes the HomeKit accessory and adds services for each subdevice.
//
// Parameters:
//   - client: A pointer to the deCONZ API client for communication with the gateway
//   - config: A pointer to the deCONZ device configuration
//
// Returns:
//   - *Device: A pointer to the initialized Device
//   - error: An error if the device could not be created or has no services
func NewDevice(client *deconz.ApiClient, config *deconz.Device) (*Device, error) {
	d := new(Device)
	d.client = client
	d.ID = config.UniqueId
	d.Services = make(map[string]DeviceService)

	// Create a new HomeKit accessory with information from the deCONZ device
	d.Accessory = accessory.New(accessory.Info{
		Name:         config.Name,
		Manufacturer: config.Manufacturer,
		Model:        config.Model,
		Firmware:     config.SwVersion,
		SerialNumber: config.UniqueId,
	}, accessory.TypeUnknown)

	// Convert the deCONZ unique ID to a HomeKit ID format
	d.Accessory.Id = uniqueIdToHomeKitId(config.UniqueId)

	// Initialize a logger for this device
	d.log = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		Prefix:          config.Name,
	})

	// Log device discovery and process each subdevice
	d.log.Infof("discovered device (%s)", config.UniqueId)
	for _, sub := range config.Subdevices {
		if err := addSubdevice(d, &sub); err != nil {
			d.log.Warnf("failed to add the service %s: %+v", sub.Type, err)
		}
	}

	// Ensure the device has at least one service
	if len(d.Services) == 0 {
		d.log.Warn("the device has no active services and will not be added to HomeKit")
		return nil, errors.New("no services found")
	}

	return d, nil
}

// addSubdevice adds a service to a device based on the subdevice type.
// It maps deCONZ device types to HomeKit service types and creates the appropriate service.
//
// Parameters:
//   - dev: A pointer to the Device to add the service to
//   - config: A pointer to the deCONZ subdevice configuration
//
// Returns:
//   - error: An error if the service could not be created or the device type is not supported
func addSubdevice(dev *Device, config *deconz.Subdevice) error {
	// Create the appropriate service based on the device type
	switch config.Type {
	case deconz.DimmableLightDevice:
		return dev.NewDimmableLight(config)
	case deconz.ColorTemperatureLightDevice:
		return dev.NewColorTemperatureLight(config)
	case deconz.PresenceSensorDevice:
		return dev.NewPresenceSensor(config)
	case deconz.OpenCloseSensorDevice:
		return dev.NewOpenCloseSensor(config)
	case deconz.OnOffOutputDevice:
		return dev.NewOnOffPlugDevice(config)
	case deconz.OnOffPlugInUnitDevice:
		return dev.NewOnOffPlugDevice(config)
	case deconz.SmartPlugDevice:
		return dev.NewOnOffPlugDevice(config)
	case deconz.OnOffSwitchDevice:
		return dev.NewOnOffPlugDevice(config)
	case deconz.OnOffLightDevice:
		return dev.NewOnOffLight(config)
	case deconz.OnOffLightSwitchDevice:
		return dev.NewOnOffLight(config)
	case deconz.SwitchDevice:
		return dev.NewSwitch(config)
	case deconz.WaterDevice:
		return dev.NewWaterSensor(config)
	case deconz.DimmablePlugInUnitDevice:
		return dev.NewDimmableLight(config)

	default:
		return fmt.Errorf("not implemented")
	}
}

// addDeviceService adds a service to a device and registers it with the HomeKit accessory.
//
// Parameters:
//   - id: The unique identifier of the service
//   - s: The DeviceService to add
func (device *Device) addDeviceService(id string, s DeviceService) {
	device.Services[id] = s
	device.Accessory.AddS(s.S())
}
