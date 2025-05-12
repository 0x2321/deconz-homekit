// Package accessoryManager provides functionality for creating and managing HomeKit accessories
// that represent deCONZ devices.
package accessoryManager

import (
	"math/big"
	"strings"
)

// uniqueIdToHomeKitId converts a deCONZ unique ID (which is typically a MAC address or similar
// identifier in hexadecimal format with colons or hyphens) to a uint64 that can be used as
// a HomeKit accessory ID.
//
// The function removes colons and hyphens from the ID, interprets the resulting string as
// a hexadecimal number, and converts it to a uint64.
//
// Parameters:
//   - id: The deCONZ unique ID to convert
//
// Returns:
//   - uint64: The converted HomeKit accessory ID
func uniqueIdToHomeKitId(id string) uint64 {
	// Remove colons and hyphens from the ID
	numberStr := strings.ReplaceAll(id, ":", "")
	numberStr = strings.ReplaceAll(numberStr, "-", "")

	// Convert the hexadecimal string to a big integer
	n := new(big.Int)
	n.SetString(numberStr, 16)

	// Return the uint64 representation
	return n.Uint64()
}

// onOffStr is a map that converts boolean values to "on" or "off" strings.
// This is used for logging and for setting device states in a human-readable format.
var onOffStr = map[bool]string{
	true:  "on",
	false: "off",
}

// boolToInt is a map that converts boolean values to 1 or 0 integers.
// This is used for converting boolean states to numeric values for HomeKit characteristics.
var boolToInt = map[bool]int{
	true:  1,
	false: 0,
}
