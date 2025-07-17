//go:build windows
// +build windows

package temperature

import (
	"fmt"
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
	return fmt.Errorf("Temperature monitoring unavailable!")
	/*
		cmd := exec.Command("powershell", "-Command",
			"$temps = Get-WmiObject -Namespace 'root\\wmi' -Class 'MSAcpi_ThermalZoneTemperature' | "+
				"ForEach-Object { [Math]::Round(($_.CurrentTemperature / 10.0 - 273.15), 1) }; "+
				"$cpuTemp = Get-WmiObject -Namespace 'root\\OpenHardwareMonitor' -Class 'Sensor' | "+
				"Where-Object { $_.Name -like '*CPU*' -and $_.SensorType -eq 'Temperature' } | "+
				"Select-Object -First 1 -ExpandProperty Value; "+
				"if ($cpuTemp) { Write-Output \"CPU:$cpuTemp\" }; "+
				"$temps | ForEach-Object { Write-Output \"THERMAL:$_\" }")

		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("PowerShell temperature query failed: %v", err)
		}

		result := strings.TrimSpace(string(output))
		if result == "" {
			return fmt.Errorf("no temperature data found")
		}

		lines := strings.Split(result, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
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
			}
		}

		return nil
	*/
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
