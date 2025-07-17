//go:build darwin
// +build darwin

package temperature

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func GetDarwinTemperatureInfo() (*TemperatureData, error) {
	tempData := &TemperatureData{
		Sensors: make([]TemperatureSensor, 0),
	}

	if err := getDarwinPowermetricsTemperatures(tempData); err == nil && len(tempData.Sensors) > 0 {
		calculateDarwinTemperatureStats(tempData)
		return tempData, nil
	}

	if err := getDarwinIstatsTemperatures(tempData); err == nil && len(tempData.Sensors) > 0 {
		calculateDarwinTemperatureStats(tempData)
		return tempData, nil
	}
	if err := getDarwinSysctlTemperatures(tempData); err == nil && len(tempData.Sensors) > 0 {
		calculateDarwinTemperatureStats(tempData)
		return tempData, nil
	}

	if err := getDarwinIOKitTemperatures(tempData); err == nil && len(tempData.Sensors) > 0 {
		calculateDarwinTemperatureStats(tempData)
		return tempData, nil
	}

	return tempData, fmt.Errorf("no temperature sensors available")
}

func getDarwinPowermetricsTemperatures(tempData *TemperatureData) error {
	cmd := exec.Command("powermetrics", "-n", "1", "-s", "thermal")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("powermetrics command failed: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	tempRegex := regexp.MustCompile(`([^:]+):\s*([0-9.]+)°C`)

	for _, line := range lines {
		matches := tempRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			sensorName := strings.TrimSpace(matches[1])
			if temp, err := strconv.ParseFloat(matches[2], 64); err == nil {
				if temp > 0 && temp < 200 {
					sensor := TemperatureSensor{
						SensorKey:   fmt.Sprintf("powermetrics_%s", sensorName),
						Temperature: temp,
					}
					tempData.Sensors = append(tempData.Sensors, sensor)
				}
			}
		}
	}

	return nil
}

func getDarwinIstatsTemperatures(tempData *TemperatureData) error {
	cmd := exec.Command("istats", "temp")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("istats command failed: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	tempRegex := regexp.MustCompile(`([^:]+):\s*([0-9.]+)°C`)

	for _, line := range lines {
		matches := tempRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			sensorName := strings.TrimSpace(matches[1])
			if temp, err := strconv.ParseFloat(matches[2], 64); err == nil {
				if temp > 0 && temp < 200 {
					sensor := TemperatureSensor{
						SensorKey:   fmt.Sprintf("istats_%s", sensorName),
						Temperature: temp,
					}
					tempData.Sensors = append(tempData.Sensors, sensor)
				}
			}
		}
	}

	return nil
}

func getDarwinSysctlTemperatures(tempData *TemperatureData) error {
	cmd := exec.Command("sysctl", "-n", "machdep.xcpm.cpu_thermal_state")
	if output, err := cmd.Output(); err == nil {
		if temp, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64); err == nil {
			if temp > 0 && temp < 200 {
				sensor := TemperatureSensor{
					SensorKey:   "sysctl_cpu_thermal_state",
					Temperature: temp,
				}
				tempData.Sensors = append(tempData.Sensors, sensor)
			}
		}
	}

	cmd = exec.Command("sysctl", "-n", "machdep.xcpm.cpu_thermal_pressure")
	if output, err := cmd.Output(); err == nil {
		if pressure, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64); err == nil {
			if pressure > 0 {
				estimatedTemp := 40.0 + (pressure * 60.0 / 100.0)
				sensor := TemperatureSensor{
					SensorKey:   "sysctl_cpu_thermal_pressure",
					Temperature: estimatedTemp,
				}
				tempData.Sensors = append(tempData.Sensors, sensor)
			}
		}
	}

	return nil
}

func getDarwinIOKitTemperatures(tempData *TemperatureData) error {
	cmd := exec.Command("ioreg", "-c", "IOPMrootDomain", "-d", "1")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("ioreg command failed: %v", err)
	}

	outputStr := string(output)

	tempRegex := regexp.MustCompile(`"([^"]*[Tt]emp[^"]*)"[^=]*=\s*([0-9]+)`)
	matches := tempRegex.FindAllStringSubmatch(outputStr, -1)

	for _, match := range matches {
		if len(match) == 3 {
			sensorName := match[1]
			if rawTemp, err := strconv.ParseFloat(match[2], 64); err == nil {
				var celsius float64

				if rawTemp > 1000 {
					celsius = rawTemp / 1000.0
				} else if rawTemp > 100 {
					celsius = rawTemp / 100.0
				} else {
					celsius = rawTemp
				}

				if celsius > 0 && celsius < 200 {
					sensor := TemperatureSensor{
						SensorKey:   fmt.Sprintf("iokit_%s", sensorName),
						Temperature: celsius,
					}
					tempData.Sensors = append(tempData.Sensors, sensor)
				}
			}
		}
	}

	return nil
}

func calculateDarwinTemperatureStats(tempData *TemperatureData) {
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
		if strings.Contains(sensorKey, "cpu") || strings.Contains(sensorKey, "core") || strings.Contains(sensorKey, "processor") || strings.Contains(sensorKey, "thermal") {
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
