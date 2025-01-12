package utils

import "fmt"

// Convert version information to a string
func ParseVersion(v uint32) string {
	// Major version: 1 byte, Minor version: 1 byte, Patch version: 2 bytes
	// Upper 8 bits
	major := v >> 24
	// Next 8 bits
	minor := (v >> 16) & 0xff
	// Lower 16 bits
	patch := v & 0xffff
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
