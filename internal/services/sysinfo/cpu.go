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

	for i, cpu := range info {
		output += fmt.Sprintf("\n--- CPU %d ---\n", i)
		output += fmt.Sprintf("Vendor: %s\n", cpu.VendorID)
		output += fmt.Sprintf("Model: %s\n", cpu.ModelName)
		output += fmt.Sprintf("Family: %s\n", cpu.Family)
		output += fmt.Sprintf("Model ID: %s\n", cpu.Model)
		output += fmt.Sprintf("Stepping: %d\n", cpu.Stepping)
		output += fmt.Sprintf("Cores: %d\n", cpu.Cores)
		output += fmt.Sprintf("Threads: %d\n", len(percents))
		output += fmt.Sprintf("Base Frequency: %.0f MHz\n", cpu.Mhz)
		output += fmt.Sprintf("Cache Size: %d KB\n", cpu.CacheSize)
		output += fmt.Sprintf("Architecture Support: %s\n", getArchitectureInfo(cpu.Flags))

		importantFlags := []string{"sse", "sse2", "sse3", "sse4_1", "sse4_2", "avx", "avx2", "avx512", "aes", "vmx", "hypervisor"}
		var relevantFlags []string

		for _, flag := range cpu.Flags {
			for _, important := range importantFlags {
				if strings.Contains(strings.ToLower(flag), important) {
					relevantFlags = append(relevantFlags, flag)
					break
				}
			}
		}

		if len(relevantFlags) > 0 {
			output += fmt.Sprintf("Key Features: %s\n", strings.Join(relevantFlags, ", "))
		}
		output += fmt.Sprintf("Total Features: %d\n", len(cpu.Flags))
	}

	output += fmt.Sprintf("--- Average CPU Usage: %.1f%%\n", totalUsage)

	return output
}

func getArchitectureInfo(flags []string) string {
	features := []string{}
	if containsFlag(flags, "avx512") {
		features = append(features, "AVX-512")
	} else if containsFlag(flags, "avx2") {
		features = append(features, "AVX2")
	} else if containsFlag(flags, "avx") {
		features = append(features, "AVX")
	}

	if containsFlag(flags, "aes") {
		features = append(features, "AES")
	}

	if containsFlag(flags, "vmx") || containsFlag(flags, "svm") {
		features = append(features, "Virtualization")
	}

	if len(features) == 0 {
		return "Basic"
	}
	return strings.Join(features, ", ")
}

func containsFlag(flags []string, target string) bool {
	for _, flag := range flags {
		if strings.Contains(strings.ToLower(flag), strings.ToLower(target)) {
			return true
		}
	}
	return false
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
		color := d.Theme.CPU.BarLow
		if totalUsage > 80 {
			color = d.Theme.CPU.BarHigh
		}

		totalText := fmt.Sprintf("Total: %s %.0f%%", utils.BarColor(utils.BAR, w/3, color), totalUsage)
		tview.Print(screen, totalText, x+2, y+1, w-2, h-1, utils.GetColorFromName(d.Theme.Layout.CPU.ForegroundColor))

		currentY := y + 2

		colWidth := 26
		maxCols := (w - 4) / colWidth
		if maxCols < 1 {
			maxCols = 1
		}

		for i, p := range d.CpuData {
			color := d.Theme.CPU.BarLow
			if p > 80 {
				color = d.Theme.CPU.BarHigh
			}
			coreText := fmt.Sprintf("Core %d: %s %.0f%%", i, utils.BarColor(utils.BAR, w/3, color), p)

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
