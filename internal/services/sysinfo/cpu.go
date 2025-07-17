package sysinfo

import (
	"fmt"
	"strings"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/cpu"
)

func GetCpuInfo() []cpu.InfoStat {
	info, err := cpu.Info()
	if err != nil {
		return nil
	}

	return info
}

func GetCpuFormattedInfo() string {
	info, err := cpu.Info()
	if err != nil {
		return "Unknown CPU"
	}

	// Get current CPU usage
	percents, err := cpu.Percent(0, true)
	if err != nil {
		percents = []float64{}
	}

	var totalUsage float64
	if len(percents) > 0 {
		var sum float64
		for _, p := range percents {
			sum += p
		}
		totalUsage = sum / float64(len(percents))
	}

	var output string

	output += "=== CPU Information ===\n"
	for i, cpu := range info {
		if len(info) > 1 {
			output += fmt.Sprintf("\n--- CPU %d ---\n", i)
		}
		output += fmt.Sprintf("Vendor: %s\n", cpu.VendorID)
		output += fmt.Sprintf("Model: %s\n", cpu.ModelName)
		output += fmt.Sprintf("Family: %s\n", cpu.Family)
		output += fmt.Sprintf("Model ID: %s\n", cpu.Model)
		output += fmt.Sprintf("Stepping: %d\n", cpu.Stepping)
		output += fmt.Sprintf("Cores: %d\n", cpu.Cores)
		output += fmt.Sprintf("Base Frequency: %.0f MHz\n", cpu.Mhz)
		output += fmt.Sprintf("Cache Size: %d KB\n", cpu.CacheSize)
		output += fmt.Sprintf("CPU Index: %d\n", cpu.CPU)

		if len(cpu.Flags) > 0 {
			output += fmt.Sprintf("Features: %s\n", strings.Join(cpu.Flags[:min(10, len(cpu.Flags))], ", "))
			if len(cpu.Flags) > 10 {
				output += fmt.Sprintf("... and %d more features\n", len(cpu.Flags)-10)
			}
		}
	}

	output += "\n=== Current CPU Usage ===\n"
	output += fmt.Sprintf("Average CPU Usage: %.1f%%\n", totalUsage)

	output += "\n=== CPU Health Status ===\n"
	if totalUsage < 30 {
		output += "• CPU Load: Low - system is running smoothly\n"
	} else if totalUsage < 60 {
		output += "• CPU Load: Moderate - normal usage levels\n"
	} else if totalUsage < 80 {
		output += "• CPU Load: High - system is working hard\n"
	} else {
		output += "• CPU Load: Very High - system may be under stress\n"
	}

	return output
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func UpdateCPU(d *utils.Dashboard) {
	if d.CpuWidget == nil {
		return
	}

	percents, err := cpu.Percent(0, true)
	if err != nil {
		return
	}

	var totalUsage float64
	if len(percents) > 0 {
		var sum float64
		for _, p := range percents {
			sum += p
		}
		totalUsage = sum / float64(len(percents))
	}

	d.CpuData = percents
	d.CpuWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		barCount := int(totalUsage / 12)
		color := d.Theme.CPU.BarLow
		if totalUsage > 80 {
			color = d.Theme.CPU.BarHigh
		}

		totalText := fmt.Sprintf("Total CPU usage: %s %.0f%%", utils.BarColor(utils.BAR, barCount, color), totalUsage)
		tview.Print(screen, totalText, x+2, y+1, w-2, h-1, utils.GetColorFromName(d.Theme.Layout.CPU.ForegroundColor))

		currentY := y + 2

		colWidth := 26
		maxCols := (w - 4) / colWidth
		if maxCols < 1 {
			maxCols = 1
		}

		for i, p := range d.CpuData {
			barCount := int(p / 10)
			color := d.Theme.CPU.BarLow
			if p > 80 {
				color = d.Theme.CPU.BarHigh
			}
			coreText := fmt.Sprintf("Core %d: %s %.0f%%", i, utils.BarColor(utils.BAR, barCount, color), p)

			col := i % maxCols
			row := i / maxCols

			xPos := x + 2 + (col * colWidth)
			yPos := currentY + 1 + row

			if yPos < h {
				tview.Print(screen, coreText, xPos, yPos, colWidth-2, h-yPos+2, utils.GetColorFromName(d.Theme.Layout.CPU.ForegroundColor))
			}
		}
		return x, y, w, h
	})
}

func GetCpuName() string {
	data := GetCpuInfo()
	if len(data) == 0 {
		return "Unknown CPU"
	}

	fullName := data[0].ModelName
	return extractCpuBrandSeries(fullName)
}

func extractCpuBrandSeries(fullName string) string {

	if fullName == "" {
		return "Unknown CPU"
	}

	if containsAny(fullName, []string{"Intel", "INTEL"}) {
		if containsAny(fullName, []string{"Core", "CORE"}) {
			return "Intel Core"
		} else if containsAny(fullName, []string{"Pentium", "PENTIUM"}) {
			return "Intel Pentium"
		} else if containsAny(fullName, []string{"Celeron", "CELERON"}) {
			return "Intel Celeron"
		} else if containsAny(fullName, []string{"Xeon", "XEON"}) {
			return "Intel Xeon"
		} else if containsAny(fullName, []string{"Atom", "ATOM"}) {
			return "Intel Atom"
		} else {
			return "Intel"
		}
	}

	if containsAny(fullName, []string{"AMD", "amd"}) {
		if containsAny(fullName, []string{"Ryzen", "RYZEN"}) {
			return "AMD Ryzen"
		} else if containsAny(fullName, []string{"FX", "fx"}) {
			return "AMD FX"
		} else if containsAny(fullName, []string{"A-Series", "A Series", "A4", "A6", "A8", "A10", "A12"}) {
			return "AMD A-Series"
		} else if containsAny(fullName, []string{"Athlon", "ATHLON"}) {
			return "AMD Athlon"
		} else if containsAny(fullName, []string{"Phenom", "PHENOM"}) {
			return "AMD Phenom"
		} else if containsAny(fullName, []string{"Opteron", "OPTERON"}) {
			return "AMD Opteron"
		} else if containsAny(fullName, []string{"EPYC", "epyc"}) {
			return "AMD EPYC"
		} else {
			return "AMD"
		}
	}

	return fullName
}

func containsAny(str string, substrings []string) bool {
	for _, substring := range substrings {
		if strings.Contains(str, substring) {
			return true
		}
	}
	return false
}
