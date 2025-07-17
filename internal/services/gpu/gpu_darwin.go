//go:build darwin
// +build darwin

package gpu

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// GetDarwinGPUInfo retrieves GPU information on macOS
func GetDarwinGPUInfo() ([]GPUInfo, error) {
	var gpus []GPUInfo

	if profilerGPUs, err := getDarwinSystemProfilerGPUs(); err == nil {
		gpus = append(gpus, profilerGPUs...)
	}

	if len(gpus) == 0 {
		if ioregGPUs, err := getDarwinIORegGPUs(); err == nil {
			gpus = append(gpus, ioregGPUs...)
		}
	}

	if metalGPUs, err := getDarwinMetalGPUs(); err == nil {
		for i, gpu := range gpus {
			for _, metalGPU := range metalGPUs {
				if strings.Contains(strings.ToLower(gpu.Name), strings.ToLower(metalGPU.Name)) {
					gpus[i] = metalGPU
					break
				}
			}
		}

		for _, metalGPU := range metalGPUs {
			found := false
			for _, gpu := range gpus {
				if strings.Contains(strings.ToLower(gpu.Name), strings.ToLower(metalGPU.Name)) {
					found = true
					break
				}
			}
			if !found {
				gpus = append(gpus, metalGPU)
			}
		}
	}

	if len(gpus) == 0 {
		return nil, fmt.Errorf("no GPUs detected on macOS")
	}

	return gpus, nil
}

func getDarwinSystemProfilerGPUs() ([]GPUInfo, error) {
	cmd := exec.Command("system_profiler", "SPDisplaysDataType")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("system_profiler failed: %v", err)
	}

	var gpus []GPUInfo
	lines := strings.Split(string(output), "\n")
	var currentGPU GPUInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "Chipset Model:") {
			if currentGPU.Name != "" {
				gpus = append(gpus, currentGPU)
			}
			currentGPU = GPUInfo{
				Name:      extractValue(line, "Chipset Model:"),
				Vendor:    "",
				Available: true,
			}
			currentGPU.Vendor = getVendorFromName(currentGPU.Name)
		} else if strings.Contains(line, "Type:") {

			gpuType := extractValue(line, "Type:")
			if gpuType != "" {
				currentGPU.Name = fmt.Sprintf("%s (%s)", currentGPU.Name, gpuType)
			}
		} else if strings.Contains(line, "VRAM (Total):") {
			memStr := extractValue(line, "VRAM (Total):")
			currentGPU.MemoryTotal = parseMemorySize(memStr)
		} else if strings.Contains(line, "VRAM (Dynamic, Max):") {
			memStr := extractValue(line, "VRAM (Dynamic, Max):")
			if currentGPU.MemoryTotal == 0 {
				currentGPU.MemoryTotal = parseMemorySize(memStr)
			}
		} else if strings.Contains(line, "Vendor:") {
			vendor := extractValue(line, "Vendor:")
			if vendor != "" {
				currentGPU.Vendor = vendor
			}
		} else if strings.Contains(line, "Device ID:") {
			deviceID := extractValue(line, "Device ID:")
			if deviceID != "" {
				currentGPU.Name = fmt.Sprintf("%s (ID: %s)", currentGPU.Name, deviceID)
			}
		} else if strings.Contains(line, "Revision ID:") {
			revisionID := extractValue(line, "Revision ID:")
			if revisionID != "" {
				currentGPU.Driver = fmt.Sprintf("Revision %s", revisionID)
			}
		}
	}

	if currentGPU.Name != "" {
		gpus = append(gpus, currentGPU)
	}

	return gpus, nil
}

func getDarwinIORegGPUs() ([]GPUInfo, error) {
	cmd := exec.Command("ioreg", "-l", "-p", "IOPCIDevice")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ioreg failed: %v", err)
	}

	var gpus []GPUInfo
	lines := strings.Split(string(output), "\n")

	var currentGPU GPUInfo
	inGPUSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "\"class-code\"") && (strings.Contains(line, "0x30000") || strings.Contains(line, "0x38000")) {
			inGPUSection = true
			currentGPU = GPUInfo{
				Available: true,
			}
		} else if inGPUSection && strings.Contains(line, "\"model\"") {
			currentGPU.Name = extractIORegValue(line, "\"model\"")
			currentGPU.Vendor = getVendorFromName(currentGPU.Name)
		} else if inGPUSection && strings.Contains(line, "\"vendor-id\"") {
			vendorID := extractIORegValue(line, "\"vendor-id\"")
			currentGPU.Vendor = getVendorFromID(vendorID)
		} else if inGPUSection && strings.Contains(line, "\"device-id\"") {
			deviceID := extractIORegValue(line, "\"device-id\"")
			if currentGPU.Name == "" {
				currentGPU.Name = fmt.Sprintf("GPU Device %s", deviceID)
			}
		} else if inGPUSection && strings.Contains(line, "}") && strings.Count(line, "}") > strings.Count(line, "{") {
			if currentGPU.Name != "" {
				gpus = append(gpus, currentGPU)
			}
			inGPUSection = false
			currentGPU = GPUInfo{}
		}
	}

	return gpus, nil
}

func getDarwinMetalGPUs() ([]GPUInfo, error) {
	return []GPUInfo{}, nil
}

func getDarwinGPUPowerUsage() (float64, error) {
	cmd := exec.Command("powermetrics", "-n", "1", "-s", "gpu_power")
	output, err := cmd.Output()
	if err != nil {
		return 0.0, fmt.Errorf("powermetrics not available: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "GPU Power") {
			re := regexp.MustCompile(`(\d+\.?\d*)\s*mW`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				if power, err := strconv.ParseFloat(matches[1], 64); err == nil {
					return power / 1000.0, nil
				}
			}
		}
	}

	return 0.0, fmt.Errorf("GPU power usage not available")
}

func extractValue(line, key string) string {
	if idx := strings.Index(line, key); idx != -1 {
		value := strings.TrimSpace(line[idx+len(key):])
		return value
	}
	return ""
}

func extractIORegValue(line, key string) string {
	if idx := strings.Index(line, key); idx != -1 {
		remaining := line[idx+len(key):]
		if idx := strings.Index(remaining, "="); idx != -1 {
			value := strings.TrimSpace(remaining[idx+1:])
			value = strings.Trim(value, "\"<>")
			return value
		}
	}
	return ""
}

func getVendorFromID(vendorID string) string {
	vendorMap := map[string]string{
		"0x1002": "AMD",
		"0x10de": "NVIDIA",
		"0x8086": "Intel",
		"0x1414": "Microsoft",
		"0x5333": "S3",
		"0x1039": "SiS",
		"0x102b": "Matrox",
		"0x121a": "3dfx",
	}

	if vendor, exists := vendorMap[vendorID]; exists {
		return vendor
	}

	return "Unknown"
}

func parseMemorySize(memStr string) uint64 {
	memStr = strings.ToLower(strings.TrimSpace(memStr))

	var multiplier uint64 = 1
	var numStr string

	if strings.Contains(memStr, "gb") {
		multiplier = 1024 * 1024 * 1024
		numStr = strings.Replace(memStr, "gb", "", -1)
	} else if strings.Contains(memStr, "mb") {
		multiplier = 1024 * 1024
		numStr = strings.Replace(memStr, "mb", "", -1)
	} else if strings.Contains(memStr, "kb") {
		multiplier = 1024
		numStr = strings.Replace(memStr, "kb", "", -1)
	} else {
		numStr = memStr
	}

	numStr = strings.TrimSpace(numStr)
	if num, err := strconv.ParseFloat(numStr, 64); err == nil {
		return uint64(num * float64(multiplier))
	}

	return 0
}
