package ui

import (
	"fmt"
	"strings"
	"syspulse/internal/services/processes"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/process"
)

func (d *Dashboard) showHelpModal() {
	helpText := `Global Shortcuts:
- TAB/Shift+TAB - Cycle through widgets
- Q - Quit application
- H - Show this help screen
- I - Show more information
- Y - Change process sorting (CPU/Memory)

Quick Navigation:
- C - Focus CPU widget
- M - Focus Memory widget
- D - Focus Disk widget
- N - Focus Network widget
- P - Focus Process widget
- G - Focus GPU widget

Widget Information:
- I or ENTER - Show detailed information modal for focused widget

Process Management:
- K - Kill selected process (platform-specific methods)
- F - Search/filter processes
- Up/Down or W/S - Navigate process list
- I - View process details

Process Kill Methods:
- Kill - Graceful termination (recommended)
- Force Kill - Immediate termination
- Platform-specific optimizations for Windows/Linux
`

	modal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"Close"}).
		SetBackgroundColor(tcell.ColorBlack).
		SetTextColor(tcell.ColorWhite).
		SetButtonBackgroundColor(tcell.ColorBlue).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			d.App.SetRoot(d.MainWidget, true)
		})

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(modal, 0, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	d.App.SetRoot(flex, true).SetFocus(modal)
}

func (d *Dashboard) showProcessSearch() {
	form := tview.NewForm().
		SetButtonsAlign(tview.AlignCenter).
		SetFieldTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetButtonBackgroundColor(tcell.ColorBlue)

	searchInput := tview.NewInputField().
		SetLabel("Search term: ").
		SetFieldWidth(30).
		SetPlaceholder("Enter process name, PID, or CPU/MEM %")

	originalProcesses := make([]string, 0)
	for i := 0; i < d.ProcessWidget.GetItemCount(); i++ {
		text, _ := d.ProcessWidget.GetItemText(i)
		originalProcesses = append(originalProcesses, text)
	}

	filterState := &struct {
		filterType string
	}{
		filterType: "all",
	}

	form.AddDropDown("Filter by:", []string{"All", "High CPU (>50%)", "High Memory (>50%)", "System Processes", "User Processes"}, 0,
		func(option string, index int) {
			switch index {
			case 0:
				filterState.filterType = "all"
			case 1:
				filterState.filterType = "highcpu"
			case 2:
				filterState.filterType = "highmem"
			case 3:
				filterState.filterType = "system"
			case 4:
				filterState.filterType = "user"
			}
			applyProcessFilter(d, searchInput.GetText(), filterState.filterType, originalProcesses)
		})

	searchInput.SetChangedFunc(func(text string) {
		applyProcessFilter(d, text, filterState.filterType, originalProcesses)
	})

	form.AddFormItem(searchInput)
	form.AddButton("Apply", func() {
		d.App.SetRoot(d.MainWidget, true).SetFocus(d.ProcessWidget)
	}).
		AddButton("Reset", func() {
			d.ProcessFilterActive = false
			d.ProcessFilterTerm = ""
			d.ProcessFilterType = "all"

			d.ProcessWidget.Clear()
			for _, proc := range originalProcesses {
				d.ProcessWidget.AddItem(proc, "", 0, nil)
			}
			d.App.SetRoot(d.MainWidget, true).SetFocus(d.ProcessWidget)
		}).
		AddButton("Close", func() {
			d.App.SetRoot(d.MainWidget, true).SetFocus(d.ProcessWidget)
		})

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 0, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	d.App.SetRoot(flex, true).SetFocus(searchInput)
}

func applyProcessFilter(d *Dashboard, searchTerm, filterType string, originalProcesses []string) {
	d.ProcessFilterActive = true
	d.ProcessFilterTerm = searchTerm
	d.ProcessFilterType = filterType

	d.ProcessWidget.Clear()
	for _, proc := range originalProcesses {
		var cpu, mem float64
		var pid int32

		pidStart := strings.Index(proc, "(PID: ")
		if pidStart != -1 {
			pidEnd := strings.Index(proc[pidStart:], ")")
			if pidEnd != -1 {
				pidStr := proc[pidStart+6 : pidStart+pidEnd]
				fmt.Sscanf(pidStr, "%d", &pid)
			}

			cpuStart := strings.Index(proc, "CPU: ")
			if cpuStart != -1 {
				cpuEnd := strings.Index(proc[cpuStart:], "%")
				if cpuEnd != -1 {
					cpuStr := proc[cpuStart+5 : cpuStart+cpuEnd]
					fmt.Sscanf(cpuStr, "%f", &cpu)
				}
			}

			memStart := strings.Index(proc, "MEM: ")
			if memStart != -1 {
				memEnd := strings.Index(proc[memStart:], "%")
				if memEnd != -1 {
					memStr := proc[memStart+5 : memStart+memEnd]
					fmt.Sscanf(memStr, "%f", &mem)
				}
			}
		}

		show := true
		switch filterType {
		case "highcpu":
			show = cpu > 50
		case "highmem":
			show = mem > 50
		case "system":
			show = pid < 1000
		case "user":
			show = pid >= 1000
		}

		if show && (searchTerm == "" || utils.CaseInsensitiveContains(proc, searchTerm)) {
			d.ProcessWidget.AddItem(proc, "", 0, nil)
		}
	}
}

func (d *Dashboard) quitModal() {
	modal := tview.NewModal().
		SetText("Are you sure you want to quit?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				d.App.Stop()
			} else {
				d.App.SetRoot(d.MainWidget, true)
			}
		})
	d.App.SetRoot(modal, false).SetFocus(modal)
}

func (d *Dashboard) showProcessKillModal(selectedPID int32) {
	canKill, reason := processes.CanKillProcess(selectedPID)
	if !canKill {
		modal := tview.NewModal().
			SetText(fmt.Sprintf("Cannot kill process PID %d: %s", selectedPID, reason)).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				d.App.SetRoot(d.MainWidget, true).SetFocus(d.ProcessWidget)
			})
		d.App.SetRoot(modal, false).SetFocus(modal)
		return
	}

	proc, err := process.NewProcess(selectedPID)
	var procName string
	if err == nil {
		procName, _ = proc.Name()
	}

	displayText := fmt.Sprintf("Kill process PID: %d", selectedPID)
	if procName != "" {
		displayText = fmt.Sprintf("Kill process: %s (PID: %d)", procName, selectedPID)
	}

	modal := tview.NewModal().
		SetText(displayText).
		AddButtons([]string{"Kill", "Force Kill", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Kill":
				result := processes.KillProcByID(selectedPID)
				if result != "" {
					d.showKillResultModal(selectedPID, result)
					return
				}
			case "Force Kill":
				result := processes.ForceKillProcByID(selectedPID)
				if result != "" {
					d.showKillResultModal(selectedPID, result)
					return
				}
			}
			d.App.SetRoot(d.MainWidget, true).SetFocus(d.ProcessWidget)
		})
	d.App.SetRoot(modal, false).SetFocus(modal)
}

func (d *Dashboard) showKillResultModal(pid int32, errorMsg string) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Failed to kill process PID %d:\n%s", pid, errorMsg)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			d.App.SetRoot(d.MainWidget, true).SetFocus(d.ProcessWidget)
		})
	d.App.SetRoot(modal, false).SetFocus(modal)
}

func (d *Dashboard) showProcessTreeModal() {
	tree, err := processes.GetProcessTree()
	if err != nil {
		modal := tview.NewModal().
			SetText(fmt.Sprintf("Error loading process tree: %v", err)).
			AddButtons([]string{"Ok"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				d.App.SetRoot(d.MainWidget, true).SetFocus(d.ProcessTreeWidget)
			})
		d.App.SetRoot(modal, false).SetFocus(modal)
		return
	}

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)

	textView.SetBorder(true).
		SetTitle("Process Tree (Arrow keys to scroll, ESC to close)").
		SetTitleAlign(tview.AlignCenter)

	content := "[yellow]Process Tree Structure[white]\n\n"
	content += fmt.Sprintf("Total Processes: %d\n", tree.TotalCount)
	content += fmt.Sprintf("Last Update: %s\n\n", tree.LastUpdate.Format("15:04:05"))

	statusCounts := make(map[string]int)
	countProcessesByStatus(tree.Roots, statusCounts)

	content += "[blue]Process Status Summary:[white]\n"
	for status, count := range statusCounts {
		content += fmt.Sprintf("  %s: %d\n", status, count)
	}
	content += "\n"

	content += "[blue]Process Tree:[white]\n"
	for _, root := range tree.Roots {
		content += buildProcessTreeString(root, "", true)
	}

	content += "\n[blue]Process Status Legend:[white]\n"
	content += "- [green]Running:[white] Currently executing\n"
	content += "- [blue]Sleeping:[white] Waiting for resources\n"
	content += "- [yellow]Stopped:[white] Suspended process\n"
	content += "- [red]Zombie:[white] Terminated but not cleaned up\n"
	content += "- [cyan]Idle:[white] Waiting for work\n"

	textView.SetText(content)

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			d.App.SetRoot(d.MainWidget, true).SetFocus(d.ProcessTreeWidget)
			return nil
		}

		switch event.Rune() {
		case 'q', 'Q':
			d.App.SetRoot(d.MainWidget, true).SetFocus(d.ProcessTreeWidget)
			return nil
		}

		return event
	})

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(textView, 0, 10, true).
			AddItem(nil, 0, 1, false), 0, 10, true).
		AddItem(nil, 0, 1, false)

	d.App.SetRoot(flex, true).SetFocus(textView)
}

func countProcessesByStatus(nodes []*processes.ProcessNode, statusCounts map[string]int) {
	for _, node := range nodes {
		statusCounts[node.Status]++
		countProcessesByStatus(node.Children, statusCounts)
	}
}

func buildProcessTreeString(node *processes.ProcessNode, prefix string, isLast bool) string {
	result := ""

	statusColor := getProcessStatusColor(node.Status)
	memoryStr := formatMemoryBytes(node.Memory)

	result += fmt.Sprintf("%s[%s]%s[white] (PID:%d) - CPU:%.1f%% MEM:%s [%s]%s[white]\n",
		prefix, statusColor, node.Name, node.PID, node.CPUPct, memoryStr, statusColor, node.Status)

	childPrefix := prefix
	if isLast {
		childPrefix += "  "
	} else {
		childPrefix += "│ "
	}

	for i, child := range node.Children {
		isLastChild := i == len(node.Children)-1
		if isLastChild {
			result += buildProcessTreeString(child, childPrefix+"└─", true)
		} else {
			result += buildProcessTreeString(child, childPrefix+"├─", false)
		}
	}

	return result
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

func formatMemoryBytes(bytes uint64) string {
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
