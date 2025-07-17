//go:build linux
// +build linux

package battery

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	powerSupplyPath = "/sys/class/power_supply"
	acpiPath        = "/proc/acpi/battery"
)

func GetLinuxBatteryInfo() (*BatteryInfo, error) {
	if info, err := getLinuxBatteryInfoSysfs(); err == nil {
		return info, nil
	}

	if info, err := getLinuxBatteryInfoACPI(); err == nil {
		return info, nil
	}

	return &BatteryInfo{
		IsPresent:     false,
		Level:         -1,
		Status:        "No Battery",
		Health:        "Unknown",
		PowerSource:   "AC Power",
		TimeRemaining: "N/A",
		LastUpdate:    time.Now(),
	}, nil
}

func getLinuxBatteryInfoSysfs() (*BatteryInfo, error) {
	entries, err := ioutil.ReadDir(powerSupplyPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read power supply directory: %v", err)
	}

	var batteryPath string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "BAT") {
			batteryPath = filepath.Join(powerSupplyPath, entry.Name())
			break
		}
	}

	if batteryPath == "" {
		return nil, fmt.Errorf("no battery found in power supply directory")
	}

	info := &BatteryInfo{
		IsPresent:  true,
		LastUpdate: time.Now(),
	}

	if capacity, err := readSysfsInt(batteryPath, "capacity"); err == nil {
		info.Level = float64(capacity)
	}

	if status, err := readSysfsString(batteryPath, "status"); err == nil {
		info.Status = strings.TrimSpace(status)
		info.IsCharging = strings.EqualFold(info.Status, "Charging")
		info.PowerSource = getPowerSourceFromLinuxStatus(info.Status)
	}

	if voltage, err := readSysfsInt(batteryPath, "voltage_now"); err == nil {
		info.Voltage = float64(voltage) / 1000000.0
	}

	if designCapacity, err := readSysfsInt(batteryPath, "energy_full_design"); err == nil {
		if fullCapacity, err := readSysfsInt(batteryPath, "energy_full"); err == nil {
			percentage := float64(fullCapacity) / float64(designCapacity) * 100
			info.Health = getBatteryHealthFromPercentage(percentage)
		}
	}

	if currentNow, err := readSysfsInt(batteryPath, "current_now"); err == nil && currentNow > 0 {
		if energyNow, err := readSysfsInt(batteryPath, "energy_now"); err == nil {
			hoursRemaining := float64(energyNow) / float64(currentNow)
			info.TimeRemaining = formatLinuxTimeRemaining(hoursRemaining)
		}
	}

	return info, nil
}

func getLinuxBatteryInfoACPI() (*BatteryInfo, error) {
	entries, err := ioutil.ReadDir(acpiPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read ACPI battery directory: %v", err)
	}

	var batteryDir string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "BAT") {
			batteryDir = filepath.Join(acpiPath, entry.Name())
			break
		}
	}

	if batteryDir == "" {
		return nil, fmt.Errorf("no battery found in ACPI directory")
	}

	info := &BatteryInfo{
		IsPresent:  true,
		LastUpdate: time.Now(),
	}

	if stateData, err := os.ReadFile(filepath.Join(batteryDir, "state")); err == nil {
		lines := strings.Split(string(stateData), "\n")
		for _, line := range lines {
			if strings.Contains(line, "charging state:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					state := strings.TrimSpace(parts[1])
					info.Status = state
					info.IsCharging = strings.EqualFold(state, "charging")
					info.PowerSource = getPowerSourceFromLinuxStatus(state)
				}
			} else if strings.Contains(line, "remaining capacity:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					capacityStr := strings.TrimSpace(strings.Fields(parts[1])[0])
					if capacity, err := strconv.Atoi(capacityStr); err == nil {
						info.Level = float64(capacity)
					}
				}
			}
		}
	}

	if infoData, err := ioutil.ReadFile(filepath.Join(batteryDir, "info")); err == nil {
		lines := strings.Split(string(infoData), "\n")
		var fullCapacity int
		for _, line := range lines {
			if strings.Contains(line, "last full capacity:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					capacityStr := strings.TrimSpace(strings.Fields(parts[1])[0])
					if capacity, err := strconv.Atoi(capacityStr); err == nil {
						fullCapacity = capacity
					}
				}
			}
		}

		if fullCapacity > 0 && info.Level > 0 {
			info.Level = (info.Level / float64(fullCapacity)) * 100
		}
	}

	return info, nil
}

func readSysfsInt(basePath, filename string) (int, error) {
	data, err := os.ReadFile(filepath.Join(basePath, filename))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func readSysfsString(basePath, filename string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(basePath, filename))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getPowerSourceFromLinuxStatus(status string) string {
	switch strings.ToLower(status) {
	case "charging":
		return "AC Power"
	case "discharging", "not charging":
		return "Battery"
	default:
		return "Unknown"
	}
}

func getBatteryHealthFromPercentage(percentage float64) string {
	switch {
	case percentage >= 80:
		return "Excellent"
	case percentage >= 60:
		return "Good"
	case percentage >= 40:
		return "Fair"
	default:
		return "Poor"
	}
}

func formatLinuxTimeRemaining(hours float64) string {
	if hours <= 0 {
		return "Unknown"
	}

	h := int(hours)
	m := int((hours - float64(h)) * 60)

	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}

	return fmt.Sprintf("%dm", m)
}
