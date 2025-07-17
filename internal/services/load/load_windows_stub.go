//go:build !windows
// +build !windows

package load

import (
	"fmt"
	"runtime"
)

func GetWindowsLoadAverage() (*LoadAverage, error) {
	return nil, fmt.Errorf("Windows load average not supported on %s", runtime.GOOS)
}
