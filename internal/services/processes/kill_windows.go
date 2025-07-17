//go:build windows
// +build windows

package processes

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/shirou/gopsutil/process"
)

func KillProcByID(pid int32) string {
	err := terminateProcessGracefully(pid)
	if err == nil {
		return ""
	}

	err = terminateProcessWithTaskkill(pid)
	if err == nil {
		return ""
	}

	err = terminateProcessWithAPI(pid)
	if err != nil {
		return fmt.Sprintf("failed to kill process: %v", err)
	}

	return ""
}

func terminateProcessGracefully(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}

	err = proc.Terminate()
	if err != nil {
		return fmt.Errorf("graceful termination failed: %v", err)
	}

	return nil
}

func terminateProcessWithTaskkill(pid int32) error {
	cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("taskkill failed: %v", err)
	}

	return nil
}

func terminateProcessWithAPI(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}

	err = proc.Kill()
	if err != nil {
		return fmt.Errorf("Windows API kill failed: %v", err)
	}

	return nil
}

func ForceKillProcByID(pid int32) string {
	cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("force kill failed: %v", err)
	}

	return ""
}

func GetProcessKillMethods() []string {
	return []string{
		"Graceful (Terminate)",
		"Taskkill",
		"Force Kill",
		"API Kill",
	}
}

func CanKillProcess(pid int32) (bool, string) {
	if pid <= 0 {
		return false, "Invalid PID"
	}

	if pid <= 4 {
		return false, "Cannot kill system processes"
	}

	proc, err := process.NewProcess(pid)
	if err != nil {
		return false, "Process not found"
	}

	name, err := proc.Name()
	if err != nil {
		return false, "Cannot get process name"
	}

	systemProcesses := []string{
		"System", "Registry", "smss.exe", "csrss.exe", "winlogon.exe",
		"services.exe", "lsass.exe", "svchost.exe", "dwm.exe", "wininit.exe",
	}

	for _, sysProc := range systemProcesses {
		if name == sysProc {
			return false, "Cannot kill Windows system process"
		}
	}

	return true, ""
}
