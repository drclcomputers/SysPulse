package processes

import (
	"testing"
)

func TestCanKillProcess(t *testing.T) {
	canKill, reason := CanKillProcess(-1)
	if canKill {
		t.Errorf("Expected false for invalid PID, got true")
	}
	if reason == "" {
		t.Errorf("Expected error message for invalid PID")
	}

	canKill, reason = CanKillProcess(0)
	if canKill {
		t.Errorf("Expected false for zero PID, got true")
	}
	if reason == "" {
		t.Errorf("Expected error message for zero PID")
	}
}

func TestGetProcessKillMethods(t *testing.T) {
	methods := GetProcessKillMethods()
	if len(methods) == 0 {
		t.Errorf("Expected at least one kill method, got none")
	}

	for _, method := range methods {
		if method == "" {
			t.Errorf("Found empty kill method")
		}
	}
}

func TestKillProcByID(t *testing.T) {
	result := KillProcByID(-1)
	if result == "" {
		t.Errorf("Expected error message for invalid PID, got empty string")
	}

	result = KillProcByID(999999)
	if result == "" {
		t.Errorf("Expected error message for non-existent PID, got empty string")
	}
}

func TestForceKillProcByID(t *testing.T) {
	result := ForceKillProcByID(-1)
	if result == "" {
		t.Errorf("Expected error message for invalid PID, got empty string")
	}

	result = ForceKillProcByID(999999)
	if result == "" {
		t.Errorf("Expected error message for non-existent PID, got empty string")
	}
}
