package processes

import (
	"fmt"
	"sort"
	"strings"
	"syspulse/internal/utils"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/shirou/gopsutil/process"
)

type ProcessNode struct {
	PID        int32          `json:"pid"`
	PPID       int32          `json:"ppid"`
	Name       string         `json:"name"`
	CPUPct     float64        `json:"cpu_percent"`
	Memory     uint64         `json:"memory"`
	Status     string         `json:"status"`
	CreateTime time.Time      `json:"create_time"`
	Children   []*ProcessNode `json:"children"`
}

type ProcessTree struct {
	Roots      []*ProcessNode `json:"roots"`
	TotalCount int            `json:"total_count"`
	LastUpdate time.Time      `json:"last_update"`
}

func GetProcessTree() (*ProcessTree, error) {
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}

	processMap := make(map[int32]*ProcessNode)
	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		name, _ := proc.Name()
		ppid, _ := proc.Ppid()
		cpuPct, _ := proc.CPUPercent()
		memInfo, _ := proc.MemoryInfo()
		status, _ := proc.Status()
		createTime, _ := proc.CreateTime()

		var memory uint64
		if memInfo != nil {
			memory = memInfo.RSS
		}

		node := &ProcessNode{
			PID:        pid,
			PPID:       ppid,
			Name:       name,
			CPUPct:     cpuPct,
			Memory:     memory,
			Status:     status,
			CreateTime: time.Unix(createTime/1000, 0),
			Children:   make([]*ProcessNode, 0),
		}

		processMap[pid] = node
	}

	roots := make([]*ProcessNode, 0)
	for _, node := range processMap {
		if parent, exists := processMap[node.PPID]; exists {
			parent.Children = append(parent.Children, node)
		} else {
			roots = append(roots, node)
		}
	}

	sort.Slice(roots, func(i, j int) bool {
		return roots[i].Name < roots[j].Name
	})

	sortChildren(roots)

	return &ProcessTree{
		Roots:      roots,
		TotalCount: len(processMap),
		LastUpdate: time.Now(),
	}, nil
}

func sortChildren(nodes []*ProcessNode) {
	for _, node := range nodes {
		sort.Slice(node.Children, func(i, j int) bool {
			return node.Children[i].Name < node.Children[j].Name
		})
		sortChildren(node.Children)
	}
}

func UpdateProcessTree(d *utils.Dashboard) {
	if d.ProcessTreeWidget == nil {
		return
	}

	tree, err := GetProcessTree()
	if err != nil {
		d.ProcessTreeWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
			utils.SafePrintText(screen, "Process tree unavailable", x+3, y+1, w-2, h-(y+1-y), tcell.ColorRed)
			return x, y, w, h
		})
		return
	}

	d.ProcessTreeData = tree
	d.ProcessTreeWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		currentY := y + 1
		maxY := y + h - 1

		for _, root := range tree.Roots {
			if currentY >= maxY {
				break
			}
			currentY = drawProcessNode(screen, root, "", x, currentY, w, maxY, true, d)
		}

		if currentY < maxY-1 {
			currentY = maxY - 1
			summary := fmt.Sprintf("Total: %d processes", tree.TotalCount)
			utils.SafePrintText(screen, summary, x+3, currentY, w-2, maxY-currentY, utils.GetColorFromName(d.Theme.Layout.ProcessTree.ForegroundColor))
		}

		return x, y, w, h
	})
}

func drawProcessNode(screen tcell.Screen, node *ProcessNode, prefix string, x, y, w, maxY int, isLast bool, d *utils.Dashboard) int {
	if y >= maxY {
		return y
	}

	truncateString := func(s string, maxLen int) string {
		if maxLen <= 3 {
			return "..."[:maxLen]
		}
		if len(s) > maxLen {
			return s[:maxLen-3] + "..."
		}
		return s
	}

	displayText := fmt.Sprintf("%s%s (PID: %d)", prefix, node.Name, node.PID)

	maxTextWidth := w - 4
	if maxTextWidth > 0 {
		displayText = truncateString(displayText, maxTextWidth)
	}

	utils.SafePrintText(screen, displayText, x+1, y, w-2, maxY-y, utils.GetColorFromName(d.Theme.Layout.ProcessTree.ForegroundColor))

	currentY := y + 1

	for i, child := range node.Children {
		if currentY >= maxY {
			break
		}

		childPrefix := prefix
		if isLast {
			childPrefix += "  "
		} else {
			childPrefix += "│ "
		}

		isLastChild := i == len(node.Children)-1
		if isLastChild {
			childPrefix += "└─"
		} else {
			childPrefix += "├─"
		}

		currentY = drawProcessNode(screen, child, childPrefix, x, currentY, w, maxY, isLastChild, d)
	}

	return currentY
}

func formatMemory(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	if bytes >= GB {
		return fmt.Sprintf("%.1fG", float64(bytes)/float64(GB))
	} else if bytes >= MB {
		return fmt.Sprintf("%.1fM", float64(bytes)/float64(MB))
	} else if bytes >= KB {
		return fmt.Sprintf("%.1fK", float64(bytes)/float64(KB))
	}
	return fmt.Sprintf("%dB", bytes)
}

func getProcessStatusColor(status string) string {
	switch strings.ToLower(status) {
	case "running":
		return "green"
	case "sleeping", "sleep":
		return "blue"
	case "stopped", "stop":
		return "yellow"
	case "zombie":
		return "red"
	case "idle":
		return "cyan"
	default:
		return "white"
	}
}

func GetProcessTreeFormattedInfo() string {
	tree, err := GetProcessTree()
	if err != nil {
		return fmt.Sprintf("Process Tree: Error - %v", err)
	}

	info := "Process Tree Structure\n\n"
	info += fmt.Sprintf("Total Processes: %d\n", tree.TotalCount)
	info += fmt.Sprintf("Last Update: %s\n\n", tree.LastUpdate.Format("15:04:05"))

	statusCounts := make(map[string]int)
	countProcessesByStatus(tree.Roots, statusCounts)

	info += "Process Status Summary:\n"
	for status, count := range statusCounts {
		info += fmt.Sprintf("  %s: %d\n", status, count)
	}

	info += "\nTop-level Processes:\n"
	for i, root := range tree.Roots {
		if i >= 10 {
			info += fmt.Sprintf("  ... and %d more\n", len(tree.Roots)-i)
			break
		}
		childCount := countTotalChildren(root)
		info += fmt.Sprintf("  %s (PID:%d) - %d children\n", root.Name, root.PID, childCount)
	}

	info += "\nProcess Status Legend:\n"
	info += "• Running: Currently executing\n"
	info += "• Sleeping: Waiting for resources\n"
	info += "• Stopped: Suspended process\n"
	info += "• Zombie: Terminated but not cleaned up\n"
	info += "• Idle: Waiting for work\n"

	return info
}

func countProcessesByStatus(nodes []*ProcessNode, statusCounts map[string]int) {
	for _, node := range nodes {
		statusCounts[node.Status]++
		countProcessesByStatus(node.Children, statusCounts)
	}
}

func countTotalChildren(node *ProcessNode) int {
	count := len(node.Children)
	for _, child := range node.Children {
		count += countTotalChildren(child)
	}
	return count
}

func GetProcessTreeSummary() string {
	tree, err := GetProcessTree()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	runningCount := 0
	sleepingCount := 0
	otherCount := 0

	var countByStatus func([]*ProcessNode)
	countByStatus = func(nodes []*ProcessNode) {
		for _, node := range nodes {
			switch strings.ToLower(node.Status) {
			case "running":
				runningCount++
			case "sleeping", "sleep":
				sleepingCount++
			default:
				otherCount++
			}
			countByStatus(node.Children)
		}
	}

	countByStatus(tree.Roots)

	return fmt.Sprintf("Total: %d | Running: %d | Sleeping: %d | Other: %d",
		tree.TotalCount, runningCount, sleepingCount, otherCount)
}
