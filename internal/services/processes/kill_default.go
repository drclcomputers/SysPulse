//go:build !windows && !linux && !darwin && !freebsd && !openbsd && !netbsd
// +build !windows,!linux,!darwin,!freebsd,!openbsd,!netbsd

package processes

import (
	"fmt"

	"github.com/shirou/gopsutil/process"
)

func KillProcByID(pid int32) string {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Sprintf("failed to find process: %v", err)
	}

	err = proc.Kill()
	if err != nil {
		return fmt.Sprintf("failed to kill process: %v", err)
	}

	return ""
}

func ForceKillProcByID(pid int32) string {
	return KillProcByID(pid)
}

func GetProcessKillMethods() []string {
	return []string{
		"Basic Kill",
	}
}

func CanKillProcess(pid int32) (bool, string) {
	if pid <= 0 {
		return false, "Invalid PID"
	}

	proc, err := process.NewProcess(pid)
	if err != nil {
		return false, "Process not found"
	}

	_, err = proc.Name()
	if err != nil {
		return false, "Cannot access process"
	}

	return true, ""
}
