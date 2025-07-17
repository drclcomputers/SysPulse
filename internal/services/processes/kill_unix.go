//go:build linux || darwin || freebsd || openbsd || netbsd
// +build linux darwin freebsd openbsd netbsd

package processes

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/process"
)

func KillProcByID(pid int32) string {
	err := terminateProcessGracefully(pid)
	if err == nil {
		return ""
	}

	err = forceKillProcess(pid)
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

	time.Sleep(2 * time.Second)

	running, err := proc.IsRunning()
	if err != nil {
		return fmt.Errorf("failed to check if process is running: %v", err)
	}

	if running {
		return fmt.Errorf("process still running after SIGTERM")
	}

	return nil
}

func forceKillProcess(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}

	err = proc.Kill()
	if err != nil {
		return fmt.Errorf("force kill failed: %v", err)
	}

	return nil
}

func KillProcessWithSignal(pid int32, signal syscall.Signal) string {
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return fmt.Sprintf("failed to find process: %v", err)
	}

	err = process.Signal(signal)
	if err != nil {
		return fmt.Sprintf("failed to send signal %v: %v", signal, err)
	}

	return ""
}

func ForceKillProcByID(pid int32) string {
	return KillProcessWithSignal(pid, syscall.SIGKILL)
}

func KillProcessTree(pid int32) string {
	cmd := exec.Command("pkill", "-P", fmt.Sprintf("%d", pid))
	err := cmd.Run()
	if err != nil {
		return KillProcByID(pid)
	}

	return KillProcByID(pid)
}

func GetProcessKillMethods() []string {
	return []string{
		"Graceful (SIGTERM)",
		"Force Kill (SIGKILL)",
		"Interrupt (SIGINT)",
		"Hangup (SIGHUP)",
		"Kill Process Tree",
	}
}

func CanKillProcess(pid int32) (bool, string) {
	if pid <= 0 {
		return false, "Invalid PID"
	}

	if pid == 1 {
		return false, "Cannot kill init process"
	}

	if pid == 2 {
		return false, "Cannot kill kernel thread"
	}

	proc, err := process.NewProcess(pid)
	if err != nil {
		return false, "Process not found"
	}

	username, err := proc.Username()
	if err != nil {
		return false, "Cannot get process owner"
	}

	currentUser := os.Getenv("USER")
	if username == "root" && currentUser != "root" {
		return false, "Cannot kill root processes without root privileges"
	}

	name, err := proc.Name()
	if err != nil {
		return false, "Cannot get process name"
	}

	systemProcesses := []string{
		"systemd", "kthreadd", "ksoftirqd", "migration", "rcu_", "watchdog",
		"NetworkManager", "dbus", "ssh", "init",
	}

	for _, sysProc := range systemProcesses {
		if name == sysProc {
			return false, "Cannot kill system-critical process"
		}
	}

	return true, ""
}
