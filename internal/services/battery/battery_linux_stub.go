//go:build !linux
// +build !linux

package battery

import (
	"time"
)

func GetLinuxBatteryInfo() (*BatteryInfo, error) {
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
