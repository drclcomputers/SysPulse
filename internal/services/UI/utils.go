package ui

import (
	"fmt"
	"syspulse/internal/services/sysinfo"
	"syspulse/internal/utils"
)

type Dashboard utils.Dashboard

func formatSort(sorttype string) string {
	if sorttype == "mem" {
		return "RAM"
	} else {
		return "CPU"
	}
}

func createHeaderTitle() string {
	return fmt.Sprint("SysPulse v", utils.VER, " | ", sysinfo.GetPlatform(), " | User: ", sysinfo.GetCurrentUser(), " | Uptime: ", sysinfo.GetUptime())
}

func updateHeaderTitle(d *utils.Dashboard) {
	d.HeaderWidget.SetTitle(createHeaderTitle())
}
