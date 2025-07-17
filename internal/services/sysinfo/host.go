package sysinfo

import (
	"fmt"
	"syspulse/internal/utils"
	"time"

	"github.com/shirou/gopsutil/host"
)

func formatDuration(seconds uint64) string {
	duration := time.Duration(seconds) * time.Second

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	secs := int(duration.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

func HostInfo() {
	info, _ := host.Info()
	fmt.Printf("Hostname: %s, OS: %s %s\n", info.Hostname, info.Platform, info.PlatformVersion)
}

func GetUptime() string {
	info, err := host.Info()
	if err != nil {
		return "Unknown"
	}

	return utils.FormatTime(int64(info.Uptime))
}

func GetPlatform() string {
	info, err := host.Info()
	if err != nil {
		return "Unknown"
	}

	return fmt.Sprint(info.Platform)
}

func GetSystemInfo() string {
	info, err := host.Info()
	if err != nil {
		return "Unable to get system info!"
	}

	uptime := formatDuration(info.Uptime)

	bootTime := utils.FormatTime(int64(info.BootTime))

	output := fmt.Sprintf("Hostname: %s\nOS: %s\nPlatform: %s\nPlatform Version: %s\nKernel Version: %s\nKernel Arch: %s\nUptime: %s\nBoot time: %s\nVirtualization: %s\nVirtualization role: %s\nHost ID: %s",
		info.Hostname, info.OS, info.Platform, info.PlatformVersion, info.KernelVersion, info.KernelArch, uptime, bootTime, info.VirtualizationSystem, info.VirtualizationRole, info.HostID)

	return output
}
