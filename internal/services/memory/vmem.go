package memory

import (
	"fmt"
	"strings"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/mem"
)

func getMemoryBar(used, total float64, barColor string, d *utils.Dashboard) string {
	usedPercent := 0.0
	if total != 0 {
		usedPercent = (used / total) * 100
	}
	barWidth := 15
	usedWidth := int((usedPercent / 100) * float64(barWidth))

	usedBar := strings.Repeat(utils.BAR, usedWidth)
	emptyBar := strings.Repeat("░", barWidth-usedWidth)

	return fmt.Sprintf("[%s]%s[-][%s]%s[-]", barColor, usedBar, utils.GetColorFromName(d.Theme.Foreground), emptyBar)
}

func VMem() {
	vm, _ := mem.VirtualMemory()
	fmt.Printf("Used: %v B, Total: %v B, Usage: %.2f%%\n", vm.Used, vm.Total, vm.UsedPercent)
}

func UpdateVMem(d *utils.Dashboard) {
	if d.MemWidget == nil {
		return
	}

	if VmemStat, err := mem.VirtualMemory(); err == nil {
		d.VMemData = VmemStat
		d.SMemData = GetSwapMem()
		if d.SMemData == nil {
			d.SMemData = &mem.SwapMemoryStat{
				Total: 0,
				Free:  0,
				Used:  0,
			}
		}

		d.MemWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
			VMemusedGB := float64(d.VMemData.Used) / 1024 / 1024 / 1024
			VMemtotalGB := float64(d.VMemData.Total) / 1024 / 1024 / 1024
			VMembar := getMemoryBar(VMemusedGB, VMemtotalGB, d.Theme.Memory.VMemGauge, d)

			SMemusedGB := float64(d.SMemData.Used) / 1024 / 1024 / 1024
			SMemtotalGB := float64(d.SMemData.Total) / 1024 / 1024 / 1024
			SMembar := getMemoryBar(SMemusedGB, SMemtotalGB, d.Theme.Memory.SMemGauge, d)

			currentY := 3
			vMemText := fmt.Sprintf("RAM : %s %.1f/%.1fGB", VMembar, VMemusedGB, VMemtotalGB)
			tview.Print(screen, vMemText, x+2, currentY, w-2, h-1, utils.GetColorFromName(d.Theme.Layout.Memory.ForegroundColor))
			currentY += 2

			sMemText := fmt.Sprintf("Swap: %s %.1f/%.1fGB", SMembar, SMemusedGB, SMemtotalGB)
			tview.Print(screen, sMemText, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Memory.ForegroundColor))
			return x, y, w, h
		})
	}
}

func GetRAM() string {
	vm, _ := mem.VirtualMemory()
	return fmt.Sprintf("%.1f GB", float64(vm.Total)/1024/1024/1024)
}

func GetMemoryFormattedInfo() string {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Sprintf("Memory: Error - %v", err)
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		return fmt.Sprintf("Memory: Error getting swap - %v", err)
	}

	info := "Memory Usage Information\n\n"

	info += "=== RAM (Virtual Memory) ===\n"
	info += fmt.Sprintf("Total: %.2f GB\n", float64(vm.Total)/1024/1024/1024)
	info += fmt.Sprintf("Used: %.2f GB (%.1f%%)\n", float64(vm.Used)/1024/1024/1024, vm.UsedPercent)
	info += fmt.Sprintf("Free: %.2f GB\n", float64(vm.Free)/1024/1024/1024)
	info += fmt.Sprintf("Available: %.2f GB (%.1f%%)\n", float64(vm.Available)/1024/1024/1024, (float64(vm.Available)/float64(vm.Total))*100)
	info += fmt.Sprintf("Cached: %.2f GB\n", float64(vm.Cached)/1024/1024/1024)
	info += fmt.Sprintf("Buffers: %.2f GB\n", float64(vm.Buffers)/1024/1024/1024)

	info += "\n=== Swap Memory ===\n"
	info += fmt.Sprintf("Total: %.2f GB\n", float64(swap.Total)/1024/1024/1024)
	info += fmt.Sprintf("Used: %.2f GB (%.1f%%)\n", float64(swap.Used)/1024/1024/1024, swap.UsedPercent)
	info += fmt.Sprintf("Free: %.2f GB\n", float64(swap.Free)/1024/1024/1024)

	info += "\n=== Memory Health ===\n"
	if vm.UsedPercent < 70 {
		info += "• Memory Status: Good - plenty of available memory\n"
	} else if vm.UsedPercent < 85 {
		info += "• Memory Status: Moderate - consider closing unused applications\n"
	} else {
		info += "• Memory Status: High - system may be running low on memory\n"
	}

	if swap.UsedPercent > 50 {
		info += "• Swap Usage: High - system is relying heavily on swap memory\n"
	} else if swap.UsedPercent > 0 {
		info += "• Swap Usage: Some swap in use - this is normal\n"
	} else {
		info += "• Swap Usage: None - system has sufficient RAM\n"
	}

	return info
}
