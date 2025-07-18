//go:build windows
// +build windows

package temperature

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func GetWindowsTemperatureInfo() (*TemperatureData, error) {
	tempData := &TemperatureData{
		Sensors: make([]TemperatureSensor, 0),
	}

	if err := getPowerShellTemperatures(tempData); err == nil && len(tempData.Sensors) > 0 {
		calculateTemperatureStats(tempData)
	} else {
		return &TemperatureData{}, err
	}

	return tempData, nil
}

func getPowerShellTemperatures(tempData *TemperatureData) error {
	if err := getPerformanceCounterTemperatures(tempData); err == nil && len(tempData.Sensors) > 0 {
		return nil
	}

	return getWMITemperatures(tempData)
}

func getPerformanceCounterTemperatures(tempData *TemperatureData) error {
	cmd := exec.Command("powershell", "-Command",
		`try {
    $counters = Get-Counter '\Thermal Zone Information(*)\High Precision Temperature' -ErrorAction Stop
    foreach ($counter in $counters.CounterSamples) {
        $temp = [Math]::Round($counter.CookedValue / 10.0 - 273.15, 1)
        if ($temp -gt 0 -and $temp -lt 200) {
            $instanceName = $counter.InstanceName
            Write-Output "THERMAL:${temp}"
        }
    }
} catch {
    Write-Output "ERROR:Performance counters not available"
}
`)

	output, err := cmd.Output()
	if err != nil {
		return tryAlternativeTemperatureCounters(tempData)
	}

	result := strings.TrimSpace(string(output))
	if result == "" || strings.HasPrefix(result, "ERROR:") {
		return tryAlternativeTemperatureCounters(tempData)
	}

	lines := strings.Split(result, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "ERROR:") {
			continue
		}

		if strings.HasPrefix(line, "THERMAL:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				tempStr := parts[1]
				if temp, err := strconv.ParseFloat(tempStr, 64); err == nil && temp > 0 {
					sensor := TemperatureSensor{
						SensorKey:   "System",
						Temperature: temp,
						High:        85.0,
						Critical:    100.0,
					}
					tempData.Sensors = append(tempData.Sensors, sensor)
				}
			}
		}
	}

	if len(tempData.Sensors) == 0 {
		return tryAlternativeTemperatureCounters(tempData)
	}

	return nil
}

func tryAlternativeTemperatureCounters(tempData *TemperatureData) error {
	// Try alternative counter paths that might be available
	counterPaths := []string{
		"\\Thermal Zone Information(*)\\High Precision Temperature",
		"\\Thermal Zone Information(*)\\Temperature",
		"\\Thermal Zone Information(_Total)\\High Precision Temperature",
		"\\Thermal Zone Information(_Total)\\Temperature",
		"\\Processor Information(*)\\% Processor Time", // This won't give temp but can verify counters work
	}

	for _, path := range counterPaths {
		cmd := exec.Command("powershell", "-Command",
			"try { "+
				"$counters = Get-Counter '"+path+"' -ErrorAction Stop; "+
				"foreach ($counter in $counters.CounterSamples) { "+
				"if ($counter.Path -like '*Temperature*') { "+
				"$temp = [Math]::Round($counter.CookedValue / 10.0 - 273.15, 1); "+
				"if ($temp -gt 0 -and $temp -lt 200) { "+
				"$instanceName = if ($counter.InstanceName) { $counter.InstanceName } else { 'Unknown' }; "+
				"Write-Output \"THERMAL:$instanceName:$temp\" "+
				"} "+
				"} "+
				"} "+
				"} catch { "+
				"Write-Output \"ERROR:Counter path not available: "+path+"\" "+
				"}")

		output, err := cmd.Output()
		if err != nil {
			continue
		}

		result := strings.TrimSpace(string(output))
		if result == "" || strings.HasPrefix(result, "ERROR:") {
			continue
		}

		lines := strings.Split(result, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "ERROR:") {
				continue
			}

			if strings.HasPrefix(line, "THERMAL:") {
				parts := strings.Split(line, ":")
				if len(parts) == 3 {
					instanceName := parts[1]
					tempStr := parts[2]
					if temp, err := strconv.ParseFloat(tempStr, 64); err == nil && temp > 0 {
						sensor := TemperatureSensor{
							SensorKey:   fmt.Sprintf("ThermalZone_%s", instanceName),
							Temperature: temp,
							High:        85.0,
							Critical:    100.0,
						}
						tempData.Sensors = append(tempData.Sensors, sensor)
					}
				}
			}
		}

		if len(tempData.Sensors) > 0 {
			return nil
		}
	}

	return fmt.Errorf("no valid temperature sensors found with any counter path")
}

func getWMITemperatures(tempData *TemperatureData) error {
	cmd := exec.Command("powershell", "-Command",
		"try { "+
			"# Try multiple WMI approaches "+
			"$temps = @(); "+
			"$thermalZones = Get-WmiObject -Namespace 'root\\wmi' -Class 'MSAcpi_ThermalZoneTemperature' -ErrorAction SilentlyContinue | "+
			"ForEach-Object { [Math]::Round(($_.CurrentTemperature / 10.0 - 273.15), 1) }; "+
			"$thermalZones | ForEach-Object { if ($_ -gt 0 -and $_ -lt 200) { Write-Output \"THERMAL:$_\" } }; "+
			"# Try OpenHardwareMonitor WMI "+
			"$ohm = Get-WmiObject -Namespace 'root\\OpenHardwareMonitor' -Class 'Sensor' -ErrorAction SilentlyContinue | "+
			"Where-Object { $_.SensorType -eq 'Temperature' -and $_.Value -gt 0 -and $_.Value -lt 200 }; "+
			"$ohm | ForEach-Object { Write-Output \"OHM:$($_.Name):$($_.Value)\" }; "+
			"# Try LibreHardwareMonitor WMI "+
			"$lhm = Get-WmiObject -Namespace 'root\\LibreHardwareMonitor' -Class 'Sensor' -ErrorAction SilentlyContinue | "+
			"Where-Object { $_.SensorType -eq 'Temperature' -and $_.Value -gt 0 -and $_.Value -lt 200 }; "+
			"$lhm | ForEach-Object { Write-Output \"LHM:$($_.Name):$($_.Value)\" }; "+
			"# Try getting CPU temperature from WMI "+
			"try { "+
			"$cpuTemp = Get-WmiObject -Namespace 'root\\wmi' -Class 'MSAcpi_ThermalZoneTemperature' -ErrorAction Stop | "+
			"Select-Object -First 1 | ForEach-Object { [Math]::Round(($_.CurrentTemperature / 10.0 - 273.15), 1) }; "+
			"if ($cpuTemp -gt 0 -and $cpuTemp -lt 200) { Write-Output \"CPU:$cpuTemp\" } "+
			"} catch { } "+
			"} catch { "+
			"Write-Output \"ERROR:WMI access denied or unavailable\" "+
			"}")

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("PowerShell WMI query failed: %v", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "" || strings.HasPrefix(result, "ERROR:") {
		return fmt.Errorf("no temperature data from WMI: %s", result)
	}

	lines := strings.Split(result, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "ERROR:") {
			continue
		}

		if strings.HasPrefix(line, "CPU:") {
			tempStr := strings.TrimPrefix(line, "CPU:")
			if temp, err := strconv.ParseFloat(tempStr, 64); err == nil && temp > 0 {
				sensor := TemperatureSensor{
					SensorKey:   "CPU",
					Temperature: temp,
					High:        85.0,
					Critical:    100.0,
				}
				tempData.Sensors = append(tempData.Sensors, sensor)
			}
		} else if strings.HasPrefix(line, "THERMAL:") {
			tempStr := strings.TrimPrefix(line, "THERMAL:")
			if temp, err := strconv.ParseFloat(tempStr, 64); err == nil && temp > 0 {
				sensor := TemperatureSensor{
					SensorKey:   fmt.Sprintf("ThermalZone%d", len(tempData.Sensors)),
					Temperature: temp,
					High:        85.0,
					Critical:    100.0,
				}
				tempData.Sensors = append(tempData.Sensors, sensor)
			}
		} else if strings.HasPrefix(line, "OHM:") || strings.HasPrefix(line, "LHM:") {
			// OpenHardwareMonitor or LibreHardwareMonitor sensor
			parts := strings.Split(line, ":")
			if len(parts) == 3 {
				sensorType := parts[0]
				sensorName := parts[1]
				tempStr := parts[2]
				if temp, err := strconv.ParseFloat(tempStr, 64); err == nil && temp > 0 {
					sensor := TemperatureSensor{
						SensorKey:   fmt.Sprintf("%s_%s", sensorType, sensorName),
						Temperature: temp,
						High:        85.0,
						Critical:    100.0,
					}
					tempData.Sensors = append(tempData.Sensors, sensor)
				}
			}
		}
	}

	if len(tempData.Sensors) == 0 {
		return fmt.Errorf("no valid temperature sensors found")
	}

	return nil
}

func calculateTemperatureStats(tempData *TemperatureData) {
	if len(tempData.Sensors) == 0 {
		return
	}

	var totalTemp float64
	var maxTemp float64
	var cpuTemp float64
	var gpuTemp float64

	for _, sensor := range tempData.Sensors {
		if sensor.Temperature > maxTemp {
			maxTemp = sensor.Temperature
		}

		totalTemp += sensor.Temperature

		sensorKey := strings.ToLower(sensor.SensorKey)
		if strings.Contains(sensorKey, "cpu") || strings.Contains(sensorKey, "core") || strings.Contains(sensorKey, "processor") {
			if cpuTemp == 0 || sensor.Temperature > cpuTemp {
				cpuTemp = sensor.Temperature
			}
		}

		if strings.Contains(sensorKey, "gpu") || strings.Contains(sensorKey, "graphics") || strings.Contains(sensorKey, "video") {
			if gpuTemp == 0 || sensor.Temperature > gpuTemp {
				gpuTemp = sensor.Temperature
			}
		}
	}

	tempData.CPUTemp = cpuTemp
	tempData.GPUTemp = gpuTemp
	tempData.MaxTemp = maxTemp
	tempData.AvgTemp = totalTemp / float64(len(tempData.Sensors))
}
