package gpu

import (
	"fmt"
	"runtime"
	"strings"

	"syspulse/internal/errors"
)

type GPUInfo struct {
	Name        string  `json:"name"`
	Vendor      string  `json:"vendor"`
	MemoryTotal uint64  `json:"memory_total"`
	MemoryUsed  uint64  `json:"memory_used"`
	MemoryFree  uint64  `json:"memory_free"`
	Temperature float64 `json:"temperature"`
	Usage       float64 `json:"usage"`
	Driver      string  `json:"driver"`
	ClockSpeed  uint64  `json:"clock_speed"`
	FanSpeed    uint64  `json:"fan_speed"`
	PowerDraw   uint64  `json:"power_draw"`
	PowerLimit  uint64  `json:"power_limit"`
	Available   bool    `json:"available"`
}

func GetGPUInfo() ([]GPUInfo, error) {
	switch runtime.GOOS {
	case "windows":
		return GetWindowsGPUInfo()
	case "linux":
		return GetLinuxGPUInfo()
	case "darwin":
		return GetDarwinGPUInfo()
	default:
		return nil, errors.NewAppError(errors.MonitorError,
			fmt.Sprintf("Unsupported operating system: %s", runtime.GOOS), nil)
	}
}

func getVendorFromName(name string) string {
	name = strings.ToLower(name)

	vendors := map[string]string{
		"nvidia":  "NVIDIA",
		"geforce": "NVIDIA",
		"gtx":     "NVIDIA",
		"rtx":     "NVIDIA",
		"quadro":  "NVIDIA",
		"tesla":   "NVIDIA",
		"amd":     "AMD",
		"radeon":  "AMD",
		"rx":      "AMD",
		"vega":    "AMD",
		"ryzen":   "AMD",
		"intel":   "Intel",
		"iris":    "Intel",
		"uhd":     "Intel",
		"hd":      "Intel",
	}

	for keyword, vendor := range vendors {
		if strings.Contains(name, keyword) {
			return vendor
		}
	}

	return "Unknown"
}

func formatMemorySize(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	if bytes >= GB {
		return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
	} else if bytes >= MB {
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	} else if bytes >= KB {
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	}

	return fmt.Sprintf("%d B", bytes)
}

func GetGPUCount() int {
	gpus, err := GetGPUInfo()
	if err != nil {
		return 0
	}
	return len(gpus)
}

func IsGPUAvailable() bool {
	gpus, err := GetGPUInfo()
	if err != nil {
		return false
	}

	for _, gpu := range gpus {
		if gpu.Available {
			return true
		}
	}

	return false
}
