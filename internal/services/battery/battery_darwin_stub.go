//go:build !darwin
// +build !darwin

package battery

import (
	"time"
)

func GetDarwinBatteryInfo() (*BatteryInfo, error) {
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
