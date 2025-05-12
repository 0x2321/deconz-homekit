// Package helper provides utility functions for common operations used throughout the application.
// This includes mathematical conversions between different value representations used by
// deCONZ and HomeKit.
package helper

import "math"

// RawToDeg converts a raw 16-bit value (0-65535) to degrees (0-360).
// This is used for converting color hue values from deCONZ's raw format to degrees.
//
// Parameters:
//   - raw: The raw value to convert (0-65535)
//
// Returns:
//   - float64: The equivalent value in degrees (0-360)
func RawToDeg(raw uint16) float64 {
	return float64((360 / 65535) * raw)
}

// DegToRaw converts a degree value (0-360) to a raw 16-bit value (0-65535).
// This is used for converting color hue values from degrees to deCONZ's raw format.
//
// Parameters:
//   - deg: The degree value to convert (0-360)
//
// Returns:
//   - uint16: The equivalent raw value (0-65535)
func DegToRaw(deg float64) uint16 {
	return uint16(math.Round((65535 / 360) * deg))
}

// RawToDec converts a raw 8-bit value (0-255) to a decimal percentage (0-100).
// This is used for converting brightness values from deCONZ's raw format to percentages
// that are more user-friendly in HomeKit.
//
// Parameters:
//   - raw: The raw value to convert (0-255)
//
// Returns:
//   - float64: The equivalent value as a percentage (0-100)
func RawToDec(raw uint8) float64 {
	return float64(raw) * 100.0 / 255.0
}

// DecToRaw converts a decimal percentage (0-100) to a raw 8-bit value (0-255).
// This is used for converting brightness percentages from HomeKit to deCONZ's raw format.
//
// Parameters:
//   - dec: The percentage value to convert (0-100)
//
// Returns:
//   - uint8: The equivalent raw value (0-255)
func DecToRaw(dec int) uint8 {
	return uint8(float64(dec) * 255.0 / 100.0)
}
