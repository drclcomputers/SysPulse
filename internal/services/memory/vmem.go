package memory

import (
	"fmt"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/shirou/gopsutil/mem"
)

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
		d.MemWidget.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
			VMemusedGB := float64(d.VMemData.Used) / 1e9
			VMemtotalGB := float64(d.VMemData.Total) / 1e9
			VMembarCount := int((VMemusedGB / VMemtotalGB) * 10)

			SMemusedGB := float64(d.SMemData.Used) / 1e9
			SMemtotalGB := float64(d.SMemData.Total) / 1e9
			SMembarCount := int((SMemusedGB / SMemtotalGB) * 10)

			vMemText := fmt.Sprintf("Virtual (RAM): %s%.1f/%.1fGB", utils.BarColor("█", VMembarCount, d.Theme.Memory.VMemGauge), VMemusedGB, VMemtotalGB)
			currentY := utils.SafePrintText(screen, vMemText, x+2, y+1, w-2, h-1, utils.GetColorFromName(d.Theme.Layout.Memory.ForegroundColor))

			sMemText := fmt.Sprintf("Swap Memory  : %s%.1f/%.1fGB", utils.BarColor("█", SMembarCount, d.Theme.Memory.SMemGauge), SMemusedGB, SMemtotalGB)
			utils.SafePrintText(screen, sMemText, x+2, currentY, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Memory.ForegroundColor))
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
