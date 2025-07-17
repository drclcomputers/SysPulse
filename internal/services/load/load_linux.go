//go:build linux
// +build linux

package load

import (
	"github.com/shirou/gopsutil/load"
)

func GetLinuxLoadAverage() (*LoadAverage, error) {
	avg, err := load.Avg()
	if err != nil {
		return nil, err
	}

	return &LoadAverage{
		Load1:  avg.Load1,
		Load5:  avg.Load5,
		Load15: avg.Load15,
	}, nil
}
