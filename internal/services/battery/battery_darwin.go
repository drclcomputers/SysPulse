//go:build darwin
// +build darwin

package battery

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetDarwinBatteryInfo() (*BatteryInfo, error) {
	if info, err := getDarwinBatteryInfoPMSet(); err == nil {
		return info, nil
	}

	if info, err := getDarwinBatteryInfoSystemProfiler(); err == nil {
		return info, nil
	}

	if info, err := getDarwinBatteryInfoIOREG(); err == nil {
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

func getDarwinBatteryInfoPMSet() (*BatteryInfo, error) {
	cmd := exec.Command("pmset", "-g", "batt")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("pmset command failed: %v", err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "InternalBattery") {
		return nil, fmt.Errorf("no internal battery found")
	}

	info := &BatteryInfo{
		IsPresent:  true,
		LastUpdate: time.Now(),
	}

	percentageRe := regexp.MustCompile(`(\d+)%`)
	if matches := percentageRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if percentage, err := strconv.Atoi(matches[1]); err == nil {
			info.Level = float64(percentage)
		}
	}

	statusRe := regexp.MustCompile(`;\s*(\w+(?:\s+\w+)*);`)
	if matches := statusRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		status := strings.TrimSpace(matches[1])
		info.Status = status
		info.IsCharging = strings.Contains(strings.ToLower(status), "charging")
		info.PowerSource = getPowerSourceFromDarwinStatus(status)
	}

	timeRe := regexp.MustCompile(`(\d+:\d+)\s+remaining`)
	if matches := timeRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		info.TimeRemaining = matches[1]
	} else if strings.Contains(outputStr, "calculating") {
		info.TimeRemaining = "Calculating..."
	}

	return info, nil
}

func getDarwinBatteryInfoSystemProfiler() (*BatteryInfo, error) {
	cmd := exec.Command("system_profiler", "SPPowerDataType")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("system_profiler command failed: %v", err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "Battery Information") {
		return nil, fmt.Errorf("no battery information found")
	}

	info := &BatteryInfo{
		IsPresent:  true,
		LastUpdate: time.Now(),
	}

	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "State of Charge (%):") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				chargeStr := strings.TrimSpace(parts[1])
				if charge, err := strconv.ParseFloat(chargeStr, 64); err == nil {
					info.Level = charge
				}
			}
		} else if strings.Contains(line, "Charging:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				charging := strings.TrimSpace(parts[1])
				info.IsCharging = strings.EqualFold(charging, "Yes")
				if info.IsCharging {
					info.Status = "Charging"
					info.PowerSource = "AC Power"
				} else {
					info.Status = "Discharging"
					info.PowerSource = "Battery"
				}
			}
		} else if strings.Contains(line, "Condition:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				condition := strings.TrimSpace(parts[1])
				info.Health = condition
			}
		} else if strings.Contains(line, "Cycle Count:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				cycleStr := strings.TrimSpace(parts[1])
				if cycles, err := strconv.Atoi(cycleStr); err == nil {
					info.CycleCount = cycles
				}
			}
		}
	}

	return info, nil
}

func getDarwinBatteryInfoIOREG() (*BatteryInfo, error) {
	cmd := exec.Command("ioreg", "-rc", "AppleSmartBattery")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ioreg command failed: %v", err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "AppleSmartBattery") {
		return nil, fmt.Errorf("no AppleSmartBattery found")
	}

	info := &BatteryInfo{
		IsPresent:  true,
		LastUpdate: time.Now(),
	}

	currentCapacityRe := regexp.MustCompile(`"CurrentCapacity"\s*=\s*(\d+)`)
	maxCapacityRe := regexp.MustCompile(`"MaxCapacity"\s*=\s*(\d+)`)

	var currentCapacity, maxCapacity int

	if matches := currentCapacityRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if cap, err := strconv.Atoi(matches[1]); err == nil {
			currentCapacity = cap
		}
	}

	if matches := maxCapacityRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if cap, err := strconv.Atoi(matches[1]); err == nil {
			maxCapacity = cap
		}
	}

	if maxCapacity > 0 {
		info.Level = float64(currentCapacity) / float64(maxCapacity) * 100
	}

	chargingRe := regexp.MustCompile(`"IsCharging"\s*=\s*(\w+)`)
	if matches := chargingRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		info.IsCharging = strings.EqualFold(matches[1], "Yes")
		if info.IsCharging {
			info.Status = "Charging"
			info.PowerSource = "AC Power"
		} else {
			info.Status = "Discharging"
			info.PowerSource = "Battery"
		}
	}

	cycleCountRe := regexp.MustCompile(`"CycleCount"\s*=\s*(\d+)`)
	if matches := cycleCountRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if cycles, err := strconv.Atoi(matches[1]); err == nil {
			info.CycleCount = cycles
		}
	}

	voltageRe := regexp.MustCompile(`"Voltage"\s*=\s*(\d+)`)
	if matches := voltageRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if voltage, err := strconv.Atoi(matches[1]); err == nil {
			info.Voltage = float64(voltage) / 1000.0
		}
	}

	designCapacityRe := regexp.MustCompile(`"DesignCapacity"\s*=\s*(\d+)`)
	if matches := designCapacityRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if designCap, err := strconv.Atoi(matches[1]); err == nil && designCap > 0 {
			healthPercentage := float64(maxCapacity) / float64(designCap) * 100
			info.Health = getBatteryHealthFromDarwinPercentage(healthPercentage)
		}
	}

	return info, nil
}

func getPowerSourceFromDarwinStatus(status string) string {
	switch strings.ToLower(status) {
	case "charging", "charged":
		return "AC Power"
	case "discharging", "not charging":
		return "Battery"
	default:
		return "Unknown"
	}
}

func getBatteryHealthFromDarwinPercentage(percentage float64) string {
	switch {
	case percentage >= 90:
		return "Excellent"
	case percentage >= 80:
		return "Good"
	case percentage >= 60:
		return "Fair"
	case percentage >= 40:
		return "Poor"
	default:
		return "Replace Soon"
	}
}
