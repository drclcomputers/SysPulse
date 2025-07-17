package ui

import (
	"fmt"
	"syspulse/internal/export"
	"syspulse/internal/services/disk"
	"syspulse/internal/services/gpu"
	"syspulse/internal/services/memory"
	"syspulse/internal/services/network"
	"syspulse/internal/services/sysinfo"
	"syspulse/internal/utils"
	"time"
)

func startExportWorker(d *utils.Dashboard, quit chan struct{}) {
	go func() {
		if !d.Theme.Export.Enabled {
			return
		}

		ticker := time.NewTicker(time.Duration(d.Theme.Export.Interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				performPeriodicExport(d)
			case <-quit:
				return
			}
		}
	}()
}

func performPeriodicExport(d *utils.Dashboard) {
	sysinfo.UpdateCPU(d)
	memory.UpdateVMem(d)
	disk.UpdateDisk(d)
	network.UpdateNetwork(d)
	if d.Theme.Layout.GPU.Enabled {
		gpu.UpdateGPU(d)
	}

	snapshot := export.CreateSnapshot(d)
	dataPoints = append(dataPoints, snapshot)

	timestamp := time.Now().Format("2006-01-02_15-04-05")

	for _, format := range d.Theme.Export.Formats {
		filename := fmt.Sprintf("%s/%s_%s.%s",
			d.Theme.Export.Directory,
			d.Theme.Export.FilenamePrefix,
			timestamp,
			format)

		var exportFormat export.ExportFormat
		switch format {
		case "csv":
			exportFormat = export.CSV
		case "json":
			exportFormat = export.JSON
		default:
			log.Error(fmt.Sprintf("Unsupported export format: %s", format))
			continue
		}

		if err := export.ExportData(dataPoints, filename, exportFormat); err != nil {
			log.Error(fmt.Sprintf("Failed to export %s: %v", format, err))
		} else {
		}
	}
}

func performFinalExport(d *utils.Dashboard) {
	if !d.Theme.Export.Enabled {
		return
	}

	sysinfo.UpdateCPU(d)
	memory.UpdateVMem(d)
	disk.UpdateDisk(d)
	network.UpdateNetwork(d)
	if d.Theme.Layout.GPU.Enabled {
		gpu.UpdateGPU(d)
	}

	finalSnapshot := export.CreateSnapshot(d)
	dataPoints = append(dataPoints, finalSnapshot)

	timestamp := time.Now().Format("2006-01-02_15-04-05")

	for _, format := range d.Theme.Export.Formats {
		filename := fmt.Sprintf("%s/%s_final_%s.%s",
			d.Theme.Export.Directory,
			d.Theme.Export.FilenamePrefix,
			timestamp,
			format)

		var exportFormat export.ExportFormat
		switch format {
		case "csv":
			exportFormat = export.CSV
		case "json":
			exportFormat = export.JSON
		default:
			log.Error(fmt.Sprintf("Unsupported export format: %s", format))
			continue
		}

		if err := export.ExportData(dataPoints, filename, exportFormat); err != nil {
			log.Error(fmt.Sprintf("Failed to export final %s: %v", format, err))
		} else {

		}
	}
}
