//go:build linux
// +build linux

package temperature

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/host"
)

const (
	thermalPath = "/sys/class/thermal"
	hwmonPath   = "/sys/class/hwmon"
)

func GetLinuxTemperatureInfo() (*TemperatureData, error) {
	tempData := &TemperatureData{
		Sensors: make([]TemperatureSensor, 0),
	}

	var errors []string

	if err := getLinuxThermalZones(tempData); err != nil {
		errors = append(errors, fmt.Sprintf("thermal zones: %v", err))
	}

	if err := getLinuxHWMONSensors(tempData); err != nil {
		errors = append(errors, fmt.Sprintf("hwmon sensors: %v", err))
	}

	if err := getLinuxLMSensors(tempData); err != nil {
		errors = append(errors, fmt.Sprintf("lm-sensors: %v", err))
	}

	if len(tempData.Sensors) > 0 {
		calculateLinuxTemperatureStats(tempData)
		return tempData, nil
	}

	if gopsutilData, err := getLinuxGopsutilTemperatures(); err == nil {
		return gopsutilData, nil
	} else {
		errors = append(errors, fmt.Sprintf("gopsutil: %v", err))
	}

	if len(errors) > 0 {
		return &TemperatureData{
			Sensors: make([]TemperatureSensor, 0),
		}, fmt.Errorf("no temperature sensors found - tried: %s", strings.Join(errors, "; "))
	}

	return &TemperatureData{
		Sensors: make([]TemperatureSensor, 0),
	}, fmt.Errorf("no temperature sensors available")
}

func getLinuxThermalZones(tempData *TemperatureData) error {
	entries, err := os.ReadDir(thermalPath)
	if err != nil {
		return fmt.Errorf("cannot read thermal directory: %v", err)
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "thermal_zone") {
			zonePath := filepath.Join(thermalPath, entry.Name())

			tempFile := filepath.Join(zonePath, "temp")
			if tempBytes, err := os.ReadFile(tempFile); err == nil {
				if temp, err := strconv.ParseFloat(strings.TrimSpace(string(tempBytes)), 64); err == nil {
					celsius := temp / 1000.0

					if celsius > 0 && celsius < 200 {
						sensor := TemperatureSensor{
							SensorKey:   fmt.Sprintf("thermal_zone_%s", entry.Name()),
							Temperature: celsius,
							High:        85.0,
							Critical:    100.0,
						}

						if typeData, err := os.ReadFile(filepath.Join(zonePath, "type")); err == nil {
							sensorType := strings.TrimSpace(string(typeData))
							sensor.SensorKey = fmt.Sprintf("thermal_%s", sensorType)
						}

						if tripData, err := os.ReadFile(filepath.Join(zonePath, "trip_point_0_temp")); err == nil {
							if tripTemp, err := strconv.ParseFloat(strings.TrimSpace(string(tripData)), 64); err == nil {
								sensor.High = tripTemp / 1000.0
							}
						}

						if criticalData, err := os.ReadFile(filepath.Join(zonePath, "trip_point_1_temp")); err == nil {
							if criticalTemp, err := strconv.ParseFloat(strings.TrimSpace(string(criticalData)), 64); err == nil {
								sensor.Critical = criticalTemp / 1000.0
							}
						}

						tempData.Sensors = append(tempData.Sensors, sensor)
					}
				}
			}
		}
	}

	return nil
}

func getLinuxHWMONSensors(tempData *TemperatureData) error {
	entries, err := os.ReadDir(hwmonPath)
	if err != nil {
		return fmt.Errorf("cannot read hwmon directory: %v", err)
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "hwmon") {
			hwmonDir := filepath.Join(hwmonPath, entry.Name())

			var deviceName string
			if nameData, err := os.ReadFile(filepath.Join(hwmonDir, "name")); err == nil {
				deviceName = strings.TrimSpace(string(nameData))
			}

			if hwmonEntries, err := os.ReadDir(hwmonDir); err == nil {
				for _, hwmonEntry := range hwmonEntries {
					if strings.HasPrefix(hwmonEntry.Name(), "temp") && strings.HasSuffix(hwmonEntry.Name(), "_input") {
						tempFile := filepath.Join(hwmonDir, hwmonEntry.Name())
						if tempBytes, err := os.ReadFile(tempFile); err == nil {
							if temp, err := strconv.ParseFloat(strings.TrimSpace(string(tempBytes)), 64); err == nil {
								celsius := temp / 1000.0

								if celsius > 0 && celsius < 200 {
									sensorKey := fmt.Sprintf("hwmon_%s_%s", deviceName, hwmonEntry.Name())
									if deviceName == "" {
										sensorKey = fmt.Sprintf("hwmon_%s_%s", entry.Name(), hwmonEntry.Name())
									}

									sensor := TemperatureSensor{
										SensorKey:   sensorKey,
										Temperature: celsius,
										High:        85.0,
										Critical:    100.0,
									}

									maxFile := strings.Replace(hwmonEntry.Name(), "_input", "_max", 1)
									if maxData, err := os.ReadFile(filepath.Join(hwmonDir, maxFile)); err == nil {
										if maxTemp, err := strconv.ParseFloat(strings.TrimSpace(string(maxData)), 64); err == nil {
											sensor.High = maxTemp / 1000.0
										}
									}

									critFile := strings.Replace(hwmonEntry.Name(), "_input", "_crit", 1)
									if critData, err := os.ReadFile(filepath.Join(hwmonDir, critFile)); err == nil {
										if critTemp, err := strconv.ParseFloat(strings.TrimSpace(string(critData)), 64); err == nil {
											sensor.Critical = critTemp / 1000.0
										}
									}

									tempData.Sensors = append(tempData.Sensors, sensor)
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func getLinuxLMSensors(tempData *TemperatureData) error {
	cmd := exec.Command("sensors", "-A", "-u")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("sensors command failed: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	var currentAdapter string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, " ") && strings.Contains(line, ":") {
			currentAdapter = strings.TrimSuffix(line, ":")
			continue
		}

		if strings.Contains(line, "_input:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				tempName := parts[0]
				tempName = strings.TrimSuffix(tempName, "_input:")

				if temp, err := strconv.ParseFloat(parts[1], 64); err == nil {
					if temp > 0 && temp < 200 {
						sensor := TemperatureSensor{
							SensorKey:   fmt.Sprintf("lm_sensors_%s_%s", currentAdapter, tempName),
							Temperature: temp,
							High:        85.0,
							Critical:    100.0,
						}
						tempData.Sensors = append(tempData.Sensors, sensor)
					}
				}
			}
		}
	}

	return nil
}

func getLinuxGopsutilTemperatures() (*TemperatureData, error) {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return nil, fmt.Errorf("gopsutil failed: %v", err)
	}

	tempData := &TemperatureData{
		Sensors: make([]TemperatureSensor, 0, len(temps)),
	}

	for _, t := range temps {
		if t.Temperature > 0 && t.Temperature < 200 {
			sensor := TemperatureSensor{
				SensorKey:   t.SensorKey,
				Temperature: t.Temperature,
				High:        85.0,
				Critical:    100.0,
			}
			tempData.Sensors = append(tempData.Sensors, sensor)
		}
	}

	if len(tempData.Sensors) > 0 {
		calculateLinuxTemperatureStats(tempData)
		return tempData, nil
	}

	return nil, fmt.Errorf("no temperature sensors found via gopsutil")
}

func calculateLinuxTemperatureStats(tempData *TemperatureData) {
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
		if strings.Contains(sensorKey, "cpu") || strings.Contains(sensorKey, "core") || strings.Contains(sensorKey, "processor") || strings.Contains(sensorKey, "k10temp") || strings.Contains(sensorKey, "coretemp") {
			if cpuTemp == 0 || sensor.Temperature > cpuTemp {
				cpuTemp = sensor.Temperature
			}
		}

		if strings.Contains(sensorKey, "gpu") || strings.Contains(sensorKey, "graphics") || strings.Contains(sensorKey, "nvidia") || strings.Contains(sensorKey, "radeon") || strings.Contains(sensorKey, "amdgpu") {
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
