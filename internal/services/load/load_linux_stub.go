//go:build !linux
// +build !linux

package load

import (
	"fmt"
	"runtime"
)

func GetLinuxLoadAverage() (*LoadAverage, error) {
	return nil, fmt.Errorf("Linux load average not supported on %s", runtime.GOOS)
}
