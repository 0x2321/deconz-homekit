// Package deconz provides interfaces and types for interacting with the deCONZ REST API.
package deconz

import (
	"deconz-homekit/internal/client"
	"math"
)

// Light represents a light device in the deCONZ ecosystem.
// This struct contains all the properties and state information for a light,
// including its capabilities, identification, and current settings.
type Light struct {
	// ColorCapabilities indicates the color features supported by the light
	// Bit 0: Enhanced Hue, Bit 1: Enhanced Saturation, Bit 2: XY, Bit 3: Color Temperature
	ColorCapabilities *int `json:"colorcapabilities,omitempty"`

	// CtMax is the maximum color temperature in mireds (higher = warmer)
	CtMax *int `json:"ctmax,omitempty"`

	// CtMin is the minimum color temperature in mireds (lower = cooler)
	CtMin *int `json:"ctmin,omitempty"`

	// LastAnnounced is the timestamp when the device last announced itself
	LastAnnounced string `json:"lastannounced"`

	// LastSeen is the timestamp when the device was last seen by the gateway
	LastSeen string `json:"lastseen"`

	// ETag is used for caching and resource versioning
	ETag string `json:"etag"`

	// ManufactureName is the name of the device manufacturer
	ManufactureName string `json:"manufacturername"`

	// Name is the user-assigned name of the light
	Name string `json:"name"`

	// ModelID is the model identifier of the light
	ModelID string `json:"modelid"`

	// PowerUp defines the behavior of the light when powered on
	PowerUp *int `json:"powerup,omitempty"`

	// SwVersion is the firmware version running on the light
	SwVersion string `json:"swversion"`

	// Type is the type of the light (e.g., "Extended color light", "Dimmable light")
	Type string `json:"type"`

	// UniqueID is the unique identifier for this light
	UniqueID string `json:"uniqueid"`

	// State contains the current state of the light
	State LightState `json:"state"`
}

// LightState represents the current state of a light device.
// This includes properties like on/off status, brightness, color, and other settings.
// All fields are pointers to allow for partial updates when changing state.
type LightState struct {
	// On indicates whether the light is turned on (true) or off (false)
	On *bool `json:"on,omitempty"`

	// Brightness is the current brightness level (0-255)
	Brightness *uint8 `json:"bri,omitempty"`

	// Hue is the current hue value (0-65535)
	Hue *uint16 `json:"hue,omitempty"`

	// Saturation is the current saturation value (0-255)
	Saturation *uint8 `json:"sat,omitempty"`

	// ColorTemperature is the current color temperature in mireds
	ColorTemperature *int `json:"ct,omitempty"`

	// XY contains the current color in CIE xy color space coordinates
	XY *[2]float64 `json:"xy,omitempty"`

	// Alert is the current alert effect ("none", "select", "lselect")
	Alert *string `json:"alert,omitempty"`

	// ColorMode indicates which color mode is active ("hs", "xy", "ct")
	ColorMode *string `json:"colormode,omitempty"`

	// Effect is the current effect running on the light
	Effect *string `json:"effect,omitempty"`

	// Speed is the speed of the current effect
	Speed *uint8 `json:"speed,omitempty"`

	// Reachable indicates whether the light is reachable by the gateway
	Reachable *bool `json:"reachable,omitempty"`
}

// GetLight retrieves detailed information about a specific light from the deCONZ gateway.
//
// Parameters:
//   - id: The identifier of the light to retrieve
//
// Returns:
//   - *Light: A pointer to the retrieved Light structure
//   - error: Any error encountered during the API request
func (ac *ApiClient) GetLight(id string) (*Light, error) {
	return client.Get[Light](ac.buildUrl("/lights/" + id))
}

// SetLightState updates the state of a light with the provided settings.
// This is the base method used by other light control methods.
//
// Parameters:
//   - id: The identifier of the light to update
//   - state: A pointer to a LightState structure containing the desired state changes
//
// Returns:
//   - error: Any error encountered during the API request
func (ac *ApiClient) SetLightState(id string, state *LightState) error {
	_, err := client.Put[any](ac.buildUrl("/lights/"+id+"/state"), *state)
	return err
}

// SetLightOn turns a light on or off.
//
// Parameters:
//   - id: The identifier of the light to control
//   - on: Boolean value indicating whether to turn the light on (true) or off (false)
//
// Returns:
//   - error: Any error encountered during the API request
func (ac *ApiClient) SetLightOn(id string, on bool) error {
	return ac.SetLightState(id, &LightState{
		On: &on,
	})
}

// SetLightBrightness sets the brightness of a light.
// If brightness is 0, the light will be turned off.
// If brightness is greater than 0, the light will be turned on and set to the specified brightness.
//
// Parameters:
//   - id: The identifier of the light to control
//   - brightness: The desired brightness level as a percentage (0-100)
//
// Returns:
//   - error: Any error encountered during the API request
func (ac *ApiClient) SetLightBrightness(id string, brightness int) error {
	state := new(LightState)
	f := false
	state.On = &f

	// convert percentage to value
	value := uint8(math.Round(float64(brightness * 255.0 / 100.0)))
	if value > 0 {
		t := true
		state.On = &t
		state.Brightness = &value
	}

	return ac.SetLightState(id, state)
}

// SetLightColorTemperature sets the color temperature of a light.
// The color temperature is specified in mireds (micro reciprocal degrees).
// Lower values represent cooler (more blue) light, higher values represent warmer (more orange) light.
//
// Parameters:
//   - id: The identifier of the light to control
//   - mired: The desired color temperature in mireds
//
// Returns:
//   - error: Any error encountered during the API request
func (ac *ApiClient) SetLightColorTemperature(id string, mired int) error {
	return ac.SetLightState(id, &LightState{
		ColorTemperature: &mired,
	})
}
