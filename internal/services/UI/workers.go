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
	startIndividualWidgetWorkers(d, quit)
	startPluginUpdateWorker(d, quit)
}

func startIndividualWidgetWorkers(d *utils.Dashboard, quit chan struct{}) {
	startWidgetWorker(d, quit, "cpu", func() { sysinfo.UpdateCPU(d) }, d.Theme.Layout.CPU)
	startWidgetWorker(d, quit, "memory", func() { memory.UpdateVMem(d) }, d.Theme.Layout.Memory)
	startWidgetWorker(d, quit, "disk", func() { disk.UpdateDisk(d) }, d.Theme.Layout.Disk)
	startWidgetWorker(d, quit, "network", func() { network.UpdateNetwork(d) }, d.Theme.Layout.Network)
	startWidgetWorker(d, quit, "process", func() {
		processes.UpdateProcesses(d)
		if d.ProcessWidget != nil {
			sortLabel := "CPU"
			if d.Theme.Sorting != "" {
				sortLabel = formatSort(d.Theme.Sorting)
			}
			d.ProcessWidget.SetTitle(fmt.Sprint("Processes - ", processes.GetNrProcesses(), " Sorted by: ", sortLabel))
		}
	}, d.Theme.Layout.Process)
	startWidgetWorker(d, quit, "gpu", func() { gpu.UpdateGPU(d) }, d.Theme.Layout.GPU)
	startWidgetWorker(d, quit, "load", func() { load.UpdateLoadAverage(d) }, d.Theme.Layout.Load)
	startWidgetWorker(d, quit, "temperature", func() { temperature.UpdateTemperatures(d) }, d.Theme.Layout.Temperature)
	startWidgetWorker(d, quit, "network_connections", func() { network.UpdateNetworkConnections(d) }, d.Theme.Layout.NetworkConns)
	startWidgetWorker(d, quit, "disk_io", func() { disk.UpdateDiskIO(d) }, d.Theme.Layout.DiskIO)
	startWidgetWorker(d, quit, "process_tree", func() { processes.UpdateProcessTree(d) }, d.Theme.Layout.ProcessTree)
	startWidgetWorker(d, quit, "battery", func() { battery.UpdateBatteryStatus(d) }, d.Theme.Layout.Battery)

	startWidgetWorker(d, quit, "header", func() { updateHeaderTitle(d) }, utils.WidgetConfig{Enabled: true, UpdateInterval: 1})

	performInitialUpdates(d)
}

func startWidgetWorker(d *utils.Dashboard, quit chan struct{}, widgetName string, updateFunc func(), config utils.WidgetConfig) {
	if !config.Enabled {
		return
	}

	interval := config.UpdateInterval
	if interval <= 0 {
		interval = 1
	}

	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				updateFunc()
				d.App.QueueUpdateDraw(func() {})
			case <-quit:
				return
			}
		}
	}()
}

func performInitialUpdates(d *utils.Dashboard) {
	if d.Theme.Layout.CPU.Enabled {
		sysinfo.UpdateCPU(d)
	}
	if d.Theme.Layout.Memory.Enabled {
		memory.UpdateVMem(d)
	}
	if d.Theme.Layout.Disk.Enabled {
		disk.UpdateDisk(d)
	}
	if d.Theme.Layout.Network.Enabled {
		network.UpdateNetwork(d)
	}
	if d.Theme.Layout.Process.Enabled {
		processes.UpdateProcesses(d)
		if d.ProcessWidget != nil {
			sortLabel := "CPU"
			if d.Theme.Sorting != "" {
				sortLabel = formatSort(d.Theme.Sorting)
			}
			d.ProcessWidget.SetTitle(fmt.Sprint("Processes - ", processes.GetNrProcesses(), " Sorted by: ", sortLabel))
		}
	}
	if d.Theme.Layout.GPU.Enabled {
		gpu.UpdateGPU(d)
	}
	if d.Theme.Layout.Load.Enabled {
		load.UpdateLoadAverage(d)
	}
	if d.Theme.Layout.Temperature.Enabled {
		temperature.UpdateTemperatures(d)
	}
	if d.Theme.Layout.NetworkConns.Enabled {
		network.UpdateNetworkConnections(d)
	}
	if d.Theme.Layout.DiskIO.Enabled {
		disk.UpdateDiskIO(d)
	}
	if d.Theme.Layout.ProcessTree.Enabled {
		processes.UpdateProcessTree(d)
	}
	if d.Theme.Layout.Battery.Enabled {
		battery.UpdateBatteryStatus(d)
	}

	updateHeaderTitle(d)
}

func startPluginUpdateWorker(d *utils.Dashboard, quit chan struct{}) {
	if d.PluginManager == nil {
		return
	}

	pluginManager, ok := d.PluginManager.(*plugins.PluginManager)
	if !ok {
		return
	}

	pluginInfo := pluginManager.GetPluginInfo()
	for _, info := range pluginInfo {
		if info.Config.Enabled {
			plugin, exists := pluginManager.GetPlugin(info.Name)
			if exists {
				startPluginWorker(d, plugin, info.Config, quit)
			}
		}
	}
}

func startPluginWorker(d *utils.Dashboard, plugin plugins.Plugin, config plugins.PluginConfig, quit chan struct{}) {
	interval := config.Layout.UpdateInterval
	if interval <= 0 {
		interval = 5
	}

	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if d.PluginManager != nil {
					if pluginManager, ok := d.PluginManager.(*plugins.PluginManager); ok {
						if widget, exists := pluginManager.GetWidget(plugin.Name()); exists {
							plugin.UpdateWidget(widget)
						}
					}
				}
			case <-quit:
				return
			}
		}
	}()
}
