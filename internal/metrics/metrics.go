package metrics

import (
	"fmt"
	"sync"
	"time"
)

type MetricType string

const (
	CPUUpdate          MetricType = "cpu_update"
	MemoryUpdate       MetricType = "memory_update"
	DiskUpdate         MetricType = "disk_update"
	NetworkUpdate      MetricType = "network_update"
	ProcessUpdate      MetricType = "process_update"
	GPUUpdate          MetricType = "gpu_update"
	LoadUpdate         MetricType = "load_update"
	TemperatureUpdate  MetricType = "temperature_update"
	NetworkConnsUpdate MetricType = "network_conns_update"
	DiskIOUpdate       MetricType = "disk_io_update"
	ProcessTreeUpdate  MetricType = "process_tree_update"
	BatteryUpdate      MetricType = "battery_update"
)

type Metrics struct {
	mu              sync.RWMutex
	updateDurations map[MetricType][]time.Duration
	errorCounts     map[MetricType]int
	lastUpdate      map[MetricType]time.Time
	sampleWindow    time.Duration
}

func New(sampleWindow time.Duration) *Metrics {
	return &Metrics{
		updateDurations: make(map[MetricType][]time.Duration),
		errorCounts:     make(map[MetricType]int),
		lastUpdate:      make(map[MetricType]time.Time),
		sampleWindow:    sampleWindow,
	}
}

func (m *Metrics) RecordUpdateDuration(metricType MetricType, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.updateDurations[metricType] = append(m.updateDurations[metricType], duration)
	m.lastUpdate[metricType] = time.Now()

	m.cleanup(metricType)
}

func (m *Metrics) RecordError(metricType MetricType) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.errorCounts[metricType]++
}

func (m *Metrics) GetAverageUpdateDuration(metricType MetricType) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()

	durations := m.updateDurations[metricType]
	if len(durations) == 0 {
		return 0
	}

	var total time.Duration
	for _, d := range durations {
		total += d
	}

	return total / time.Duration(len(durations))
}

func (m *Metrics) GetErrorCount(metricType MetricType) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.errorCounts[metricType]
}

func (m *Metrics) GetLastUpdate(metricType MetricType) time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.lastUpdate[metricType]
}

func (m *Metrics) cleanup(metricType MetricType) {
	cutoff := time.Now().Add(-m.sampleWindow)
	var newDurations []time.Duration

	for i, t := range m.updateDurations[metricType] {
		if m.lastUpdate[metricType].Add(-time.Duration(len(m.updateDurations[metricType])-i) * time.Second).After(cutoff) {
			newDurations = append(newDurations, t)
		}
	}

	m.updateDurations[metricType] = newDurations
}

func (m *Metrics) GetStats() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := "Performance Metrics:\n"
	for _, metricType := range []MetricType{CPUUpdate, MemoryUpdate, DiskUpdate, NetworkUpdate, ProcessUpdate} {
		avg := m.GetAverageUpdateDuration(metricType)
		errors := m.GetErrorCount(metricType)
		last := m.GetLastUpdate(metricType)

		stats += fmt.Sprintf("%s:\n", metricType)
		stats += fmt.Sprintf("  Average Update Duration: %v\n", avg)
		stats += fmt.Sprintf("  Error Count: %d\n", errors)
		if !last.IsZero() {
			stats += fmt.Sprintf("  Last Update: %s\n", last.Format(time.RFC3339))
		}
	}

	return stats
}
