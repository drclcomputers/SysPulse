//go:build windows
// +build windows

package battery

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type win32_Battery struct {
	BatteryStatus               uint16
	EstimatedChargeRemaining    uint16
	EstimatedRunTime            uint32
	TimeToFullCharge            uint32
	DesignCapacity              uint32
	FullChargeCapacity          uint32
	Voltage                     uint32
	BatteryRechargeTime         uint32
	PowerManagementCapabilities []uint16
	SystemCreationClassName     string
	SystemName                  string
	CreationClassName           string
	DeviceID                    string
	Name                        string
	Tag                         string
	Status                      string
	Availability                uint16
	Chemistry                   uint16
	Caption                     string
	Description                 string
}

type win32_PowerSupply struct {
	Range1SwitchPosition    uint16
	Range2SwitchPosition    uint16
	TotalOutputPower        uint32
	TypeOfRangeSwitching    uint16
	ActiveInputVoltage      uint32
	InputVoltageRange       uint16
	MaxNumberOfOutputs      uint32
	NumberOfPowerCords      uint32
	OutputVoltage           uint32
	SystemCreationClassName string
	SystemName              string
	CreationClassName       string
	DeviceID                string
	Name                    string
	Status                  string
	Availability            uint16
	Caption                 string
	Description             string
}

func GetWindowsBatteryInfo() (*BatteryInfo, error) {
	return getPowerShellBatteryInfo()
}

func processWMIBatteryInfo(battery win32_Battery) (*BatteryInfo, error) {
	status := "Unknown"
	batteryStatus := battery.BatteryStatus
	switch batteryStatus {
	case 1:
		status = "Discharging"
	case 2:
		status = "Charging"
	case 3:
		status = "Charged"
	case 4:
		status = "Low"
	case 5:
		status = "Critical"
	case 6:
		status = "Charging and High"
	case 7:
		status = "Charging and Low"
	case 8:
		status = "Charging and Critical"
	case 9:
		status = "Undefined"
	case 10:
		status = "Partially Charged"
	}

	timeRemaining := "Unknown"
	estimatedRunTime := battery.EstimatedRunTime
	if estimatedRunTime != 0 && estimatedRunTime != 71582788 {
		hours := estimatedRunTime / 60
		minutes := estimatedRunTime % 60
		timeRemaining = fmt.Sprintf("%dh %dm", hours, minutes)
	}

	powerSource := "Battery"
	if batteryStatus == 2 || batteryStatus == 6 || batteryStatus == 7 || batteryStatus == 8 {
		powerSource = "AC Power"
	}

	return &BatteryInfo{
		IsPresent:     true,
		Level:         float64(battery.EstimatedChargeRemaining),
		Status:        status,
		Health:        battery.Status,
		PowerSource:   powerSource,
		TimeRemaining: timeRemaining,
		LastUpdate:    time.Now(),
	}, nil
}

func getPowerShellBatteryInfo() (*BatteryInfo, error) {
	cmd := exec.Command("powershell", "-Command",
		"$battery = Get-WmiObject -Class Win32_Battery | Select-Object -First 1; "+
			"if ($battery) { "+
			"Write-Output \"$($battery.EstimatedChargeRemaining)|$($battery.BatteryStatus)|$($battery.Status)|$($battery.EstimatedRunTime)\" "+
			"} else { "+
			"Write-Output \"NOBATTERY\" "+
			"}")

	output, err := cmd.Output()
	if err != nil {
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

	result := strings.TrimSpace(string(output))
	if result == "NOBATTERY" || result == "" {
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

	parts := strings.Split(result, "|")
	if len(parts) < 4 {
		return &BatteryInfo{
			IsPresent:     false,
			Level:         -1,
			Status:        "Parse Error",
			Health:        "Unknown",
			PowerSource:   "AC Power",
			TimeRemaining: "N/A",
			LastUpdate:    time.Now(),
		}, nil
	}

	level, _ := strconv.ParseFloat(parts[0], 64)
	batteryStatus, _ := strconv.Atoi(parts[1])
	healthStatus := parts[2]
	runTime, _ := strconv.Atoi(parts[3])

	status := "Unknown"
	switch batteryStatus {
	case 1:
		status = "Discharging"
	case 2:
		status = "Charging"
	case 3:
		status = "Charged"
	case 4:
		status = "Low"
	case 5:
		status = "Critical"
	case 6:
		status = "Charging and High"
	case 7:
		status = "Charging and Low"
	case 8:
		status = "Charging and Critical"
	case 9:
		status = "Undefined"
	case 10:
		status = "Partially Charged"
	}

	timeRemaining := "Unknown"
	if runTime != 0 && runTime != 71582788 {
		hours := runTime / 60
		minutes := runTime % 60
		timeRemaining = fmt.Sprintf("%dh %dm", hours, minutes)
	}

	powerSource := "Battery"
	if batteryStatus == 2 || batteryStatus == 6 || batteryStatus == 7 || batteryStatus == 8 {
		powerSource = "AC Power"
	}

	return &BatteryInfo{
		IsPresent:     true,
		Level:         level,
		Status:        status,
		Health:        healthStatus,
		PowerSource:   powerSource,
		TimeRemaining: timeRemaining,
		IsCharging:    batteryStatus == 2,
		LastUpdate:    time.Now(),
	}, nil
}
