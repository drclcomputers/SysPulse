package network

import (
	"fmt"
	"strings"
	"syspulse/internal/utils"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/net"
)

var (
	lastBytesSent   uint64
	lastBytesRecv   uint64
	lastMeasurement time.Time
	bytesSentPerSec float64
	bytesRecvPerSec float64
)

func getNetworkBar(bytesPerSec float64, barColor string, d *utils.Dashboard, w int) string {
	maxSpeed := 100.0 * 1024 * 1024
	percentage := (bytesPerSec / maxSpeed) * 100
	if percentage > 100 {
		percentage = 100
	}

	barWidth := w / 3
	if barWidth < 10 {
		barWidth = 10
	}
	usedWidth := int((percentage / 100) * float64(barWidth))

	usedBar := strings.Repeat(utils.BAR, usedWidth)
	emptyBar := strings.Repeat("░", barWidth-usedWidth)

	return fmt.Sprintf("[%s]%s[-][%s]%s[-]", barColor, usedBar, utils.GetColorFromName(d.Theme.Foreground), emptyBar)
}

func formatBytes(bytes float64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.0f GB/s", bytes/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.0f MB/s", bytes/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.0f KB/s", bytes/KB)
	default:
		return fmt.Sprintf("%.0f B/s", bytes)
	}
}

func UpdateNetwork(d *utils.Dashboard) {
	if d.NetWidget == nil {
		return
	}

	d.NetWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		stats, err := net.IOCounters(false)
		if err == nil && len(stats) > 0 {
			d.NetData = &stats[0]

			now := time.Now()
			if !lastMeasurement.IsZero() {
				duration := now.Sub(lastMeasurement).Seconds()
				if duration > 0 {
					bytesSentPerSec = float64(d.NetData.BytesSent-lastBytesSent) / duration
					bytesRecvPerSec = float64(d.NetData.BytesRecv-lastBytesRecv) / duration
				}
			}

			lastBytesSent = d.NetData.BytesSent
			lastBytesRecv = d.NetData.BytesRecv
			lastMeasurement = now

			uploadColor := d.Theme.Network.BarLow
			if bytesSentPerSec > 10*1024*1024 {
				uploadColor = d.Theme.Network.BarHigh
			}

			downloadColor := d.Theme.Network.BarLow
			if bytesRecvPerSec > 10*1024*1024 {
				downloadColor = d.Theme.Network.BarHigh
			}

			uploadBar := getNetworkBar(bytesSentPerSec, uploadColor, d, w)
			downloadBar := getNetworkBar(bytesRecvPerSec, downloadColor, d, w)

			uploadText := fmt.Sprintf("Upload  : %s %s", uploadBar, formatBytes(bytesSentPerSec))
			tview.Print(screen, uploadText, x+2, y+1, w-2, h-1, utils.GetColorFromName(d.Theme.Layout.Network.ForegroundColor))
			currentY := y + 2

			downloadText := fmt.Sprintf("Download: %s %s", downloadBar, formatBytes(bytesRecvPerSec))
			tview.Print(screen, downloadText, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Network.ForegroundColor))
			currentY++

			for _, iface := range GetInterfaces() {
				if currentY >= y+h-1 {
					break
				}
				currentY = utils.SafePrintText(screen, iface, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Network.ForegroundColor))
			}
		} else {
			utils.SafePrintText(screen, "Network stats unavailable", x+2, y+2, w-2, h-1, utils.GetColorFromName(d.Theme.Layout.Network.ForegroundColor))
		}
		return x, y, w, h
	})
}

func GetNetworkFormattedInfo() string {
	stats, err := net.IOCounters(false)
	if err != nil {
		return fmt.Sprintf("Network: Error - %v", err)
	}

	var info string

	if len(stats) > 0 {
		netStat := stats[0]

		info += "Overall Network Statistics\n"
		info += fmt.Sprintf("• Bytes Sent/Received: %.2f/%.2f GB\n", float64(netStat.BytesSent)/1024/1024/1024, float64(netStat.BytesRecv)/1024/1024/1024)
		info += fmt.Sprintf("• Packets Sent/Received: %d/%d\n", netStat.PacketsSent, netStat.PacketsRecv)
		info += fmt.Sprintf("• Send/Receive Errors: %d/%d\n", netStat.Errin, netStat.Errout)
		info += fmt.Sprintf("• Dropped Packets In/Out: %d/%d\n", netStat.Dropin, netStat.Dropout)
		info += "\n"
	}

	info += fmt.Sprintf("Up/Down Speed: %s / %s\n", formatBytes(bytesSentPerSec), formatBytes(bytesRecvPerSec))
	info += "\n"

	info += "Network Interfaces\n"
	interfaces := GetInterfaces()
	for _, iface := range interfaces {
		info += fmt.Sprintf("• %s\n", iface)
	}

	interfaceStats, err := net.IOCounters(true)
	if err == nil && len(interfaceStats) > 0 {
		info += "\nPer-Interface Statistics\n"
		for _, iface := range interfaceStats {
			if iface.BytesSent > 0 || iface.BytesRecv > 0 {
				info += fmt.Sprintf("Interface: %s\n", iface.Name)
				info += fmt.Sprintf("• Sent: %.2f MB (%d packets)\n", float64(iface.BytesSent)/1024/1024, iface.PacketsSent)
				info += fmt.Sprintf("• Received: %.2f MB (%d packets)\n", float64(iface.BytesRecv)/1024/1024, iface.PacketsRecv)
				if iface.Errin > 0 || iface.Errout > 0 {
					info += fmt.Sprintf("• Errors: %d in, %d out\n", iface.Errin, iface.Errout)
				}
				if iface.Dropin > 0 || iface.Dropout > 0 {
					info += fmt.Sprintf("• Drops: %d in, %d out\n", iface.Dropin, iface.Dropout)
				}
				info += "\n"
			}
		}
	}

	return info
}
