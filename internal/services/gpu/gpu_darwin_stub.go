//go:build !darwin
// +build !darwin

package gpu

import "fmt"

// GetDarwinGPUInfo for non-Darwin platforms (remove idiot errors)
func GetDarwinGPUInfo() ([]GPUInfo, error) {
	return nil, fmt.Errorf("Darwin GPU monitoring not available on this platform")
}
