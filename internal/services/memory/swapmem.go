package memory

import (
	"github.com/shirou/gopsutil/mem"
)

func GetSwapMem() *mem.SwapMemoryStat {
	swap, err := mem.SwapMemory()
	if err != nil {
		return nil
	}
	return swap
}
