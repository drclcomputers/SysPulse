package network

import (
	"fmt"
	"syspulse/internal/utils"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/shirou/gopsutil/net"
)

var (
	lastBytesSent   uint64
	lastBytesRecv   uint64
	lastMeasurement time.Time
	bytesSentPerSec float64
	bytesRecvPerSec float64
)

func formatBytes(bytes float64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB/s", bytes/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB/s", bytes/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB/s", bytes/KB)
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

			sentBarCount := utils.Between(int(bytesSentPerSec/1024/10), 0, 10)
			recvBarCount := utils.Between(int(bytesRecvPerSec/1024/10), 0, 10)

			color := d.Theme.Network.BarLow
			if sentBarCount > 8 {
				color = d.Theme.Network.BarHigh
			}
			uploadText := fmt.Sprintf("Upload  : %s %s", utils.BarColor(utils.BAR, sentBarCount, color), formatBytes(bytesSentPerSec))
			currentY := utils.SafePrintText(screen, uploadText, x+2, y+1, w-2, h-1, utils.GetColorFromName(d.Theme.Layout.Network.ForegroundColor))

			color = d.Theme.Network.BarLow
			if recvBarCount > 8 {
				color = d.Theme.Network.BarHigh
			}
			downloadText := fmt.Sprintf("Download: %s %s", utils.BarColor(utils.BAR, recvBarCount, color), formatBytes(bytesRecvPerSec))
			currentY = utils.SafePrintText(screen, downloadText, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Network.ForegroundColor))

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

	info := "Network Activity Information\n\n"

	if len(stats) > 0 {
		netStat := stats[0]

		info += "Overall Network Statistics\n"
		info += fmt.Sprintf("Bytes Sent/Received: %.2f/%.2f GB\n", float64(netStat.BytesSent)/1024/1024/1024, float64(netStat.BytesRecv)/1024/1024/1024)
		info += fmt.Sprintf("Packets Sent/Received: %d/%d\n", netStat.PacketsSent, netStat.PacketsRecv)
		info += fmt.Sprintf("Send/Receive Errors: %d/%d\n", netStat.Errin, netStat.Errout)
		info += fmt.Sprintf("Dropped Packets In/Out: %d/%d\n", netStat.Dropin, netStat.Dropout)
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
				info += fmt.Sprintf("  Sent: %.2f MB (%d packets)\n", float64(iface.BytesSent)/1024/1024, iface.PacketsSent)
				info += fmt.Sprintf("  Received: %.2f MB (%d packets)\n", float64(iface.BytesRecv)/1024/1024, iface.PacketsRecv)
				if iface.Errin > 0 || iface.Errout > 0 {
					info += fmt.Sprintf("  Errors: %d in, %d out\n", iface.Errin, iface.Errout)
				}
				if iface.Dropin > 0 || iface.Dropout > 0 {
					info += fmt.Sprintf("  Drops: %d in, %d out\n", iface.Dropin, iface.Dropout)
				}
				info += "\n"
			}
		}
	}

	return info
}
