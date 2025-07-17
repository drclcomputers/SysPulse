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
			utils.SafePrintText(screen, sMemText, x+2, currentY+1, w-2, h-(currentY-y), utils.GetColorFromName(d.Theme.Layout.Memory.ForegroundColor))
			return x, y, w, h
		})
	}
}

func GetRAM() string {
	vm, _ := mem.VirtualMemory()
	return fmt.Sprintf("%.1f GB", float64(vm.Total)/1024/1024/1024)
}
