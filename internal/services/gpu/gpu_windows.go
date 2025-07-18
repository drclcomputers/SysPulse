//go:build windows
// +build windows

package gpu

import (
	"fmt"
	"os/exec"
	"strconv"
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

type win32_PerfRawData_GPUPerformanceCounters_GPUMemory struct {
	Name           string
	DedicatedUsage uint64
	SharedUsage    uint64
}

type win32_PerfRawData_GPUPerformanceCounters_GPUProcessMemory struct {
	Name           string
	DedicatedUsage uint64
	SharedUsage    uint64
}

type win32_PerfRawData_GPUPerformanceCounters_GPUAdapter struct {
	Name           string
	DedicatedUsage uint64
	SharedUsage    uint64
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
	for i, gpu := range win32GPUs {
		gpuInfo := GPUInfo{
			Name:      gpu.Name,
			Vendor:    getVendorFromName(gpu.Name),
			Driver:    gpu.DriverVersion,
			Available: gpu.Status == "OK",
		}

		if gpu.AdapterRAM > 0 {
			memoryBytes := gpu.AdapterRAM

			if memoryBytes > 24*1024*1024*1024 {
				if strings.Contains(strings.ToLower(gpu.Name), "mx250") {
					memoryBytes = 2 * 1024 * 1024 * 1024
				} else if strings.Contains(strings.ToLower(gpu.Name), "gtx 1050") {
					memoryBytes = 2 * 1024 * 1024 * 1024
				} else if strings.Contains(strings.ToLower(gpu.Name), "gtx 1060") {
					memoryBytes = 6 * 1024 * 1024 * 1024
				} else if strings.Contains(strings.ToLower(gpu.Name), "gtx") {
					memoryBytes = 6 * 1024 * 1024 * 1024
				} else if strings.Contains(strings.ToLower(gpu.Name), "rtx") {
					memoryBytes = 8 * 1024 * 1024 * 1024
				} else {
					memoryBytes = 4 * 1024 * 1024 * 1024
				}
			}

			gpuInfo.MemoryTotal = memoryBytes

			if memUsed, err := getWindowsGPUMemoryUsage(i); err == nil {
				// Ensure memory used doesn't exceed total memory
				if memUsed > memoryBytes {
					memUsed = memoryBytes / 2 // Fallback to 50% if unrealistic
				}
				gpuInfo.MemoryUsed = memUsed
				gpuInfo.MemoryFree = memoryBytes - memUsed
			} else {
				// Try alternative approach for memory usage
				if memUsed, err := getAlternativeGPUMemoryUsage(); err == nil && memUsed > 0 {
					if memUsed > memoryBytes {
						memUsed = memoryBytes / 2
					}
					gpuInfo.MemoryUsed = memUsed
					gpuInfo.MemoryFree = memoryBytes - memUsed
				} else {
					gpuInfo.MemoryFree = memoryBytes
					gpuInfo.MemoryUsed = 0
				}
			}
		}

		if usage, err := getWindowsGPUUsage(i); err == nil {
			gpuInfo.Usage = usage
		}

		if temp, err := getWindowsGPUTemperature(); err == nil {
			gpuInfo.Temperature = temp
		}

		gpus = append(gpus, gpuInfo)
	}

	return gpus, nil
}

func getWindowsGPUMemoryUsage(gpuIndex int) (uint64, error) {
	if usage, err := getNvidiaMemoryUsage(); err == nil {
		return usage, nil
	}

	if usage, err := getGPUMemoryFromCommand(); err == nil {
		return usage, nil
	}

	var memCounters []win32_PerfRawData_GPUPerformanceCounters_GPUMemory
	query := "SELECT * FROM Win32_PerfRawData_GPUPerformanceCounters_GPUMemory"

	if err := wmi.Query(query, &memCounters); err == nil && len(memCounters) > gpuIndex {
		dedicatedUsage := memCounters[gpuIndex].DedicatedUsage
		if dedicatedUsage > 0 && dedicatedUsage < 100*1024*1024*1024 { // Sanity check < 100GB
			return dedicatedUsage, nil
		}
	}

	var adapterCounters []win32_PerfRawData_GPUPerformanceCounters_GPUAdapter
	query = "SELECT * FROM Win32_PerfRawData_GPUPerformanceCounters_GPUAdapter"

	if err := wmi.Query(query, &adapterCounters); err == nil && len(adapterCounters) > gpuIndex {
		dedicatedUsage := adapterCounters[gpuIndex].DedicatedUsage
		if dedicatedUsage > 0 && dedicatedUsage < 100*1024*1024*1024 { // Sanity check < 100GB
			return dedicatedUsage, nil
		}
	}

	if usage, err := getGPUMemoryFromWMI(); err == nil {
		return usage, nil
	}

	return 0, fmt.Errorf("GPU memory usage not available")
}

func getNvidiaMemoryUsage() (uint64, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=memory.used", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 0 {
		memStr := strings.TrimSpace(lines[0])
		if memStr == "" || memStr == "N/A" {
			return 0, fmt.Errorf("nvidia-smi returned empty memory usage")
		}

		if memMB, err := strconv.ParseUint(memStr, 10, 64); err == nil {
			memBytes := memMB * 1024 * 1024
			if memBytes > 100*1024*1024*1024 {
				return 0, fmt.Errorf("nvidia-smi returned unreasonable memory usage: %d GB", memBytes/(1024*1024*1024))
			}
			return memBytes, nil
		}
	}

	return 0, fmt.Errorf("failed to parse nvidia-smi output")
}

func getGPUMemoryFromCommand() (uint64, error) {
	cmd := exec.Command("powershell", "-Command",
		"Get-WmiObject -Class Win32_PerfRawData_GPUPerformanceCounters_GPUMemory | Select-Object -ExpandProperty DedicatedUsage")

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 0 {
		memStr := strings.TrimSpace(lines[0])
		if memBytes, err := strconv.ParseUint(memStr, 10, 64); err == nil {
			return memBytes, nil
		}
	}

	return 0, fmt.Errorf("failed to get GPU memory from PowerShell")
}

func getGPUMemoryFromWMI() (uint64, error) {
	cmd := exec.Command("powershell", "-Command", `
		$gpuMemory = Get-WmiObject -Query "SELECT * FROM Win32_PerfRawData_GPUPerformanceCounters_GPUMemory";
		if ($gpuMemory) {
			$gpuMemory | ForEach-Object { $_.DedicatedUsage }
		} else {
			$adapter = Get-WmiObject -Query "SELECT * FROM Win32_PerfRawData_GPUPerformanceCounters_GPUAdapter";
			if ($adapter) {
				$adapter | ForEach-Object { $_.DedicatedUsage }
			}
		}
	`)

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if memBytes, err := strconv.ParseUint(line, 10, 64); err == nil && memBytes > 0 {
			return memBytes, nil
		}
	}

	return 0, fmt.Errorf("failed to get GPU memory from WMI")
}

func getWindowsGPUUsage(gpuIndex int) (float64, error) {
	// Try nvidia-smi first for NVIDIA GPUs
	if usage, err := getNvidiaGPUUsage(); err == nil {
		return usage, nil
	}

	// Try PowerShell GPU usage approach
	if usage, err := getGPUUsageFromWPT(); err == nil {
		return usage, nil
	}

	// Try WMI GPU engine counters with proper normalization
	var engineCounters []win32_PerfRawData_GPUPerformanceCounters_GPUEngine
	query := "SELECT * FROM Win32_PerfRawData_GPUPerformanceCounters_GPUEngine"

	if err := wmi.Query(query, &engineCounters); err == nil {
		var totalUsage float64
		var engineCount int

		for _, counter := range engineCounters {
			if strings.Contains(strings.ToLower(counter.Name), "3d") ||
				strings.Contains(strings.ToLower(counter.Name), "graphics") {

				// The UtilizationPercentage from WMI is in raw form and needs normalization
				// It's typically a value that needs to be divided by 100 or 10000
				rawUsage := float64(counter.UtilizationPercentage)

				// Normalize the raw value
				var normalizedUsage float64
				if rawUsage > 10000 {
					normalizedUsage = rawUsage / 10000.0 // Some counters use 10000 as max
				} else if rawUsage > 100 {
					normalizedUsage = rawUsage / 100.0 // Some counters use 100 as max
				} else {
					normalizedUsage = rawUsage // Already normalized
				}

				// Cap at 100% to avoid crazy values
				if normalizedUsage > 100 {
					normalizedUsage = 100
				}

				totalUsage += normalizedUsage
				engineCount++
			}
		}

		if engineCount > 0 {
			avgUsage := totalUsage / float64(engineCount)
			// Additional sanity check
			if avgUsage > 100 {
				avgUsage = 100
			}
			return avgUsage, nil
		}
	}

	// Try alternative PowerShell approach for GPU usage
	if usage, err := getGPUUsageFromTaskManager(); err == nil {
		return usage, nil
	}

	return 0.0, fmt.Errorf("GPU usage not available")
}

func getNvidiaGPUUsage() (float64, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		return 0.0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 0 {
		usageStr := strings.TrimSpace(lines[0])
		if usage, err := strconv.ParseFloat(usageStr, 64); err == nil {
			return usage, nil
		}
	}

	return 0.0, fmt.Errorf("failed to parse nvidia-smi output")
}

func getGPUUsageFromWPT() (float64, error) {
	cmd := exec.Command("powershell", "-Command",
		"Get-Counter -Counter \"\\GPU Engine(*)\\Utilization Percentage\" -SampleInterval 1 -MaxSamples 1 | Select-Object -ExpandProperty CounterSamples | Select-Object -ExpandProperty CookedValue")

	output, err := cmd.Output()
	if err != nil {
		return 0.0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var totalUsage float64
	var sampleCount int

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if usage, err := strconv.ParseFloat(line, 64); err == nil {
			// Normalize the usage value
			if usage > 100 {
				usage = usage / 100.0 // Some values might be in 0-10000 range
			}
			if usage > 100 {
				usage = 100 // Cap at 100%
			}
			totalUsage += usage
			sampleCount++
		}
	}

	if sampleCount > 0 {
		avgUsage := totalUsage / float64(sampleCount)
		if avgUsage > 100 {
			avgUsage = 100
		}
		return avgUsage, nil
	}

	return 0.0, fmt.Errorf("failed to get GPU usage from WPT")
}

func getGPUUsageFromTaskManager() (float64, error) {
	// Try to get GPU usage via PowerShell using Task Manager-like approach
	cmd := exec.Command("powershell", "-Command", `
		$gpu = Get-WmiObject -Query "SELECT * FROM Win32_PerfRawData_GPUPerformanceCounters_GPUEngine";
		if ($gpu) {
			$totalUsage = 0;
			$count = 0;
			foreach ($engine in $gpu) {
				if ($engine.Name -like "*3D*" -or $engine.Name -like "*Graphics*") {
					$usage = $engine.UtilizationPercentage;
					if ($usage -gt 10000) { $usage = $usage / 10000 }
					elseif ($usage -gt 100) { $usage = $usage / 100 }
					if ($usage -gt 100) { $usage = 100 }
					$totalUsage += $usage;
					$count++;
				}
			}
			if ($count -gt 0) {
				$totalUsage / $count;
			}
		}
	`)

	output, err := cmd.Output()
	if err != nil {
		return 0.0, err
	}

	usageStr := strings.TrimSpace(string(output))
	if usageStr == "" {
		return 0.0, fmt.Errorf("no usage data returned")
	}

	if usage, err := strconv.ParseFloat(usageStr, 64); err == nil {
		if usage > 100 {
			usage = 100
		}
		return usage, nil
	}

	return 0.0, fmt.Errorf("failed to parse GPU usage from Task Manager approach")
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

func getAlternativeGPUMemoryUsage() (uint64, error) {
	cmd := exec.Command("powershell", "-Command", `
		Add-Type -TypeDefinition '
			using System;
			using System.Runtime.InteropServices;
			
			public class DXGI {
				[DllImport("dxgi.dll")]
				public static extern int CreateDXGIFactory(ref Guid riid, out IntPtr ppFactory);
				
				[DllImport("dxgi.dll")]
				public static extern int CreateDXGIFactory1(ref Guid riid, out IntPtr ppFactory);
			}
		' -ReferencedAssemblies System.Runtime.InteropServices -ErrorAction SilentlyContinue;
		
		# Fallback to WMI approach with better error handling
		try {
			$processes = Get-WmiObject -Query "SELECT * FROM Win32_PerfRawData_GPUPerformanceCounters_GPUProcessMemory" -ErrorAction SilentlyContinue;
			if ($processes) {
				$totalMemory = 0;
				foreach ($process in $processes) {
					if ($process.DedicatedUsage -gt 0) {
						$totalMemory += $process.DedicatedUsage;
					}
				}
				if ($totalMemory -gt 0) {
					return $totalMemory;
				}
			}
		} catch {}
		
		# Try GPU adapter memory
		try {
			$adapter = Get-WmiObject -Query "SELECT * FROM Win32_PerfRawData_GPUPerformanceCounters_GPUAdapter" -ErrorAction SilentlyContinue;
			if ($adapter -and $adapter.DedicatedUsage -gt 0) {
				return $adapter.DedicatedUsage;
			}
		} catch {}
		
		return 0;
	`)

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	memStr := strings.TrimSpace(string(output))
	if memStr == "" || memStr == "0" {
		return 0, fmt.Errorf("no memory usage data available")
	}

	if memBytes, err := strconv.ParseUint(memStr, 10, 64); err == nil && memBytes > 0 {
		// Sanity check
		if memBytes > 100*1024*1024*1024 { // > 100GB is unreasonable
			return 0, fmt.Errorf("unreasonable memory usage value: %d GB", memBytes/(1024*1024*1024))
		}
		return memBytes, nil
	}

	return 0, fmt.Errorf("failed to parse alternative GPU memory usage")
}
