//go:build linux
// +build linux

package gpu

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// GetLinuxGPUInfo retrieves GPU information on Linux
func GetLinuxGPUInfo() ([]GPUInfo, error) {
	var gpus []GPUInfo

	if nvidiaGPUs, err := getLinuxNVIDIAGPUs(); err == nil {
		gpus = append(gpus, nvidiaGPUs...)
	}

	if amdGPUs, err := getLinuxAMDGPUs(); err == nil {
		gpus = append(gpus, amdGPUs...)
	}

	if intelGPUs, err := getLinuxIntelGPUs(); err == nil {
		gpus = append(gpus, intelGPUs...)
	}

	if len(gpus) == 0 {
		if pciGPUs, err := getLinuxPCIGPUs(); err == nil {
			gpus = append(gpus, pciGPUs...)
		}
	}

	if len(gpus) == 0 {
		return nil, fmt.Errorf("no GPUs detected on Linux")
	}

	return gpus, nil
}

func getLinuxNVIDIAGPUs() ([]GPUInfo, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=name,driver_version,memory.total,memory.used,memory.free,temperature.gpu,utilization.gpu,power.draw,power.limit,clocks.current.graphics", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("nvidia-smi not available or failed: %v", err)
	}

	var gpus []GPUInfo
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 10 {
			continue
		}

		gpu := GPUInfo{
			Name:      strings.TrimSpace(parts[0]),
			Vendor:    "NVIDIA",
			Driver:    strings.TrimSpace(parts[1]),
			Available: true,
		}

		if memTotal, err := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 64); err == nil {
			gpu.MemoryTotal = memTotal * 1024 * 1024
		}
		if memUsed, err := strconv.ParseUint(strings.TrimSpace(parts[3]), 10, 64); err == nil {
			gpu.MemoryUsed = memUsed * 1024 * 1024
		}
		if memFree, err := strconv.ParseUint(strings.TrimSpace(parts[4]), 10, 64); err == nil {
			gpu.MemoryFree = memFree * 1024 * 1024
		}

		if temp, err := strconv.ParseFloat(strings.TrimSpace(parts[5]), 64); err == nil {
			gpu.Temperature = temp
		}

		if usage, err := strconv.ParseFloat(strings.TrimSpace(parts[6]), 64); err == nil {
			gpu.Usage = usage
		}

		if power, err := strconv.ParseUint(strings.TrimSpace(parts[7]), 10, 64); err == nil {
			gpu.PowerDraw = power
		}

		if powerLimit, err := strconv.ParseUint(strings.TrimSpace(parts[8]), 10, 64); err == nil {
			gpu.PowerLimit = powerLimit
		}

		if clockSpeed, err := strconv.ParseUint(strings.TrimSpace(parts[9]), 10, 64); err == nil {
			gpu.ClockSpeed = clockSpeed
		}

		gpus = append(gpus, gpu)
	}

	return gpus, nil
}

func getLinuxAMDGPUs() ([]GPUInfo, error) {
	var gpus []GPUInfo

	cardDirs, err := filepath.Glob("/sys/class/drm/card*/device")
	if err != nil {
		return nil, fmt.Errorf("failed to glob DRM card directories: %v", err)
	}

	for _, cardDir := range cardDirs {
		vendorPath := filepath.Join(cardDir, "vendor")
		if vendorData, err := os.ReadFile(vendorPath); err == nil {
			vendor := strings.TrimSpace(string(vendorData))
			if vendor != "0x1002" {
				continue
			}
		}

		gpu := GPUInfo{
			Vendor:    "AMD",
			Available: true,
		}

		if devicePath := filepath.Join(cardDir, "device"); fileExists(devicePath) {
			if deviceData, err := os.ReadFile(devicePath); err == nil {
				deviceID := strings.TrimSpace(string(deviceData))
				gpu.Name = fmt.Sprintf("AMD GPU (Device ID: %s)", deviceID)
			}
		}

		tempPath := filepath.Join(cardDir, "hwmon/hwmon*/temp1_input")
		if matches, _ := filepath.Glob(tempPath); len(matches) > 0 {
			if tempData, err := os.ReadFile(matches[0]); err == nil {
				if temp, err := strconv.ParseInt(strings.TrimSpace(string(tempData)), 10, 64); err == nil {
					gpu.Temperature = float64(temp) / 1000.0
				}
			}
		}

		powerPath := filepath.Join(cardDir, "hwmon/hwmon*/power1_average")
		if matches, _ := filepath.Glob(powerPath); len(matches) > 0 {
			if powerData, err := os.ReadFile(matches[0]); err == nil {
				if power, err := strconv.ParseUint(strings.TrimSpace(string(powerData)), 10, 64); err == nil {
					gpu.PowerDraw = power / 1000000
				}
			}
		}

		gpus = append(gpus, gpu)
	}

	return gpus, nil
}

func getLinuxIntelGPUs() ([]GPUInfo, error) {
	var gpus []GPUInfo

	cmd := exec.Command("intel_gpu_top", "-o", "-")
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Intel") {
				gpu := GPUInfo{
					Name:      "Intel Integrated GPU",
					Vendor:    "Intel",
					Available: true,
				}
				gpus = append(gpus, gpu)
				break
			}
		}
	}

	if len(gpus) == 0 {
		if pciGPUs, err := getLinuxPCIGPUs(); err == nil {
			for _, gpu := range pciGPUs {
				if gpu.Vendor == "Intel" {
					gpus = append(gpus, gpu)
				}
			}
		}
	}

	return gpus, nil
}

func getLinuxPCIGPUs() ([]GPUInfo, error) {
	cmd := exec.Command("lspci", "-nn")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("lspci not available: %v", err)
	}

	var gpus []GPUInfo
	lines := strings.Split(string(output), "\n")

	vgaRegex := regexp.MustCompile(`(?i)(vga|3d|display).*controller`)

	for _, line := range lines {
		if vgaRegex.MatchString(line) {
			gpu := GPUInfo{
				Name:      extractGPUNameFromPCI(line),
				Vendor:    getVendorFromName(line),
				Available: true,
			}
			gpus = append(gpus, gpu)
		}
	}

	return gpus, nil
}

func extractGPUNameFromPCI(line string) string {
	parts := strings.Split(line, ":")
	if len(parts) >= 3 {
		name := strings.Join(parts[2:], ":")
		name = strings.TrimSpace(name)

		re := regexp.MustCompile(`\[[0-9a-fA-F:]+\]`)
		name = re.ReplaceAllString(name, "")
		name = strings.TrimSpace(name)

		return name
	}

	return "Unknown GPU"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func readLinuxGPUFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func readLinuxGPUFileUint(path string) (uint64, error) {
	data, err := readLinuxGPUFile(path)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(data, 10, 64)
}
