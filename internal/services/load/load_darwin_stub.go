//go:build !darwin
// +build !darwin

package load

import (
	"fmt"
	"runtime"
)

func GetDarwinLoadAverage() (*LoadAverage, error) {
	return nil, fmt.Errorf("Darwin load average not supported on %s", runtime.GOOS)
}
