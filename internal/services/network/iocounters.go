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
