package gpu

import (
	"fmt"
	"strings"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
)

func UpdateGPU(d *utils.Dashboard) error {
	if d.GPUWidget == nil {
		return fmt.Errorf("GPU widget not initialized")
	}

	gpus, err := GetGPUInfo()
	if err != nil {
		d.GPUData = nil
		d.GPUWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
			utils.SafePrintText(screen, "GPU monitoring unavailable", x+1, y+1, w-2, h-(y+1-y), tcell.ColorRed)
			utils.SafePrintText(screen, fmt.Sprintf("Error: %s", err.Error()), x+1, y+2, w-2, h-(y+2-y), tcell.ColorRed)
			return x, y, w, h
		})
		return err
	}

	if len(gpus) == 0 {
		d.GPUData = nil
		d.GPUWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
			utils.SafePrintText(screen, "No GPUs detected", x+1, y+1, w-2, h-(y+1-y), tcell.ColorYellow)
			return x, y, w, h
		})
		return nil
	}

	var gpuDataSlice []interface{}
	for _, gpu := range gpus {
		gpuMap := map[string]interface{}{
			"name":         gpu.Name,
			"vendor":       gpu.Vendor,
			"memory_total": gpu.MemoryTotal,
			"memory_used":  gpu.MemoryUsed,
			"memory_free":  gpu.MemoryFree,
			"temperature":  gpu.Temperature,
			"usage":        gpu.Usage,
			"available":    gpu.Available,
		}
		gpuDataSlice = append(gpuDataSlice, gpuMap)
	}
	d.GPUData = gpuDataSlice

	d.GPUWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		currentY := y + 1

		for i, gpu := range gpus {
			if currentY >= y+h-1 {
				break
			}

			gpuTitle := fmt.Sprintf("%s (%s)", gpu.Name, gpu.Vendor)
			color := utils.GetColorFromName(d.Theme.Layout.GPU.ForegroundColor)
			if !gpu.Available {
				color = tcell.ColorRed
			}

			currentY = utils.SafePrintText(screen, gpuTitle, x, currentY, w-2, h-(currentY-y), color)

			if gpu.MemoryTotal > 0 {
				memoryPercent := 0.0
				if gpu.MemoryTotal > 0 {
					memoryPercent = float64(gpu.MemoryUsed) / float64(gpu.MemoryTotal) * 100
				}

				memoryText := fmt.Sprintf("Memory: %s / %s (%.1f%%)",
					formatMemorySize(gpu.MemoryUsed),
					formatMemorySize(gpu.MemoryTotal),
					memoryPercent)

				barWidth := w - 4
				if barWidth > 0 {
					bar := createUsageBar(memoryPercent, barWidth)
					currentY = utils.SafePrintText(screen, bar, x, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.GPU.ForegroundColor))
				}

				if currentY < y+h-1 {
					currentY = utils.SafePrintText(screen, memoryText, x, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.GPU.ForegroundColor))
				}
			}

			if gpu.Usage > 0 {
				usageText := fmt.Sprintf("GPU Usage: %.1f%%", gpu.Usage)
				bar := createUsageBar(gpu.Usage, w-4)

				if currentY < y+h-1 {
					currentY = utils.SafePrintText(screen, bar, x, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.GPU.ForegroundColor))
				}

				if currentY < y+h-1 {
					currentY = utils.SafePrintText(screen, usageText, x, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.GPU.ForegroundColor))
				}
			}

			if gpu.Temperature > 0 {
				tempText := fmt.Sprintf("Temperature: %.1f°C", gpu.Temperature)
				tempColor := tcell.ColorGreen
				if gpu.Temperature > 80 {
					tempColor = tcell.ColorRed
				} else if gpu.Temperature > 65 {
					tempColor = tcell.ColorYellow
				}

				if currentY < y+h-1 {
					currentY = utils.SafePrintText(screen, tempText, x, currentY, w-2, h-(currentY-y), tempColor)
				}
			}

			if gpu.PowerDraw > 0 {
				powerText := fmt.Sprintf("Power: %dW", gpu.PowerDraw)
				if gpu.PowerLimit > 0 {
					powerText = fmt.Sprintf("Power: %dW / %dW", gpu.PowerDraw, gpu.PowerLimit)
				}

				if currentY < y+h-1 {
					currentY = utils.SafePrintText(screen, powerText, x, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.GPU.ForegroundColor))
				}
			}

			if gpu.ClockSpeed > 0 {
				clockText := fmt.Sprintf("Clock: %d MHz", gpu.ClockSpeed)
				if currentY < y+h-1 {
					currentY = utils.SafePrintText(screen, clockText, x, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.GPU.ForegroundColor))
				}
			}

			if gpu.Driver != "" {
				driverText := fmt.Sprintf("Driver: %s", gpu.Driver)
				if currentY < y+h-1 {
					currentY = utils.SafePrintText(screen, driverText, x, currentY, w-2, h-(currentY-y), tcell.ColorGray)
				}
			}

			if i < len(gpus)-1 && currentY < y+h-1 {
				separator := strings.Repeat("-", w-2)
				currentY = utils.SafePrintText(screen, separator, x, currentY, w-2, h-(currentY-y), tcell.ColorGray)
			}
		}

		return x, y, w, h
	})

	return nil
}

func createUsageBar(percentage float64, width int) string {
	if width <= 0 {
		return ""
	}

	filled := int(percentage / 100.0 * float64(width))
	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return bar
}

func GetGPUTitle() string {
	gpus, err := GetGPUInfo()
	if err != nil {
		return "GPU - Error"
	}

	if len(gpus) == 0 {
		return "GPU - None detected"
	}

	if len(gpus) == 1 {
		return fmt.Sprintf("GPU - %s", gpus[0].Name)
	}

	return fmt.Sprintf("GPU - %d devices", len(gpus))
}

func GetGPUFormattedInfo() string {
	gpus, err := GetGPUInfo()
	if err != nil {
		return fmt.Sprintf("Error getting GPU information: %v", err)
	}

	if len(gpus) == 0 {
		return "No GPUs detected on this system."
	}

	var info strings.Builder

	for i, gpu := range gpus {
		if i > 0 {
			info.WriteString("\n")
		}

		info.WriteString(fmt.Sprintf("GPU %d: %s\n", i+1, gpu.Name))
		info.WriteString(fmt.Sprintf("Vendor: %s\n", gpu.Vendor))

		if gpu.Driver != "" {
			info.WriteString(fmt.Sprintf("Driver: %s\n", gpu.Driver))
		}

		if gpu.MemoryTotal > 0 {
			info.WriteString(fmt.Sprintf("Memory: %s", formatMemorySize(gpu.MemoryTotal)))
			if gpu.MemoryUsed > 0 {
				usedPercent := float64(gpu.MemoryUsed) / float64(gpu.MemoryTotal) * 100
				info.WriteString(fmt.Sprintf(" (%.1f%% used)", usedPercent))
			}
			info.WriteString("\n")
		}

		if gpu.Usage > 0 {
			info.WriteString(fmt.Sprintf("Usage: %.1f%%\n", gpu.Usage))
		}

		if gpu.Temperature > 0 {
			info.WriteString(fmt.Sprintf("Temperature: %.1f°C\n", gpu.Temperature))
		}

		if gpu.PowerDraw > 0 {
			info.WriteString(fmt.Sprintf("Power Draw: %dW", gpu.PowerDraw))
			if gpu.PowerLimit > 0 {
				info.WriteString(fmt.Sprintf(" / %dW", gpu.PowerLimit))
			}
			info.WriteString("\n")
		}

		if gpu.ClockSpeed > 0 {
			info.WriteString(fmt.Sprintf("Clock Speed: %d MHz\n", gpu.ClockSpeed))
		}

		if gpu.FanSpeed > 0 {
			info.WriteString(fmt.Sprintf("Fan Speed: %d RPM\n", gpu.FanSpeed))
		}

		info.WriteString(fmt.Sprintf("Available: %v\n", gpu.Available))
	}

	return info.String()
}
