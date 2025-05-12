// Package deviceConfiguration provides functionality for loading, parsing, and managing
// device configuration files. These configurations define how different Zigbee devices
// (particularly remote controls and switches) map their button events to HomeKit actions.
package deviceConfiguration

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
	"os"
	"path/filepath"
)

// ButtonEvent represents a type of button press event.
// This is used to map deCONZ button events to HomeKit button events.
type ButtonEvent string

// Constants defining the different types of button press events.
const (
	// ButtonSinglePress represents a single press of a button
	ButtonSinglePress ButtonEvent = "SINGLE_PRESS"

	// ButtonDoublePress represents a double press of a button
	ButtonDoublePress ButtonEvent = "DOUBLE_PRESS"

	// ButtonLongPress represents a long press of a button
	ButtonLongPress ButtonEvent = "LONG_PRESS"
)

// ButtonConfiguration represents the configuration for a single button on a device.
// It defines the button's name and how its raw events map to button press types.
type ButtonConfiguration struct {
	// Name is a human-readable name for the button (e.g., "Top Button", "Power Button")
	Name string `json:"name"`

	// EventMap maps raw deCONZ event codes to button press types
	// The keys are strings like "1001" and the values are ButtonEvent constants
	EventMap map[string]ButtonEvent `json:"eventMap"`
}

// DeviceConfiguration represents the complete configuration for a device model.
// It includes metadata about the device and configurations for all its buttons.
type DeviceConfiguration struct {
	// SchemaVersion is the version of the configuration schema
	SchemaVersion string `json:"schemaVersion"`

	// Manufacturer is the name of the device manufacturer
	Manufacturer string `json:"manufacturer"`

	// Models is a list of model identifiers that this configuration applies to
	Models []string `json:"models"`

	// Description is a human-readable description of the device
	Description string `json:"description"`

	// Buttons is a list of button configurations for this device
	Buttons []ButtonConfiguration `json:"buttons"`
}

// SaveToFile saves the device configuration to a JSON file.
// The file is formatted with pretty-printing for readability.
//
// Parameters:
//   - file: The path to the file to save to
//
// Returns:
//   - error: An error if the file could not be saved
func (dc *DeviceConfiguration) SaveToFile(file string) error {
	// Convert the configuration to JSON
	data, err := json.Marshal(dc)
	if err != nil {
		return err
	}

	// Format the JSON for readability
	prettyData := pretty.Pretty(data)

	// Write the formatted JSON to the file
	return os.WriteFile(file, prettyData, 0644)
}

// LoadFromDirectory loads all device configurations from JSON files in a directory.
// It returns a map of model identifiers to their corresponding configurations.
//
// Parameters:
//   - dir: The directory to load configuration files from
//
// Returns:
//   - map[string]DeviceConfiguration: A map of model identifiers to device configurations
//   - error: An error if the directory could not be read
func LoadFromDirectory(dir string) (map[string]DeviceConfiguration, error) {
	configMap := make(map[string]DeviceConfiguration)

	// Find all JSON files in the specified directory
	files, err := filepath.Glob(dir + "/*.json")
	if err != nil {
		return nil, err
	}

	// Process each configuration file
	for _, fileName := range files {
		// Read the file contents
		if file, err := os.ReadFile(fileName); err == nil {
			// Parse the JSON into a DeviceConfiguration
			config := new(DeviceConfiguration)
			if err = json.Unmarshal(file, config); err == nil {
				// Add the configuration to the map for each model it applies to
				for _, model := range config.Models {
					configMap[model] = *config
				}
			} else {
				// Log an error if the file couldn't be parsed
				fmt.Printf("Error reading device configuration file %s: %s\n", fileName, err)
			}
		} else {
			// Log an error if the file couldn't be read
			fmt.Printf("Error reading device configuration file %s: %s\n", fileName, err)
		}
	}

	return configMap, nil
}

// SplitEventId splits a button event ID into a button number and an event code.
// For example, "1001" would be split into "10" (button number) and "01" (event code).
// This is used to identify which button was pressed and what type of press it was.
//
// Parameters:
//   - event: The button event ID to split
//
// Returns:
//   - string: The button number
//   - string: The event code
func SplitEventId(event string) (string, string) {
	// The last 3 characters are the event code, the rest is the button number
	prefix := event[:len(event)-3]
	suffix := event[len(event)-3:]
	return prefix, suffix
}
