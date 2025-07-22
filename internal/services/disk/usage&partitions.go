package disk

import (
	"fmt"
	"strings"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/disk"
)

func getDiskBar(used, total float64, theme utils.DISKModel, w int) string {
	usedPercent := (used / total) * 100
	barWidth := w / 3
	usedWidth := int((usedPercent / 100) * float64(barWidth))

	var barColor string
	switch {
	case usedPercent >= 90:
		barColor = theme.BarHigh
	case usedPercent >= 70:
		barColor = theme.BarMedium
	default:
		barColor = theme.BarLow
	}

	usedBar := strings.Repeat(utils.BAR, usedWidth)
	emptyBar := strings.Repeat("░", barWidth-usedWidth)

	return fmt.Sprintf("[%s]%s[-]%s", barColor, usedBar, emptyBar)
}

func UpdateDisk(d *utils.Dashboard) {
	if d.DiskWidget == nil {
		return
	}

	if partitions, err := disk.Partitions(false); err == nil {
		d.DiskData = nil
		for _, p := range partitions {
			if usage, err := disk.Usage(p.Mountpoint); err == nil {
				d.DiskData = append(d.DiskData, usage)
			}
		}
		d.DiskWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
			currentY := y + 1
			for i, u := range d.DiskData {
				used := float64(u.Used) / 1024 / 1024 / 1024
				total := float64(u.Total) / 1024 / 1024 / 1024
				//free := float64(u.Free) / 1024 / 1024 / 1024
				bar := getDiskBar(used, total, d.Theme.Disk, w)

				if currentY >= y+h-1 {
					break
				}

				fs := ""
				if i < len(partitions) && partitions[i].Fstype != "" {
					fs = partitions[i].Fstype
				} else {
					fs = "Unknown"
				}

				line1 := fmt.Sprintf("%s (%s) %s",
					u.Path,
					fs,
					bar,
				)
				tview.Print(screen, line1, x+2, currentY, w-2, y+h-1, utils.GetColorFromName(d.Theme.Layout.Disk.ForegroundColor))

				currentY++
				if currentY >= y+h-1 {
					break
				}

				line2 := fmt.Sprintf("%.1f/%.1fGB",
					used,
					total,
				)

				tview.Print(screen, line2, x+3, currentY, w-2, y+h-1, utils.GetColorFromName(d.Theme.Layout.Disk.ForegroundColor))
				currentY++
			}
			return x, y, w, h
		})
	}
}

func GetNumberofPartitions() string {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return "Unknown"
	}

	return fmt.Sprint(len(partitions), " Partitions")
}

func GetDiskFormattedInfo() string {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return fmt.Sprintf("Disk: Error - %v", err)
	}

	info := "Disk Usage Information\n\n"

	totalUsed := uint64(0)
	totalSize := uint64(0)

	for i, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		info += fmt.Sprintf("=== Partition %d ===\n", i+1)
		info += fmt.Sprintf("Device: %s\n", partition.Device)
		info += fmt.Sprintf("Mountpoint: %s\n", partition.Mountpoint)
		info += fmt.Sprintf("Filesystem: %s\n", partition.Fstype)
		info += fmt.Sprintf("Used/Total: %.2f/%.2f GB (%.1f%%)\n", float64(usage.Used)/1024/1024/1024, float64(usage.Total)/1024/1024/1024, usage.UsedPercent)
		info += fmt.Sprintf("Free: %.2f GB\n", float64(usage.Free)/1024/1024/1024)
		//info += fmt.Sprintf("Inodes Total: %d\n", usage.InodesTotal)
		//info += fmt.Sprintf("Inodes Used: %d (%.1f%%)\n", usage.InodesUsed, usage.InodesUsedPercent)
		//info += fmt.Sprintf("Inodes Free: %d\n", usage.InodesFree)

		/*if usage.UsedPercent < 70 {
			info += "Status: Good - plenty of space available\n"
		} else if usage.UsedPercent < 85 {
			info += "Status: Moderate - consider cleaning up files\n"
		} else if usage.UsedPercent < 95 {
			info += "Status: Warning - running low on space\n"
		} else {
			info += "Status: Critical - almost full!\n"
		}*/

		totalUsed += usage.Used
		totalSize += usage.Total
		info += "\n"
	}

	if totalSize > 0 {
		overallPercent := (float64(totalUsed) / float64(totalSize)) * 100
		info += "=== Overall System ===\n"
		info += fmt.Sprintf("Used/Total Storage: %.2f/%.2f GB (%.1f%%)\n", float64(totalUsed)/1024/1024/1024, float64(totalSize)/1024/1024/1024, overallPercent)
		info += fmt.Sprintf("Free Storage: %.2f GB\n", float64(totalSize-totalUsed)/1024/1024/1024)
		info += "\n"
	}

	return info
}
