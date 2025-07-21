package temperature

import (
	"fmt"
	"runtime"
	"strings"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/shirou/gopsutil/host"
)

type TemperatureSensor struct {
	SensorKey   string  `json:"sensor_key"`
	Temperature float64 `json:"temperature"`
	High        float64 `json:"high"`
	Critical    float64 `json:"critical"`
}

type TemperatureData struct {
	CPUTemp float64             `json:"cpu_temp"`
	GPUTemp float64             `json:"gpu_temp"`
	Sensors []TemperatureSensor `json:"sensors"`
	MaxTemp float64             `json:"max_temp"`
	AvgTemp float64             `json:"avg_temp"`
}

func GetTemperatures() (*TemperatureData, error) {
	switch runtime.GOOS {
	case "windows":
		return GetWindowsTemperatureInfo()
	case "linux":
		return GetLinuxTemperatureInfo()
	case "darwin":
		return GetDarwinTemperatureInfo()
	default:
		return getGopsutilTemperatures()
	}
}

func getGopsutilTemperatures() (*TemperatureData, error) {
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		return nil, err
	}

	tempData := &TemperatureData{
		Sensors: make([]TemperatureSensor, 0),
	}

	var totalTemp float64
	var count int

	for _, sensor := range sensors {
		if sensor.Temperature > 0 {
			tempSensor := TemperatureSensor{
				SensorKey:   sensor.SensorKey,
				Temperature: sensor.Temperature,
				High:        0,
				Critical:    0,
			}
			tempData.Sensors = append(tempData.Sensors, tempSensor)

			if contains(sensor.SensorKey, "cpu", "core", "processor") {
				if tempData.CPUTemp == 0 || sensor.Temperature > tempData.CPUTemp {
					tempData.CPUTemp = sensor.Temperature
				}
			}

			if contains(sensor.SensorKey, "gpu", "graphics", "nvidia", "amd", "radeon") {
				if tempData.GPUTemp == 0 || sensor.Temperature > tempData.GPUTemp {
					tempData.GPUTemp = sensor.Temperature
				}
			}

			totalTemp += sensor.Temperature
			count++

			if sensor.Temperature > tempData.MaxTemp {
				tempData.MaxTemp = sensor.Temperature
			}
		}
	}

	if count > 0 {
		tempData.AvgTemp = totalTemp / float64(count)
	}

	return tempData, nil
}

func contains(str string, substrings ...string) bool {
	for _, substr := range substrings {
		if len(str) >= len(substr) {
			for i := 0; i <= len(str)-len(substr); i++ {
				if str[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

func UpdateTemperatures(d *utils.Dashboard) {
	if d.TemperatureWidget == nil {
		return
	}

	tempData, err := GetTemperatures()
	if err != nil {
		d.TemperatureWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
			message := "Temperature monitoring unavailable on this system"

			lines := wrapText(message, w-6)

			currentY := y + 1
			for _, line := range lines {
				if currentY >= y+h-1 {
					break
				}
				currentY = utils.SafePrintText(screen, line, x+3, currentY, w-6, h-(currentY-y), tcell.ColorRed)
			}

			return x, y, w, h
		})
		return
	}

	d.TemperatureData = map[string]interface{}{
		"cpu_temp": tempData.CPUTemp,
		"gpu_temp": tempData.GPUTemp,
		"max_temp": tempData.MaxTemp,
		"avg_temp": tempData.AvgTemp,
		"sensors":  tempData.Sensors,
	}
	d.TemperatureWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		currentY := y + 1

		if tempData.CPUTemp > 0 {
			color := getTemperatureColor(tempData.CPUTemp)
			foregroundColor := utils.GetColorFromName(d.Theme.Layout.Temperature.ForegroundColor)
			text := fmt.Sprintf("CPU: [%s]%.1f°C[-]", color, tempData.CPUTemp)
			currentY = utils.SafePrintText(screen, text, x+2, currentY, w-2, h-(currentY-y), foregroundColor)
		}

		if tempData.GPUTemp > 0 {
			color := getTemperatureColor(tempData.GPUTemp)
			foregroundColor := utils.GetColorFromName(d.Theme.Layout.Temperature.ForegroundColor)
			text := fmt.Sprintf("GPU: [%s]%.1f°C[-]", color, tempData.GPUTemp)
			currentY = utils.SafePrintText(screen, text, x+2, currentY, w-2, h-(currentY-y), foregroundColor)
		}

		if tempData.MaxTemp > 0 {
			color := getTemperatureColor(tempData.MaxTemp)
			foregroundColor := utils.GetColorFromName(d.Theme.Layout.Temperature.ForegroundColor)
			text := fmt.Sprintf("Max: [%s]%.1f°C[-]", color, tempData.MaxTemp)
			currentY = utils.SafePrintText(screen, text, x+2, currentY, w-2, h-(currentY-y), foregroundColor)
		}

		if tempData.AvgTemp > 0 {
			color := getTemperatureColor(tempData.AvgTemp)
			foregroundColor := utils.GetColorFromName(d.Theme.Layout.Temperature.ForegroundColor)
			text := fmt.Sprintf("Avg: [%s]%.1f°C[-]", color, tempData.AvgTemp)
			currentY = utils.SafePrintText(screen, text, x+2, currentY, w-2, h-(currentY-y), foregroundColor)
		}

		maxSensors := (h - 6)
		if len(tempData.Sensors) > 0 && maxSensors > 0 {
			currentY = utils.SafePrintText(screen, "Sensors:", x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Temperature.ForegroundColor))

			for i, sensor := range tempData.Sensors {
				if i >= maxSensors {
					break
				}

				color := getTemperatureColor(sensor.Temperature)
				sensorName := truncateString(sensor.SensorKey, 8)
				text := fmt.Sprintf("%s: [%s]%.1f°C[-]", sensorName, color, sensor.Temperature)
				currentY = utils.SafePrintText(screen, text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Temperature.ForegroundColor))
			}
		}

		return x, y, w, h
	})
}

func getTemperatureColor(temp float64) string {
	if temp >= 85.0 {
		return "red"
	} else if temp >= 75.0 {
		return "orange"
	} else if temp >= 65.0 {
		return "yellow"
	} else {
		return "green"
	}
}

func truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen-3] + "..."
}

func GetTemperatureFormattedInfo() string {
	tempData, err := GetTemperatures()
	if err != nil {
		return fmt.Sprintf("Temperature monitoring: Error - %v", err)
	}

	info := "System Temperature Monitoring\n\n"

	if tempData.CPUTemp > 0 {
		info += fmt.Sprintf("CPU Temperature: %.1f°C\n", tempData.CPUTemp)
	}

	if tempData.GPUTemp > 0 {
		info += fmt.Sprintf("GPU Temperature: %.1f°C\n", tempData.GPUTemp)
	}

	if tempData.MaxTemp > 0 {
		info += fmt.Sprintf("Maximum Temperature: %.1f°C\n", tempData.MaxTemp)
	}

	if tempData.AvgTemp > 0 {
		info += fmt.Sprintf("Average Temperature: %.1f°C\n", tempData.AvgTemp)
	}

	if len(tempData.Sensors) > 0 {
		info += "\nAll Sensors:\n"
		for _, sensor := range tempData.Sensors {
			info += fmt.Sprintf("• %s: %.1f°C", sensor.SensorKey, sensor.Temperature)
			info += "\n"
		}
	}

	return info
}

func wrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	words := strings.Split(text, " ")
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine) == 0 {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= maxWidth {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}
