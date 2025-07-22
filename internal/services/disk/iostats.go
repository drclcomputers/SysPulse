package disk

import (
	"fmt"
	"syspulse/internal/utils"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/disk"
)

type DiskIOStats struct {
	ReadCount        uint64  `json:"read_count"`
	WriteCount       uint64  `json:"write_count"`
	ReadBytes        uint64  `json:"read_bytes"`
	WriteBytes       uint64  `json:"write_bytes"`
	ReadTime         uint64  `json:"read_time"`
	WriteTime        uint64  `json:"write_time"`
	ReadBytesPerSec  float64 `json:"read_bytes_per_sec"`
	WriteBytesPerSec float64 `json:"write_bytes_per_sec"`
	ReadOpsPerSec    float64 `json:"read_ops_per_sec"`
	WriteOpsPerSec   float64 `json:"write_ops_per_sec"`
	UtilizationPct   float64 `json:"utilization_pct"`
}

type DiskIOData struct {
	Disks    []*DiskIODevice `json:"disks"`
	LastTime time.Time       `json:"last_time"`
}

type DiskIODevice struct {
	Name  string       `json:"name"`
	Stats *DiskIOStats `json:"stats"`
}

var (
	lastIOStats map[string]disk.IOCountersStat
	lastIOTime  time.Time
)

func init() {
	lastIOStats = make(map[string]disk.IOCountersStat)
}

func GetDiskIOStats() (*DiskIOData, error) {
	ioStats, err := disk.IOCounters()
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	ioData := &DiskIOData{
		Disks:    make([]*DiskIODevice, 0),
		LastTime: currentTime,
	}

	deviceNames := make([]string, 0, len(ioStats))
	for device := range ioStats {
		deviceNames = append(deviceNames, device)
	}

	for _, device := range deviceNames {
		stat := ioStats[device]
		diskStat := &DiskIOStats{
			ReadCount:  stat.ReadCount,
			WriteCount: stat.WriteCount,
			ReadBytes:  stat.ReadBytes,
			WriteBytes: stat.WriteBytes,
			ReadTime:   stat.ReadTime,
			WriteTime:  stat.WriteTime,
		}

		if lastStat, exists := lastIOStats[device]; exists && !lastIOTime.IsZero() {
			duration := currentTime.Sub(lastIOTime).Seconds()
			if duration > 0 {
				diskStat.ReadBytesPerSec = float64(stat.ReadBytes-lastStat.ReadBytes) / duration
				diskStat.WriteBytesPerSec = float64(stat.WriteBytes-lastStat.WriteBytes) / duration
				diskStat.ReadOpsPerSec = float64(stat.ReadCount-lastStat.ReadCount) / duration
				diskStat.WriteOpsPerSec = float64(stat.WriteCount-lastStat.WriteCount) / duration

				totalTime := float64(stat.ReadTime + stat.WriteTime - lastStat.ReadTime - lastStat.WriteTime)
				diskStat.UtilizationPct = (totalTime / (duration * 1000)) * 100
				if diskStat.UtilizationPct > 100 {
					diskStat.UtilizationPct = 100
				}
			}
		}

		ioData.Disks = append(ioData.Disks, &DiskIODevice{
			Name:  device,
			Stats: diskStat,
		})
		lastIOStats[device] = stat
	}

	lastIOTime = currentTime
	return ioData, nil
}

func UpdateDiskIO(d *utils.Dashboard) {
	if d.DiskIOWidget == nil {
		return
	}

	ioData, err := GetDiskIOStats()
	if err != nil {
		d.DiskIOWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
			utils.SafePrintText(screen, "Disk I/O stats unavailable", x+3, y+1, w-6, y+h-1, tcell.ColorRed)
			return x, y, w, h
		})
		return
	}

	d.DiskIOData = ioData
	d.DiskIOWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		currentY := y + 1

		for _, device := range ioData.Disks {
			if currentY >= y+h-1 {
				break
			}

			deviceName := truncateDeviceName(device.Name, 8)
			deviceLine := fmt.Sprintf("%s", deviceName)
			currentY = utils.SafePrintText(screen, deviceLine, x+2, currentY, w-6, y+h-1, utils.GetColorFromName(d.Theme.Layout.DiskIO.ForegroundColor))

			if currentY >= y+h-1 {
				break
			}

			readRate := formatBytes(device.Stats.ReadBytesPerSec)
			writeRate := formatBytes(device.Stats.WriteBytesPerSec)

			readColor := getIOColor(device.Stats.ReadBytesPerSec)
			writeColor := getIOColor(device.Stats.WriteBytesPerSec)

			readLine := fmt.Sprintf("  R: [%s]%s/s[-] (%s ops/s)", readColor, readRate, formatNumber(device.Stats.ReadOpsPerSec))
			tview.Print(screen, readLine, x+3, currentY, w-6, y+h-1, utils.GetColorFromName(d.Theme.Layout.DiskIO.ForegroundColor))
			currentY++

			if currentY >= y+h-1 {
				break
			}

			writeLine := fmt.Sprintf("  W: [%s]%s/s[-] (%s ops/s)", writeColor, writeRate, formatNumber(device.Stats.WriteOpsPerSec))
			tview.Print(screen, writeLine, x+3, currentY, w-6, y+h-1, utils.GetColorFromName(d.Theme.Layout.DiskIO.ForegroundColor))
			currentY++

			if currentY >= y+h-1 {
				break
			}

			if device.Stats.UtilizationPct > 0 {
				utilColor := getUtilizationColor(device.Stats.UtilizationPct)
				utilBar := createUtilizationBar(device.Stats.UtilizationPct, 10)
				utilLine := fmt.Sprintf("  Util: [%s]%s[-] %.1f%%", utilColor, utilBar, device.Stats.UtilizationPct)
				currentY = utils.SafePrintText(screen, utilLine, x+3, currentY, w-6, y+h-1, utils.GetColorFromName(d.Theme.Layout.DiskIO.ForegroundColor))
			}

			//currentY++
		}

		return x, y, w, h
	})
}

func formatBytes(bytes float64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	if bytes >= GB {
		return fmt.Sprintf("%.0f GB", bytes/GB)
	} else if bytes >= MB {
		return fmt.Sprintf("%.0f MB", bytes/MB)
	} else if bytes >= KB {
		return fmt.Sprintf("%.0f KB", bytes/KB)
	}
	return fmt.Sprintf("%.0f B", bytes)
}

func formatNumber(num float64) string {
	if num >= 1000000 {
		return fmt.Sprintf("%.1fM", num/1000000)
	} else if num >= 1000 {
		return fmt.Sprintf("%.1fK", num/1000)
	}
	return fmt.Sprintf("%.1f", num)
}

func getIOColor(bytesPerSec float64) string {
	const (
		MB = 1024 * 1024
		GB = MB * 1024
	)

	if bytesPerSec >= GB {
		return "red"
	} else if bytesPerSec >= 100*MB {
		return "orange"
	} else if bytesPerSec >= 10*MB {
		return "yellow"
	} else {
		return "green"
	}
}

func getUtilizationColor(util float64) string {
	if util >= 90 {
		return "red"
	} else if util >= 70 {
		return "orange"
	} else if util >= 50 {
		return "yellow"
	} else {
		return "green"
	}
}

func createUtilizationBar(utilization float64, width int) string {
	if width <= 0 {
		return ""
	}

	filled := int((utilization / 100.0) * float64(width))
	if filled > width {
		filled = width
	}

	return utils.RepeatString("█", filled) + utils.RepeatString("░", width-filled)
}

func truncateDeviceName(name string, maxLen int) string {
	if len(name) <= maxLen {
		return name
	}
	return name[:maxLen-3] + "..."
}

func GetDiskIOFormattedInfo() string {
	ioData, err := GetDiskIOStats()
	if err != nil {
		return fmt.Sprintf("Disk I/O: Error - %v", err)
	}

	var info string

	for _, device := range ioData.Disks {
		info += fmt.Sprintf("Device: %s\n", device.Name)
		info += fmt.Sprintf("  Read Rate: %s/s (%.1f ops/s)\n", formatBytes(device.Stats.ReadBytesPerSec), device.Stats.ReadOpsPerSec)
		info += fmt.Sprintf("  Write Rate: %s/s (%.1f ops/s)\n", formatBytes(device.Stats.WriteBytesPerSec), device.Stats.WriteOpsPerSec)
		info += fmt.Sprintf("  Total Read: %s (%d operations)\n", formatBytes(float64(device.Stats.ReadBytes)), device.Stats.ReadCount)
		info += fmt.Sprintf("  Total Write: %s (%d operations)\n", formatBytes(float64(device.Stats.WriteBytes)), device.Stats.WriteCount)
		if device.Stats.UtilizationPct > 0 {
			info += fmt.Sprintf("  Utilization: %.1f%%\n", device.Stats.UtilizationPct)
		}
		info += "\n"
	}

	info += "Performance Indicators:\n"
	info += "• < 10 MB/s: Low activity\n"
	info += "• 10-100 MB/s: Moderate activity\n"
	info += "• > 100 MB/s: High activity\n"
	info += "• > 90% utilization: Disk bottleneck\n"

	return info
}
