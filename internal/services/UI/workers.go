package ui

import (
	"fmt"
	"syspulse/internal/plugins"
	"syspulse/internal/services/battery"
	"syspulse/internal/services/disk"
	"syspulse/internal/services/gpu"
	"syspulse/internal/services/load"
	"syspulse/internal/services/memory"
	"syspulse/internal/services/network"
	"syspulse/internal/services/processes"
	"syspulse/internal/services/sysinfo"
	"syspulse/internal/services/temperature"
	"syspulse/internal/utils"
	"time"
)

func startWorkers(d *utils.Dashboard, quit chan struct{}) {
	startOptimizedUpdateWorker(d, quit)
	startPluginUpdateWorker(d, quit)
}

func startOptimizedUpdateWorker(d *utils.Dashboard, quit chan struct{}) {
	performHighFrequencyUpdates(d)
	performMediumFrequencyUpdates(d)
	performLowFrequencyUpdates(d)
	performVeryLowFrequencyUpdates(d)

	go func() {
		updateInterval := time.Second
		if d.Theme.UpdateTime > 0 {
			updateInterval = time.Duration(d.Theme.UpdateTime) * time.Second
		}
		ticker := time.NewTicker(updateInterval)
		defer ticker.Stop()

		updateCount := 0

		for {
			select {
			case <-ticker.C:
				updateCount++

				performHighFrequencyUpdates(d)

				if updateCount%3 == 0 {
					performMediumFrequencyUpdates(d)
				}

				if updateCount%6 == 0 {
					performLowFrequencyUpdates(d)
				}

				if updateCount%9 == 0 {
					performVeryLowFrequencyUpdates(d)
				}

				d.App.QueueUpdateDraw(func() {})

			case <-quit:
				return
			}
		}
	}()
}

func performHighFrequencyUpdates(d *utils.Dashboard) {
	sysinfo.UpdateCPU(d)
	network.UpdateNetwork(d)
	updateHeaderTitle(d)
}

func performMediumFrequencyUpdates(d *utils.Dashboard) {
	processes.UpdateProcesses(d)
	if d.ProcessWidget != nil {
		sortLabel := "CPU"
		if d.Theme.Sorting != "" {
			sortLabel = formatSort(d.Theme.Sorting)
		}
		d.ProcessWidget.SetTitle(fmt.Sprint("Processes - ", processes.GetNrProcesses(), " Sorted by: ", sortLabel))
	}
	memory.UpdateVMem(d)
}

func performLowFrequencyUpdates(d *utils.Dashboard) {
	load.UpdateLoadAverage(d)
	network.UpdateNetworkConnections(d)
	disk.UpdateDisk(d)
	disk.UpdateDiskIO(d)
}

func performVeryLowFrequencyUpdates(d *utils.Dashboard) {
	gpu.UpdateGPU(d)
	temperature.UpdateTemperatures(d)
	battery.UpdateBatteryStatus(d)
	processes.UpdateProcessTree(d)
}

func startPluginUpdateWorker(d *utils.Dashboard, quit chan struct{}) {
	if d.PluginManager == nil {
		return
	}

	pluginManager, ok := d.PluginManager.(*plugins.PluginManager)
	if !ok {
		return
	}

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				pluginManager.UpdatePlugins()
			case <-quit:
				return
			}
		}
	}()
}
