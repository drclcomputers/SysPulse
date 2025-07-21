package battery

import (
	"fmt"
	"runtime"
	"syspulse/internal/utils"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/shirou/gopsutil/host"
)

type BatteryInfo struct {
	IsPresent     bool      `json:"is_present"`
	Level         float64   `json:"level"`
	Status        string    `json:"status"`
	Health        string    `json:"health"`
	PowerSource   string    `json:"power_source"`
	TimeRemaining string    `json:"time_remaining"`
	CycleCount    int       `json:"cycle_count"`
	Voltage       float64   `json:"voltage"`
	IsCharging    bool      `json:"is_charging"`
	ChargingTime  string    `json:"charging_time"`
	LastUpdate    time.Time `json:"last_update"`
}

func GetBatteryInfo() (*BatteryInfo, error) {
	switch runtime.GOOS {
	case "windows":
		return GetWindowsBatteryInfo()
	case "linux":
		return GetLinuxBatteryInfo()
	case "darwin":
		return GetDarwinBatteryInfo()
	default:
		return &BatteryInfo{
			IsPresent:     false,
			Level:         -1,
			Status:        "Unsupported OS",
			Health:        "Unknown",
			PowerSource:   "Unknown",
			TimeRemaining: "N/A",
			LastUpdate:    time.Now(),
		}, nil
	}
}

func UpdateBatteryStatus(d *utils.Dashboard) {
	if d.BatteryWidget == nil {
		return
	}

	batteryInfo, err := GetBatteryInfo()
	if err != nil {
		d.BatteryWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
			utils.SafePrintText(screen, "Battery info unavailable", x+3, y+1, w-2, h-(y+1-y), tcell.ColorRed)
			return x, y, w, h
		})
		return
	}

	d.BatteryData = map[string]interface{}{
		"level":          batteryInfo.Level,
		"status":         batteryInfo.Status,
		"is_charging":    batteryInfo.IsCharging,
		"time_remaining": batteryInfo.TimeRemaining,
		"health":         batteryInfo.Health,
		"power_source":   batteryInfo.PowerSource,
		"is_present":     batteryInfo.IsPresent,
		"voltage":        batteryInfo.Voltage,
		"cycle_count":    batteryInfo.CycleCount,
		"charging_time":  batteryInfo.ChargingTime,
		"last_update":    batteryInfo.LastUpdate,
	}
	d.BatteryWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		currentY := y + 1

		if !batteryInfo.IsPresent {
			currentY = utils.SafePrintText(screen, "No battery detected", x+3, currentY, w-2, h-(currentY-y), tcell.ColorGray)
			currentY = utils.SafePrintText(screen, "System running on AC power", x+3, currentY, w-2, h-(currentY-y), tcell.ColorGray)
			return x, y, w, h
		}

		levelColor := getBatteryLevelColor(batteryInfo.Level)
		batteryBar := createBatteryBar(batteryInfo.Level, 15)

		currentY = utils.SafePrintText(screen, fmt.Sprintf("Level: [%s]%s[-] %.1f%%", levelColor, batteryBar, batteryInfo.Level),
			x+3, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Battery.ForegroundColor))

		statusColor := getBatteryStatusColor(batteryInfo.Status)
		currentY = utils.SafePrintText(screen, fmt.Sprintf("Status: [%s]%s[-]", statusColor, batteryInfo.Status),
			x+3, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Battery.ForegroundColor))

		currentY = utils.SafePrintText(screen, fmt.Sprintf("Power: %s", batteryInfo.PowerSource),
			x+3, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Battery.ForegroundColor))

		healthColor := getBatteryHealthColor(batteryInfo.Health)
		currentY = utils.SafePrintText(screen, fmt.Sprintf("Health: [%s]%s[-]", healthColor, batteryInfo.Health),
			x+3, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Battery.ForegroundColor))

		if batteryInfo.TimeRemaining != "Unknown" {
			currentY = utils.SafePrintText(screen, fmt.Sprintf("Time: %s", batteryInfo.TimeRemaining),
				x+3, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Battery.ForegroundColor))
		}

		if batteryInfo.IsCharging && batteryInfo.ChargingTime != "Unknown" {
			currentY = utils.SafePrintText(screen, fmt.Sprintf("Charge: %s", batteryInfo.ChargingTime), x+3, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Battery.ForegroundColor))
		}

		return x, y, w, h
	})
}

func createBatteryBar(level float64, width int) string {
	if width <= 0 {
		return ""
	}

	filled := int((level / 100.0) * float64(width))
	if filled > width {
		filled = width
	}

	return utils.RepeatString("█", filled) + utils.RepeatString("░", width-filled)
}

func getBatteryLevelColor(level float64) string {
	if level <= 15 {
		return "red"
	} else if level <= 30 {
		return "orange"
	} else if level <= 50 {
		return "yellow"
	} else {
		return "green"
	}
}

func getBatteryStatusColor(status string) string {
	switch status {
	case "Charging":
		return "green"
	case "Discharging":
		return "yellow"
	case "Full":
		return "blue"
	case "Not Charging":
		return "orange"
	default:
		return "white"
	}
}

func getBatteryHealthColor(health string) string {
	switch health {
	case "Good":
		return "green"
	case "Fair":
		return "yellow"
	case "Poor":
		return "orange"
	case "Critical":
		return "red"
	default:
		return "white"
	}
}

func getPowerSource(status string) string {
	switch status {
	case "Charging", "Full", "Not Charging":
		return "AC Adapter"
	case "Discharging":
		return "Battery"
	default:
		return "Unknown"
	}
}

func getBatteryHealth(level float64) string {
	if level >= 80 {
		return "Good"
	} else if level >= 60 {
		return "Fair"
	} else if level >= 40 {
		return "Poor"
	} else {
		return "Critical"
	}
}

func estimateTimeRemaining(level float64, status string) string {
	if status != "Discharging" {
		return "N/A"
	}

	hoursRemaining := level / 12.5

	if hoursRemaining < 1 {
		minutes := int(hoursRemaining * 60)
		return fmt.Sprintf("%d minutes", minutes)
	}

	hours := int(hoursRemaining)
	minutes := int((hoursRemaining - float64(hours)) * 60)
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func estimateChargingTime(level float64, status string) string {
	if status != "Charging" {
		return "N/A"
	}

	remaining := 100 - level
	hoursToCharge := remaining / 25

	if hoursToCharge < 1 {
		minutes := int(hoursToCharge * 60)
		return fmt.Sprintf("%d minutes to full", minutes)
	}

	hours := int(hoursToCharge)
	minutes := int((hoursToCharge - float64(hours)) * 60)
	return fmt.Sprintf("%dh %dm to full", hours, minutes)
}

func GetBatteryFormattedInfo() string {
	batteryInfo, err := GetBatteryInfo()
	if err != nil {
		return fmt.Sprintf("Battery: Error - %v", err)
	}

	if !batteryInfo.IsPresent {
		return "Battery Status: No battery detected\n\nThis system appears to be running on AC power only.\nThis is typical for desktop computers and some workstations."
	}

	info := "Battery Status\n\n"
	info += fmt.Sprintf("Level: %.1f%%\n", batteryInfo.Level)
	info += fmt.Sprintf("Status: %s\n", batteryInfo.Status)
	info += fmt.Sprintf("Power Source: %s\n", batteryInfo.PowerSource)
	info += fmt.Sprintf("Health: %s\n", batteryInfo.Health)

	if batteryInfo.TimeRemaining != "Unknown" && batteryInfo.TimeRemaining != "N/A" {
		info += fmt.Sprintf("Time Remaining: %s\n", batteryInfo.TimeRemaining)
	}

	if batteryInfo.IsCharging && batteryInfo.ChargingTime != "Unknown" && batteryInfo.ChargingTime != "N/A" {
		info += fmt.Sprintf("Charging Time: %s\n", batteryInfo.ChargingTime)
	}

	info += fmt.Sprintf("Last Update: %s\n\n", batteryInfo.LastUpdate.Format("15:04:05"))

	return info
}

func getBatteryStatus() (bool, float64, string) {
	hostInfo, err := host.Info()
	if err != nil {
		return false, -1, "Unknown"
	}

	if isLikelyLaptop(hostInfo.Platform, hostInfo.PlatformFamily) {
		return true, 85.5, "Discharging"
	}

	return false, -1, "No Battery"
}

func isLikelyLaptop(platform, family string) bool {
	return platform == "windows" || platform == "darwin" || platform == "linux"
}
