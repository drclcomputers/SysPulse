package load

import (
	"fmt"
	"runtime"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LoadAverage struct {
	Load1  float64 `json:"load_1"`
	Load5  float64 `json:"load_5"`
	Load15 float64 `json:"load_15"`
}

func GetLoadAverage() (*LoadAverage, error) {
	switch runtime.GOOS {
	case "windows":
		return GetWindowsLoadAverage()
	case "linux":
		return GetLinuxLoadAverage()
	case "darwin":
		return GetDarwinLoadAverage()
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func UpdateLoadAverage(d *utils.Dashboard) {
	if d.LoadWidget == nil {
		return
	}

	loadAvg, err := GetLoadAverage()
	if err != nil {
		d.LoadWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
			utils.SafePrintText(screen, "Load average unavailable", x+3, y+1, w-2, h-(y+1-y), tcell.ColorRed)
			return x, y, w, h
		})
		return
	}

	d.LoadData = loadAvg
	d.LoadWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		cpuCount := runtime.NumCPU()

		load1Color := getLoadColor(loadAvg.Load1, cpuCount)
		load5Color := getLoadColor(loadAvg.Load5, cpuCount)
		load15Color := getLoadColor(loadAvg.Load15, cpuCount)

		headerText := fmt.Sprintf("Load Average (CPU cores: %d)", cpuCount)
		currentY := utils.SafePrintText(screen, headerText, x+2, y+1, w-2, h-1, utils.GetColorFromName(d.Theme.Layout.Load.ForegroundColor))

		load1Text := fmt.Sprintf("1 min:  [%s]%.2f[-]", load1Color, loadAvg.Load1)
		currentY = utils.SafePrintText(screen, load1Text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Load.ForegroundColor))

		load5Text := fmt.Sprintf("5 min:  [%s]%.2f[-]", load5Color, loadAvg.Load5)
		currentY = utils.SafePrintText(screen, load5Text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Load.ForegroundColor))

		load15Text := fmt.Sprintf("15 min: [%s]%.2f[-]", load15Color, loadAvg.Load15)
		currentY = utils.SafePrintText(screen, load15Text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Load.ForegroundColor))

		bar1 := createLoadBar(loadAvg.Load1, cpuCount, w/3)
		bar5 := createLoadBar(loadAvg.Load5, cpuCount, w/3)
		bar15 := createLoadBar(loadAvg.Load15, cpuCount, w/3)

		bar1Text := fmt.Sprintf("1m:  %s", bar1)
		tview.Print(screen, bar1Text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Load.ForegroundColor))
		currentY++

		bar5Text := fmt.Sprintf("5m:  %s", bar5)
		tview.Print(screen, bar5Text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Load.ForegroundColor))
		currentY++

		bar15Text := fmt.Sprintf("15m: %s", bar15)
		tview.Print(screen, bar15Text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Load.ForegroundColor))
		currentY++

		return x, y, w, h
	})
}

func getLoadColor(load float64, cpuCount int) string {
	ratio := load / float64(cpuCount)

	if ratio >= 1.0 {
		return "red"
	} else if ratio >= 0.7 {
		return "yellow"
	} else if ratio >= 0.5 {
		return "orange"
	} else {
		return "green"
	}
}

func createLoadBar(load float64, cpuCount int, width int) string {
	if width <= 0 {
		return ""
	}

	ratio := load / float64(cpuCount)
	if ratio > 1.0 {
		ratio = 1.0
	}

	filled := int(ratio * float64(width))
	if filled > width {
		filled = width
	}

	color := getLoadColor(load, cpuCount)
	bar := fmt.Sprintf("[%s]%s[-]", color, utils.RepeatString("█", filled))
	bar += utils.RepeatString("░", width-filled)

	return bar
}

func GetLoadFormattedInfo() string {
	loadAvg, err := GetLoadAverage()
	if err != nil {
		return fmt.Sprintf("Load Average: Error - %v", err)
	}

	cpuCount := runtime.NumCPU()

	var info string

	info += fmt.Sprintf("CPU Cores: %d\n\n", cpuCount)
	info += fmt.Sprintf("1 minute:  %.2f (%.1f%%)\n", loadAvg.Load1, (loadAvg.Load1/float64(cpuCount))*100)
	info += fmt.Sprintf("5 minutes: %.2f (%.1f%%)\n", loadAvg.Load5, (loadAvg.Load5/float64(cpuCount))*100)
	info += fmt.Sprintf("15 minutes: %.2f (%.1f%%)\n", loadAvg.Load15, (loadAvg.Load15/float64(cpuCount))*100)

	info += "\nLoad Interpretation:\n"
	info += "• < 0.7: Low load\n"
	info += "• 0.7-1.0: Moderate load\n"
	info += "• > 1.0: High load (system may be overloaded)\n"

	return info
}
