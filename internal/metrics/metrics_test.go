package metrics

import (
	"testing"
	"time"
)

func TestMetrics(t *testing.T) {
	t.Run("New Metrics", func(t *testing.T) {
		m := New(time.Hour)
		if m == nil {
			t.Fatal("Expected non-nil Metrics instance")
		}
	})

	t.Run("Update Durations", func(t *testing.T) {
		m := New(time.Hour)
		duration := time.Second

		m.RecordUpdateDuration(CPUUpdate, duration)
		avg := m.GetAverageUpdateDuration(CPUUpdate)

		if avg != duration {
			t.Errorf("Expected average duration %v, got %v", duration, avg)
		}
	})

	t.Run("Error Counts", func(t *testing.T) {
		m := New(time.Hour)

		m.RecordError(MemoryUpdate)
		m.RecordError(MemoryUpdate)

		count := m.GetErrorCount(MemoryUpdate)
		if count != 2 {
			t.Errorf("Expected error count 2, got %d", count)
		}
	})

	t.Run("Last Update Time", func(t *testing.T) {
		m := New(time.Hour)
		before := time.Now()

		m.RecordUpdateDuration(NetworkUpdate, time.Second)
		after := time.Now()

		last := m.GetLastUpdate(NetworkUpdate)
		if last.Before(before) || last.After(after) {
			t.Error("Last update time outside expected range")
		}
	})

	t.Run("Cleanup", func(t *testing.T) {
		m := New(time.Millisecond)
		m.RecordUpdateDuration(DiskUpdate, time.Second)

		time.Sleep(time.Millisecond * 2)
		m.RecordUpdateDuration(DiskUpdate, time.Second)

		durations := len(m.updateDurations[DiskUpdate])
		if durations > 1 {
			t.Errorf("Expected at most 1 duration after cleanup, got %d", durations)
		}
	})

	t.Run("GetStats", func(t *testing.T) {
		m := New(time.Hour)
		m.RecordUpdateDuration(CPUUpdate, time.Second)
		m.RecordError(CPUUpdate)

		stats := m.GetStats()
		if stats == "" {
			t.Error("Expected non-empty stats string")
		}
	})
}
