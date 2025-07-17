//go:build !windows
// +build !windows

package gpu

import "fmt"

// GetWindowsGPUInfo stub for non-Windows platforms (remove idiot errors)
func GetWindowsGPUInfo() ([]GPUInfo, error) {
	return nil, fmt.Errorf("Windows GPU monitoring not available on this platform")
}
