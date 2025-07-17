//go:build !linux
// +build !linux

package gpu

import "fmt"

// GetLinuxGPUInfo stub for non-Linux platforms (remove idiot errors)
func GetLinuxGPUInfo() ([]GPUInfo, error) {
	return nil, fmt.Errorf("Linux GPU monitoring not available on this platform")
}
