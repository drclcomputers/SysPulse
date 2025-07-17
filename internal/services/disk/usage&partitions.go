package disk

import (
	"fmt"
	"strings"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/shirou/gopsutil/disk"
)

func getDiskBar(used, total float64, theme utils.DISKModel) string {
	usedPercent := (used / total) * 100
	barWidth := 15
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
	emptyBar := strings.Repeat(utils.BAR, barWidth-usedWidth)

	return fmt.Sprintf("[%s]%s[-][%s]%s[-]", barColor, usedBar, theme.BarEmpty, emptyBar)
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
				used := float64(u.Used) / 1e9
				total := float64(u.Total) / 1e9
				//free := float64(u.Free) / 1e9
				bar := getDiskBar(used, total, d.Theme.Disk)

				if currentY >= y+h-1 {
					break
				}

				line1 := fmt.Sprintf("%s - Filesystem: %s",
					u.Path,
					partitions[i].Fstype,
				)
				currentY = utils.SafePrintText(screen, line1, x+3, currentY, w-6, y+h-1, utils.GetColorFromName(d.Theme.Layout.Disk.ForegroundColor))

				if currentY >= y+h-1 {
					break
				}

				line2 := fmt.Sprintf("  %s %.1f/%.1fGB",
					bar,
					used,
					total,
					//free,
				)
				currentY = utils.SafePrintText(screen, line2, x+3, currentY, w-6, y+h-1, utils.GetColorFromName(d.Theme.Layout.Disk.ForegroundColor))
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
