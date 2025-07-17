//go:build darwin
// +build darwin

package load

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func GetDarwinLoadAverage() (*LoadAverage, error) {
	if load, err := getDarwinSysctlLoadAverage(); err == nil {
		return load, nil
	}

	if load, err := getDarwinUptimeLoadAverage(); err == nil {
		return load, nil
	}

	return nil, fmt.Errorf("unable to get load average on Darwin")
}

func getDarwinSysctlLoadAverage() (*LoadAverage, error) {
	cmd := exec.Command("sysctl", "-n", "vm.loadavg")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("sysctl vm.loadavg failed: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))
	outputStr = strings.Trim(outputStr, "{}")

	fields := strings.Fields(outputStr)
	if len(fields) < 3 {
		return nil, fmt.Errorf("unexpected sysctl output format: %s", outputStr)
	}

	load1, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 1min load: %v", err)
	}

	load5, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 5min load: %v", err)
	}

	load15, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 15min load: %v", err)
	}

	return &LoadAverage{
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}, nil
}

func getDarwinUptimeLoadAverage() (*LoadAverage, error) {
	cmd := exec.Command("uptime")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("uptime command failed: %v", err)
	}

	outputStr := string(output)

	index := strings.Index(outputStr, "load averages:")
	if index == -1 {
		return nil, fmt.Errorf("load averages not found in uptime output")
	}

	loadStr := outputStr[index+len("load averages:"):]
	fields := strings.Fields(loadStr)

	if len(fields) < 3 {
		return nil, fmt.Errorf("insufficient load average values in uptime output")
	}

	load1, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 1min load: %v", err)
	}

	load5, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 5min load: %v", err)
	}

	load15, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 15min load: %v", err)
	}

	return &LoadAverage{
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}, nil
}

func getDarwinSyscallLoadAverage() (*LoadAverage, error) {
	type loadavg struct {
		load  [3]uint32
		scale int32
	}

	var la loadavg

	r1, _, errno := syscall.Syscall(148,
		uintptr(unsafe.Pointer(&la.load[0])), 3, 0)

	if errno != 0 {
		return nil, fmt.Errorf("getloadavg syscall failed: %v", errno)
	}

	if r1 == 0 {
		return nil, fmt.Errorf("getloadavg returned no data")
	}

	scale := 65536.0

	return &LoadAverage{
		Load1:  float64(la.load[0]) / scale,
		Load5:  float64(la.load[1]) / scale,
		Load15: float64(la.load[2]) / scale,
	}, nil
}
