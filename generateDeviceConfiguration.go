//go:build ignore
// +build ignore

package main

import (
	"deconz-homekit/internal/client"
	deviceConfiguration "deconz-homekit/internal/device_configuration"
	"fmt"
	"log"
	"maps"
	"regexp"
	"slices"
	"strings"
)

type ButtonMap struct {
	Vendor   string               `json:"vendor"`
	ModelIds []string             `json:"modelids"`
	Buttons  *[]map[string]string `json:"buttons"`
	Doc      string               `json:"doc"`
	Map      [][]interface{}      `json:"map"`
}

type MapFile struct {
	Buttons map[string]int       `json:"buttons"`
	Actions map[string]int       `json:"buttonActions"`
	Maps    map[string]ButtonMap `json:"maps"`
}

func main() {
	data, err := client.Get[MapFile]("https://raw.githubusercontent.com/dresden-elektronik/deconz-rest-plugin/master/button_maps.json")
	if err != nil {
		log.Fatalf("error getting file: %+v", err)
	}

	// for every device
	for _, device := range data.Maps {
		newDeviceConfig := &deviceConfiguration.DeviceConfiguration{}
		newDeviceConfig.SchemaVersion = "1.0"
		newDeviceConfig.Manufacturer = device.Vendor
		newDeviceConfig.Models = device.ModelIds
		newDeviceConfig.Description = device.Doc

		// for every button map in device
		buttonsMap := make(map[string]deviceConfiguration.ButtonConfiguration)
		for _, button := range device.Map {
			buttonId := button[5].(string)
			actionId := button[6].(string)
			eventId := data.Buttons[buttonId] + data.Actions[actionId]

			if eventId > 1000 {
				// create button if not exist in new configuration
				if _, ok := buttonsMap[buttonId]; !ok {

					// find name
					buttonName := fmt.Sprintf("Button %d", len(buttonsMap)+1)
					if device.Buttons != nil {
						for _, deviceButton := range *device.Buttons {
							if slices.Collect(maps.Keys(deviceButton))[0] == buttonId {
								buttonName = slices.Collect(maps.Values(deviceButton))[0]
							}
						}
					}

					// create button configuration
					buttonsMap[buttonId] = deviceConfiguration.ButtonConfiguration{
						Name:     buttonName,
						EventMap: make(map[string]deviceConfiguration.ButtonEvent),
					}
				}

				// add event
				switch button[6] {
				case "S_BUTTON_ACTION_SHORT_RELEASED":
					buttonsMap[buttonId].EventMap[fmt.Sprintf("%d", eventId)] = deviceConfiguration.ButtonSinglePress
				case "S_BUTTON_ACTION_DOUBLE_PRESS":
					buttonsMap[buttonId].EventMap[fmt.Sprintf("%d", eventId)] = deviceConfiguration.ButtonDoublePress
				case "S_BUTTON_ACTION_LONG_RELEASED":
					buttonsMap[buttonId].EventMap[fmt.Sprintf("%d", eventId)] = deviceConfiguration.ButtonLongPress
				}
			}
		}

		// add all buttons to the new configuration
		for _, button := range buttonsMap {
			if len(button.EventMap) > 0 {
				newDeviceConfig.Buttons = append(newDeviceConfig.Buttons, button)
			}
		}

		// save new configuration if there are any buttons
		if len(newDeviceConfig.Buttons) > 0 {
			name := strings.ToLower(device.Vendor) + "_" + strings.ToLower(device.ModelIds[0])
			re := regexp.MustCompile(`[^a-z0-9]+`)
			name = re.ReplaceAllString(name, "_")
			if err := newDeviceConfig.SaveToFile("./devices/" + name + ".json"); err != nil {
				fmt.Printf("error saving device configuration: %+v\n", err)
			}
		}
	}
}
