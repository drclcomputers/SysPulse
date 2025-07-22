package processes

import (
	"fmt"
	"sort"
	"strings"
	"syspulse/internal/utils"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/process"
)

func applyProcessFilterToItems(items []string, searchTerm, filterType string) []string {
	if searchTerm == "" && filterType == "all" {
		return items
	}

	var filteredItems []string
	for _, item := range items {
		var cpu, mem float64
		var pid int32

		pidStart := strings.Index(item, "(PID: ")
		if pidStart != -1 {
			pidEnd := strings.Index(item[pidStart:], ")")
			if pidEnd != -1 {
				pidStr := item[pidStart+6 : pidStart+pidEnd]
				fmt.Sscanf(pidStr, "%d", &pid)
			}

			cpuStart := strings.Index(item, "CPU: ")
			if cpuStart != -1 {
				cpuEnd := strings.Index(item[cpuStart:], "%")
				if cpuEnd != -1 {
					cpuStr := item[cpuStart+5 : cpuStart+cpuEnd]
					fmt.Sscanf(cpuStr, "%f", &cpu)
				}
			}

			memStart := strings.Index(item, "MEM: ")
			if memStart != -1 {
				memEnd := strings.Index(item[memStart:], "%")
				if memEnd != -1 {
					memStr := item[memStart+5 : memStart+memEnd]
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

		if show && (searchTerm == "" || utils.CaseInsensitiveContains(item, searchTerm)) {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func GetNrProcesses() int {
	procs, err := process.Processes()
	if err != nil {
		panic(err)
	}
	return len(procs)
}

func UpdateProcesses(d *utils.Dashboard) {
	if d.ProcessWidget == nil {
		return
	}

	var selectedPID int32
	currentItem := d.ProcessWidget.GetCurrentItem()
	if currentItem >= 0 && currentItem < d.ProcessWidget.GetItemCount() {
		text, _ := d.ProcessWidget.GetItemText(currentItem)
		fmt.Sscanf(text, "%s (PID: %d)", new(string), &selectedPID)
	}

	procs, err := process.Processes()
	if err != nil {
		return
	}

	if d.Theme.Sorting == "mem" {
		sort.Slice(procs, func(i, j int) bool {
			mem1, _ := procs[i].MemoryPercent()
			mem2, _ := procs[j].MemoryPercent()
			return mem1 > mem2
		})
	} else {
		sort.Slice(procs, func(i, j int) bool {
			cpu1, _ := procs[i].CPUPercent()
			cpu2, _ := procs[j].CPUPercent()
			return cpu1 > cpu2
		})
	}

	var items []string
	selectedIndex := 0

	for i, p := range procs {
		pid := p.Pid
		name, _ := p.Name()
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryPercent()

		mainText := fmt.Sprintf("%s (PID: %d) - CPU: %.1f%% | MEM: %.1f%%", name, pid, cpu, mem)
		items = append(items, mainText)

		if pid == selectedPID {
			selectedIndex = i
		}
	}

	if d.ProcessFilterActive {
		items = applyProcessFilterToItems(items, d.ProcessFilterTerm, d.ProcessFilterType)
		selectedIndex = 0
		for i, item := range items {
			var pid int32
			fmt.Sscanf(item, "%s (PID: %d)", new(string), &pid)
			if pid == selectedPID {
				selectedIndex = i
				break
			}
		}
	}

	if d.ProcessWidget.GetItemCount() != len(items) {
		d.ProcessWidget.Clear()
		for _, item := range items {
			d.ProcessWidget.AddItem(item, "", 0, nil)
		}
	} else {
		for i, item := range items {
			currentText, _ := d.ProcessWidget.GetItemText(i)
			if currentText != item {
				d.ProcessWidget.SetItemText(i, item, "")
			}
		}
	}

	if selectedPID != 0 {
		if currentItem != selectedIndex {
			d.ProcessWidget.SetCurrentItem(selectedIndex)
		}
	}
}

func ShowProcessDetails(d *utils.Dashboard) {
	currentItem := d.ProcessWidget.GetCurrentItem()
	if currentItem < 0 || currentItem >= d.ProcessWidget.GetItemCount() {
		return
	}

	text, _ := d.ProcessWidget.GetItemText(currentItem)
	var selectedPID int32
	fmt.Sscanf(text, "%s (PID: %d)", new(string), &selectedPID)

	proc, err := process.NewProcess(selectedPID)
	if err != nil {
		return
	}

	name, _ := proc.Name()
	cmdline, _ := proc.Cmdline()
	cpu, _ := proc.CPUPercent()
	mem, _ := proc.MemoryPercent()
	status, _ := proc.Status()
	createTime, _ := proc.CreateTime()
	numThreads, _ := proc.NumThreads()
	username, _ := proc.Username()
	memInfo, _ := proc.MemoryInfo()

	details := fmt.Sprintf(`Process Details

Basic Information:
• Name: %s
• PID: %d
• Status: %s
• Username: %s
• Created: %s

Resource Usage:
• CPU Usage: %.2f%%
• Memory Usage: %.2f%%
• Memory RSS: %d MB
• Memory VMS: %d MB
• Threads: %d

Command:
%s

[yellow]Navigation: ↑/↓ to scroll, Esc to close[white]`,
		name, selectedPID, status, username,
		time.Unix(createTime/1000, 0).Format("2006-01-02 15:04:05"),
		cpu, mem,
		memInfo.RSS/1024/1024,
		memInfo.VMS/1024/1024,
		numThreads,
		cmdline)

	modal := tview.NewTextView().
		SetText(details).
		SetScrollable(true).
		SetWrap(true).
		SetDynamicColors(true).
		SetWordWrap(true)

	utils.SetBorderStyle(modal.Box)
	modal.SetTitle("Process Details (Arrow keys to scroll, ESC to close)").
		SetTitleAlign(tview.AlignCenter)

	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			d.App.SetRoot(d.MainWidget, true).SetFocus(d.ProcessWidget)
			return nil
		}
		return event
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
