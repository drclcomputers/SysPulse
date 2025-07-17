//go:build windows
// +build windows

package load

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type win32_PerfRawData_PerfOS_Processor struct {
	Name                 string
	PercentProcessorTime uint64
	Timestamp_Sys100NS   uint64
	Frequency_Sys100NS   uint64
}

type win32_PerfRawData_PerfOS_System struct {
	ProcessorQueueLength  uint32
	Processes             uint32
	Threads               uint32
	SystemCallsPerSec     uint64
	ContextSwitchesPerSec uint64
}

var (
	lastCPUTimes   []float64
	lastUpdateTime time.Time
	loadSamples1   []float64
	loadSamples5   []float64
	loadSamples15  []float64
	sampleIndex    int
	isInitialized  bool
)

func GetWindowsLoadAverage() (*LoadAverage, error) {
	if !isInitialized {
		initializeLoadTracking()
	}

	cpuUsage, err := getCurrentCPUUsage()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	}

	queueLength, err := getProcessorQueueLength()
	if err != nil {
		queueLength = 0
	}

	currentLoad := calculateLoadFromCPU(cpuUsage, queueLength)

	updateLoadSamples(currentLoad)

	return &LoadAverage{
		Load1:  calculateAverage(loadSamples1),
		Load5:  calculateAverage(loadSamples5),
		Load15: calculateAverage(loadSamples15),
	}, nil
}

func initializeLoadTracking() {
	loadSamples1 = make([]float64, 12)
	loadSamples5 = make([]float64, 60)
	loadSamples15 = make([]float64, 180)
	sampleIndex = 0
	lastUpdateTime = time.Now()
	isInitialized = true

	if cpuUsage, err := getCurrentCPUUsage(); err == nil {
		initialLoad := calculateLoadFromCPU(cpuUsage, 0)
		for i := range loadSamples1 {
			loadSamples1[i] = initialLoad
		}
		for i := range loadSamples5 {
			loadSamples5[i] = initialLoad
		}
		for i := range loadSamples15 {
			loadSamples15[i] = initialLoad
		}
	}
}

func getCurrentCPUUsage() (float64, error) {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}

	if len(percentages) > 0 {
		return percentages[0], nil
	}

	return 0, fmt.Errorf("no CPU data available")
}

func getProcessorQueueLength() (float64, error) {
	return 0.0, fmt.Errorf("processor queue length monitoring temporarily disabled")
}

func calculateLoadFromCPU(cpuUsage float64, queueLength float64) float64 {
	cpuCount := float64(runtime.NumCPU())

	baseLoad := (cpuUsage / 100.0) * cpuCount

	queueContribution := queueLength / cpuCount

	totalLoad := baseLoad + queueContribution

	if totalLoad > cpuCount*2 {
		totalLoad = cpuCount * 2
	}

	return totalLoad
}

func updateLoadSamples(currentLoad float64) {
	now := time.Now()

	if now.Sub(lastUpdateTime) < 4*time.Second {
		return
	}

	loadSamples1[sampleIndex%len(loadSamples1)] = currentLoad

	loadSamples5[sampleIndex%len(loadSamples5)] = currentLoad

	loadSamples15[sampleIndex%len(loadSamples15)] = currentLoad

	sampleIndex++
	lastUpdateTime = now
}

func calculateAverage(samples []float64) float64 {
	if len(samples) == 0 {
		return 0
	}

	var sum float64
	var count int

	for _, sample := range samples {
		if sample > 0 {
			sum += sample
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return sum / float64(count)
}
