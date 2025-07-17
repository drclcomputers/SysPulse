package network

import (
	"fmt"
	"sort"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/net"
)

type ConnectionStat struct {
	LocalAddr  string `json:"local_addr"`
	RemoteAddr string `json:"remote_addr"`
	Status     string `json:"status"`
	PID        int32  `json:"pid"`
	Family     uint32 `json:"family"`
	Type       uint32 `json:"type"`
}

type ConnectionStats struct {
	Connections []ConnectionStat  `json:"connections"`
	Summary     ConnectionSummary `json:"summary"`
}

type ConnectionSummary struct {
	Total       int `json:"total"`
	Established int `json:"established"`
	Listen      int `json:"listen"`
	TimeWait    int `json:"time_wait"`
	CloseWait   int `json:"close_wait"`
	SynSent     int `json:"syn_sent"`
	SynRecv     int `json:"syn_recv"`
	FinWait1    int `json:"fin_wait1"`
	FinWait2    int `json:"fin_wait2"`
	Closing     int `json:"closing"`
	LastAck     int `json:"last_ack"`
}

func GetNetworkConnections() (*ConnectionStats, error) {
	connections, err := net.Connections("all")
	if err != nil {
		return nil, err
	}

	stats := &ConnectionStats{
		Connections: make([]ConnectionStat, 0),
		Summary:     ConnectionSummary{},
	}

	for _, conn := range connections {
		connStat := ConnectionStat{
			LocalAddr:  fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port),
			RemoteAddr: fmt.Sprintf("%s:%d", conn.Raddr.IP, conn.Raddr.Port),
			Status:     conn.Status,
			PID:        conn.Pid,
			Family:     conn.Family,
			Type:       conn.Type,
		}

		stats.Connections = append(stats.Connections, connStat)
		stats.Summary.Total++

		switch conn.Status {
		case "ESTABLISHED":
			stats.Summary.Established++
		case "LISTEN":
			stats.Summary.Listen++
		case "TIME_WAIT":
			stats.Summary.TimeWait++
		case "CLOSE_WAIT":
			stats.Summary.CloseWait++
		case "SYN_SENT":
			stats.Summary.SynSent++
		case "SYN_RECV":
			stats.Summary.SynRecv++
		case "FIN_WAIT1":
			stats.Summary.FinWait1++
		case "FIN_WAIT2":
			stats.Summary.FinWait2++
		case "CLOSING":
			stats.Summary.Closing++
		case "LAST_ACK":
			stats.Summary.LastAck++
		}
	}

	sort.Slice(stats.Connections, func(i, j int) bool {
		if stats.Connections[i].Status == stats.Connections[j].Status {
			return stats.Connections[i].LocalAddr < stats.Connections[j].LocalAddr
		}
		return getStatusPriority(stats.Connections[i].Status) < getStatusPriority(stats.Connections[j].Status)
	})

	return stats, nil
}

func getStatusPriority(status string) int {
	switch status {
	case "ESTABLISHED":
		return 1
	case "LISTEN":
		return 2
	case "TIME_WAIT":
		return 3
	case "CLOSE_WAIT":
		return 4
	case "SYN_SENT":
		return 5
	case "SYN_RECV":
		return 6
	default:
		return 7
	}
}

func UpdateNetworkConnections(d *utils.Dashboard) {
	if d.NetworkConnsWidget == nil {
		return
	}

	connStats, err := GetNetworkConnections()
	if err != nil {
		d.NetworkConnsWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
			utils.SafePrintText(screen, "Network connections unavailable", x+3, y+1, w-2, h-(y+1-y), tcell.ColorRed)
			return x, y, w, h
		})
		return
	}

	d.NetworkConnsData = connStats
	d.NetworkConnsWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		currentY := y + 1

		totalText := fmt.Sprintf("Total: %d", connStats.Summary.Total)
		currentY = utils.SafePrintText(screen, totalText, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.NetworkConns.ForegroundColor))

		if connStats.Summary.Established > 0 {
			text := fmt.Sprintf("Established: [green]%d[-]", connStats.Summary.Established)
			currentY = utils.SafePrintText(screen, text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.NetworkConns.ForegroundColor))
		}

		if connStats.Summary.Listen > 0 {
			text := fmt.Sprintf("Listening: [blue]%d[-]", connStats.Summary.Listen)
			currentY = utils.SafePrintText(screen, text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.NetworkConns.ForegroundColor))
		}

		if connStats.Summary.TimeWait > 0 {
			text := fmt.Sprintf("Time Wait: [yellow]%d[-]", connStats.Summary.TimeWait)
			currentY = utils.SafePrintText(screen, text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.NetworkConns.ForegroundColor))
		}

		if connStats.Summary.CloseWait > 0 {
			text := fmt.Sprintf("Close Wait: [orange]%d[-]", connStats.Summary.CloseWait)
			currentY = utils.SafePrintText(screen, text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.NetworkConns.ForegroundColor))
		}

		if currentY < y+h-3 {
			currentY = utils.SafePrintText(screen, "Recent Connections:", x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.NetworkConns.ForegroundColor))

			maxConn := (h - (currentY - y) - 1)
			for i, conn := range connStats.Connections {
				if i >= maxConn || currentY >= y+h-1 {
					break
				}

				color := getConnectionColor(conn.Status)
				localAddr := truncateAddr(conn.LocalAddr, 15)
				remoteAddr := truncateAddr(conn.RemoteAddr, 15)

				text := fmt.Sprintf("[%s]%s[-] %s→%s", color, conn.Status, localAddr, remoteAddr)
				currentY = utils.SafePrintText(screen, text, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.NetworkConns.ForegroundColor))
			}
		}

		return x, y, w, h
	})
}

func getConnectionColor(status string) string {
	switch status {
	case "ESTABLISHED":
		return "green"
	case "LISTEN":
		return "blue"
	case "TIME_WAIT":
		return "yellow"
	case "CLOSE_WAIT":
		return "orange"
	case "SYN_SENT", "SYN_RECV":
		return "cyan"
	default:
		return "gray"
	}
}

func truncateAddr(addr string, maxLen int) string {
	if len(addr) <= maxLen {
		return addr
	}
	return addr[:maxLen-3] + "..."
}

func GetNetworkConnectionsFormattedInfo() string {
	connStats, err := GetNetworkConnections()
	if err != nil {
		return fmt.Sprintf("Network connections: Error - %v", err)
	}

	info := "Network Connection Summary\n\n"
	info += fmt.Sprintf("Total Connections: %d\n", connStats.Summary.Total)
	info += fmt.Sprintf("Established: %d\n", connStats.Summary.Established)
	info += fmt.Sprintf("Listening: %d\n", connStats.Summary.Listen)
	info += fmt.Sprintf("Time Wait: %d\n", connStats.Summary.TimeWait)
	info += fmt.Sprintf("Close Wait: %d\n", connStats.Summary.CloseWait)

	if connStats.Summary.SynSent > 0 {
		info += fmt.Sprintf("SYN Sent: %d\n", connStats.Summary.SynSent)
	}
	if connStats.Summary.SynRecv > 0 {
		info += fmt.Sprintf("SYN Received: %d\n", connStats.Summary.SynRecv)
	}
	if connStats.Summary.FinWait1 > 0 {
		info += fmt.Sprintf("FIN Wait 1: %d\n", connStats.Summary.FinWait1)
	}
	if connStats.Summary.FinWait2 > 0 {
		info += fmt.Sprintf("FIN Wait 2: %d\n", connStats.Summary.FinWait2)
	}

	info += "\nConnection Details:\n"
	for i, conn := range connStats.Connections {
		if i >= 20 {
			info += fmt.Sprintf("... and %d more connections\n", len(connStats.Connections)-20)
			break
		}
		info += fmt.Sprintf("• %s: %s → %s", conn.Status, conn.LocalAddr, conn.RemoteAddr)
		if conn.PID > 0 {
			info += fmt.Sprintf(" (PID: %d)", conn.PID)
		}
		info += "\n"
	}

	return info
}

func CreateNetworkConnectionsModal(app *tview.Application, returnWidget tview.Primitive) tview.Primitive {
	connStats, err := GetNetworkConnections()
	if err != nil {
		return tview.NewModal().
			SetText(fmt.Sprintf("Error getting network connections: %v", err)).
			AddButtons([]string{"Ok"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				app.SetRoot(returnWidget, true).SetFocus(returnWidget)
			})
	}

	list := tview.NewList().
		SetHighlightFullLine(true).
		SetSelectedTextColor(tcell.ColorWhite).
		SetSelectedBackgroundColor(tcell.ColorDarkBlue)

	list.AddItem(fmt.Sprintf("Total Connections: %d", connStats.Summary.Total), "", 0, nil)
	list.AddItem(fmt.Sprintf("Established: %d", connStats.Summary.Established), "", 0, nil)
	list.AddItem(fmt.Sprintf("Listening: %d", connStats.Summary.Listen), "", 0, nil)
	list.AddItem(fmt.Sprintf("Time Wait: %d", connStats.Summary.TimeWait), "", 0, nil)
	list.AddItem(fmt.Sprintf("Close Wait: %d", connStats.Summary.CloseWait), "", 0, nil)

	if connStats.Summary.SynSent > 0 {
		list.AddItem(fmt.Sprintf("SYN Sent: %d", connStats.Summary.SynSent), "", 0, nil)
	}
	if connStats.Summary.SynRecv > 0 {
		list.AddItem(fmt.Sprintf("SYN Received: %d", connStats.Summary.SynRecv), "", 0, nil)
	}
	if connStats.Summary.FinWait1 > 0 {
		list.AddItem(fmt.Sprintf("FIN Wait 1: %d", connStats.Summary.FinWait1), "", 0, nil)
	}
	if connStats.Summary.FinWait2 > 0 {
		list.AddItem(fmt.Sprintf("FIN Wait 2: %d", connStats.Summary.FinWait2), "", 0, nil)
	}

	list.AddItem("", "", 0, nil)
	list.AddItem("--- Connection Details ---", "", 0, nil)

	for _, conn := range connStats.Connections {
		connectionLine := fmt.Sprintf("[%s]%s[white] %s → %s",
			getConnectionColor(conn.Status),
			conn.Status,
			conn.LocalAddr,
			conn.RemoteAddr)
		if conn.PID > 0 {
			connectionLine += fmt.Sprintf(" (PID: %d)", conn.PID)
		}
		list.AddItem(connectionLine, "", 0, nil)
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.SetRoot(returnWidget, true).SetFocus(returnWidget)
			return nil
		}
		return event
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Network Connections - Press ESC to close").SetTextAlign(tview.AlignCenter), 1, 0, false).
		AddItem(list, 0, 1, true)

	frame := tview.NewFrame(flex).
		SetBorders(1, 1, 1, 1, 2, 2)

	return frame
}
