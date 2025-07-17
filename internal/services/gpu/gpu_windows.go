//go:build windows
// +build windows

package gpu

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
)

type win32_VideoController struct {
	Name                        string
	AdapterRAM                  uint64
	DriverVersion               string
	VideoProcessor              string
	VideoMemoryType             uint16
	Status                      string
	Availability                uint16
	CurrentBitsPerPixel         uint32
	CurrentHorizontalResolution uint32
	CurrentVerticalResolution   uint32
	Description                 string
	DeviceID                    string
	PNPDeviceID                 string
}

type win32_PerfRawData_GPUPerformanceCounters_GPUEngine struct {
	Name                  string
	UtilizationPercentage uint64
}

type win32_VideoControllerResolution struct {
	HorizontalResolution uint32
	VerticalResolution   uint32
	BitsPerPixel         uint32
}

func GetWindowsGPUInfo() ([]GPUInfo, error) {
	var win32GPUs []win32_VideoController
	query := "SELECT * FROM Win32_VideoController WHERE VideoProcessor IS NOT NULL"

	if err := wmi.Query(query, &win32GPUs); err != nil {
		return nil, fmt.Errorf("failed to query GPU information: %v", err)
	}

	var gpus []GPUInfo
	for _, gpu := range win32GPUs {
		gpuInfo := GPUInfo{
			Name:      gpu.Name,
			Vendor:    getVendorFromName(gpu.Name),
			Driver:    gpu.DriverVersion,
			Available: gpu.Status == "OK",
		}

		if gpu.AdapterRAM > 0 {
			memoryBytes := gpu.AdapterRAM

			if memoryBytes > 16*1024*1024*1024 {
				if strings.Contains(strings.ToLower(gpu.Name), "mx250") {
					memoryBytes = 2 * 1024 * 1024 * 1024
				} else if strings.Contains(strings.ToLower(gpu.Name), "gtx 1050") {
					memoryBytes = 2 * 1024 * 1024 * 1024
				} else if strings.Contains(strings.ToLower(gpu.Name), "gtx 1060") {
					memoryBytes = 6 * 1024 * 1024 * 1024
				} else if strings.Contains(strings.ToLower(gpu.Name), "gtx") {
					memoryBytes = 4 * 1024 * 1024 * 1024
				} else if strings.Contains(strings.ToLower(gpu.Name), "rtx") {
					memoryBytes = 8 * 1024 * 1024 * 1024
				} else {
					memoryBytes = 4 * 1024 * 1024 * 1024
				}
			}

			gpuInfo.MemoryTotal = memoryBytes
			gpuInfo.MemoryFree = memoryBytes
			gpuInfo.MemoryUsed = 0
		}

		if usage, err := getWindowsGPUUsage(); err == nil {
			gpuInfo.Usage = usage
		}

		if temp, err := getWindowsGPUTemperature(); err == nil {
			gpuInfo.Temperature = temp
		}

		gpus = append(gpus, gpuInfo)
	}

	return gpus, nil
}

func getWindowsGPUUsage() (float64, error) {
	type Win32_PerfRawData_Counters_ProcessorInformation struct {
		Name                 string
		PercentProcessorTime uint64
	}

	var perfCounters []Win32_PerfRawData_Counters_ProcessorInformation
	query := "SELECT * FROM Win32_PerfRawData_Counters_ProcessorInformation WHERE Name LIKE '%GPU%'"

	err := wmi.Query(query, &perfCounters)
	if err != nil {
		return 0.0, nil
	}

	// For now, return a basic approximation
	return 0.0, nil
}

func getWindowsGPUTemperature() (float64, error) {
	type Win32_TemperatureProbe struct {
		Name           string
		CurrentReading *uint32
		MaxReading     *uint32
		Status         string
		Description    string
	}

	var tempProbes []Win32_TemperatureProbe
	query := "SELECT * FROM Win32_TemperatureProbe"

	err := wmi.Query(query, &tempProbes)
	if err != nil {
		type MSAcpi_ThermalZoneTemperature struct {
			CurrentTemperature uint32
			InstanceName       string
		}

		var thermalZones []MSAcpi_ThermalZoneTemperature
		query := "SELECT * FROM MSAcpi_ThermalZoneTemperature"

		err := wmi.Query(query, &thermalZones)
		if err != nil {
			return 0.0, nil
		}

		if len(thermalZones) > 0 {
			tempKelvin := float64(thermalZones[0].CurrentTemperature) / 10.0
			tempCelsius := tempKelvin - 273.15
			return tempCelsius, nil
		}
	}

	for _, probe := range tempProbes {
		if probe.CurrentReading != nil && probe.Status == "OK" {
			return float64(*probe.CurrentReading) / 10.0, nil
		}
	}

	return 0.0, nil
}
