package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"syspulse/internal/utils"
)

type DataPoint struct {
	Timestamp time.Time
	CPU       []float64
	Memory    struct {
		Total     uint64
		Used      uint64
		SwapTotal uint64
		SwapUsed  uint64
	}
	Disk struct {
		Path     string
		Total    uint64
		Used     uint64
		UsedPerc float64
		IOReads  uint64
		IOWrites uint64
	}
	Network struct {
		BytesSent     uint64
		BytesReceived uint64
		PacketsSent   uint64
		PacketsRecv   uint64
	}
	Load struct {
		Load1  float64
		Load5  float64
		Load15 float64
	}
	Temperature struct {
		CPUTemp float64
		GPUTemp float64
	}
	NetworkConnections struct {
		Total       int
		Established int
		Listening   int
		CloseWait   int
		TimeWait    int
	}
	DiskIO struct {
		ReadCount  uint64
		WriteCount uint64
		ReadBytes  uint64
		WriteBytes uint64
	}
	ProcessTree struct {
		ProcessCount int
		TopProcesses []string
	}
	Battery struct {
		Level         float64
		Status        string
		IsCharging    bool
		TimeRemaining string
	}
	GPU []struct {
		Name        string  `json:"name"`
		Vendor      string  `json:"vendor"`
		MemoryTotal uint64  `json:"memory_total"`
		MemoryUsed  uint64  `json:"memory_used"`
		MemoryFree  uint64  `json:"memory_free"`
		Temperature float64 `json:"temperature"`
		Usage       float64 `json:"usage"`
		Available   bool    `json:"available"`
	} `json:"gpu"`
}

type ExportFormat int

const (
	CSV ExportFormat = iota
	JSON
)

func ExportData(data []DataPoint, filename string, format ExportFormat) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create export directory: %v", err)
	}

	switch format {
	case CSV:
		return exportToCSV(data, filename)
	case JSON:
		return exportToJSON(data, filename)
	default:
		return fmt.Errorf("unsupported export format")
	}
}

func exportToCSV(data []DataPoint, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"Timestamp",
		"CPU_Total",
		"Memory_Total", "Memory_Used",
		"Swap_Total", "Swap_Used",
		"Disk_Path", "Disk_Total", "Disk_Used", "Disk_UsedPerc",
		"Disk_IOReads", "Disk_IOWrites",
		"Net_BytesSent", "Net_BytesReceived",
		"Net_PacketsSent", "Net_PacketsReceived",
		"Load_1", "Load_5", "Load_15",
		"Temp_CPU", "Temp_GPU",
		"NetConn_Total", "NetConn_Established", "NetConn_Listening",
		"DiskIO_ReadCount", "DiskIO_WriteCount", "DiskIO_ReadBytes", "DiskIO_WriteBytes",
		"Processes_Count", "Processes_Top",
		"Battery_Level", "Battery_Status", "Battery_Charging", "Battery_TimeRemaining",
		"GPU_Count", "GPU_Primary_Name", "GPU_Primary_Vendor", "GPU_Primary_MemoryTotal", "GPU_Primary_MemoryUsed", "GPU_Primary_Usage",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, d := range data {
		cpuTotal := 0.0
		for _, cpu := range d.CPU {
			cpuTotal += cpu
		}
		cpuTotal /= float64(len(d.CPU))

		gpuCount := len(d.GPU)
		primaryGPUName := ""
		primaryGPUVendor := ""
		primaryGPUMemoryTotal := uint64(0)
		primaryGPUMemoryUsed := uint64(0)
		primaryGPUUsage := 0.0

		if gpuCount > 0 {
			primaryGPUName = d.GPU[0].Name
			primaryGPUVendor = d.GPU[0].Vendor
			primaryGPUMemoryTotal = d.GPU[0].MemoryTotal
			primaryGPUMemoryUsed = d.GPU[0].MemoryUsed
			primaryGPUUsage = d.GPU[0].Usage
		}

		row := []string{
			d.Timestamp.Format(time.RFC3339),
			fmt.Sprintf("%.2f", cpuTotal),
			fmt.Sprintf("%d", d.Memory.Total),
			fmt.Sprintf("%d", d.Memory.Used),
			fmt.Sprintf("%d", d.Memory.SwapTotal),
			fmt.Sprintf("%d", d.Memory.SwapUsed),
			d.Disk.Path,
			fmt.Sprintf("%d", d.Disk.Total),
			fmt.Sprintf("%d", d.Disk.Used),
			fmt.Sprintf("%.2f", d.Disk.UsedPerc),
			fmt.Sprintf("%d", d.Disk.IOReads),
			fmt.Sprintf("%d", d.Disk.IOWrites),
			fmt.Sprintf("%d", d.Network.BytesSent),
			fmt.Sprintf("%d", d.Network.BytesReceived),
			fmt.Sprintf("%d", d.Network.PacketsSent),
			fmt.Sprintf("%d", d.Network.PacketsRecv),
			fmt.Sprintf("%.2f", d.Load.Load1),
			fmt.Sprintf("%.2f", d.Load.Load5),
			fmt.Sprintf("%.2f", d.Load.Load15),
			fmt.Sprintf("%.2f", d.Temperature.CPUTemp),
			fmt.Sprintf("%.2f", d.Temperature.GPUTemp),
			fmt.Sprintf("%d", d.NetworkConnections.Total),
			fmt.Sprintf("%d", d.NetworkConnections.Established),
			fmt.Sprintf("%d", d.NetworkConnections.Listening),
			fmt.Sprintf("%d", d.DiskIO.ReadCount),
			fmt.Sprintf("%d", d.DiskIO.WriteCount),
			fmt.Sprintf("%d", d.DiskIO.ReadBytes),
			fmt.Sprintf("%d", d.DiskIO.WriteBytes),
			fmt.Sprintf("%d", d.ProcessTree.ProcessCount),
			fmt.Sprintf("%v", d.ProcessTree.TopProcesses),
			fmt.Sprintf("%.2f", d.Battery.Level),
			d.Battery.Status,
			fmt.Sprintf("%t", d.Battery.IsCharging),
			d.Battery.TimeRemaining,
			fmt.Sprintf("%d", gpuCount),
			primaryGPUName,
			primaryGPUVendor,
			fmt.Sprintf("%d", primaryGPUMemoryTotal),
			fmt.Sprintf("%d", primaryGPUMemoryUsed),
			fmt.Sprintf("%.2f", primaryGPUUsage),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func exportToJSON(data []DataPoint, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func CreateSnapshot(d *utils.Dashboard) DataPoint {
	dp := DataPoint{
		Timestamp: time.Now(),
		CPU:       d.CpuData,
	}

	if d.VMemData != nil {
		dp.Memory.Total = d.VMemData.Total
		dp.Memory.Used = d.VMemData.Used
	}

	if d.SMemData != nil {
		dp.Memory.SwapTotal = d.SMemData.Total
		dp.Memory.SwapUsed = d.SMemData.Used
	}

	if len(d.DiskData) > 0 {
		dp.Disk.Path = d.DiskData[0].Path
		dp.Disk.Total = d.DiskData[0].Total
		dp.Disk.Used = d.DiskData[0].Used
		dp.Disk.UsedPerc = d.DiskData[0].UsedPercent
	}

	if d.NetData != nil {
		dp.Network.BytesSent = d.NetData.BytesSent
		dp.Network.BytesReceived = d.NetData.BytesRecv
		dp.Network.PacketsSent = d.NetData.PacketsSent
		dp.Network.PacketsRecv = d.NetData.PacketsRecv
	}

	if d.LoadData != nil {
		if loadData, ok := d.LoadData.(map[string]interface{}); ok {
			if load1, ok := loadData["load1"].(float64); ok {
				dp.Load.Load1 = load1
			}
			if load5, ok := loadData["load5"].(float64); ok {
				dp.Load.Load5 = load5
			}
			if load15, ok := loadData["load15"].(float64); ok {
				dp.Load.Load15 = load15
			}
		}
	}

	if d.TemperatureData != nil {
		if tempData, ok := d.TemperatureData.(map[string]interface{}); ok {
			if cpuTemp, ok := tempData["cpu_temp"].(float64); ok {
				dp.Temperature.CPUTemp = cpuTemp
			}
			if gpuTemp, ok := tempData["gpu_temp"].(float64); ok {
				dp.Temperature.GPUTemp = gpuTemp
			}
		}
	}

	if d.NetworkConnsData != nil {
		if connData, ok := d.NetworkConnsData.(map[string]interface{}); ok {
			if summary, ok := connData["summary"].(map[string]interface{}); ok {
				if total, ok := summary["total"].(int); ok {
					dp.NetworkConnections.Total = total
				}
				if established, ok := summary["established"].(int); ok {
					dp.NetworkConnections.Established = established
				}
				if listening, ok := summary["listening"].(int); ok {
					dp.NetworkConnections.Listening = listening
				}
			}
		}
	}

	if d.DiskIOData != nil {
		if diskIOData, ok := d.DiskIOData.(map[string]interface{}); ok {
			if readCount, ok := diskIOData["read_count"].(uint64); ok {
				dp.DiskIO.ReadCount = readCount
			}
			if writeCount, ok := diskIOData["write_count"].(uint64); ok {
				dp.DiskIO.WriteCount = writeCount
			}
			if readBytes, ok := diskIOData["read_bytes"].(uint64); ok {
				dp.DiskIO.ReadBytes = readBytes
			}
			if writeBytes, ok := diskIOData["write_bytes"].(uint64); ok {
				dp.DiskIO.WriteBytes = writeBytes
			}
		}
	}

	if d.ProcessTreeData != nil {
		if processData, ok := d.ProcessTreeData.(map[string]interface{}); ok {
			if processCount, ok := processData["process_count"].(int); ok {
				dp.ProcessTree.ProcessCount = processCount
			}
			if topProcesses, ok := processData["top_processes"].([]string); ok {
				dp.ProcessTree.TopProcesses = topProcesses
			}
		}
	}

	if d.BatteryData != nil {
		if batteryData, ok := d.BatteryData.(map[string]interface{}); ok {
			if level, ok := batteryData["level"].(float64); ok {
				dp.Battery.Level = level
			}
			if status, ok := batteryData["status"].(string); ok {
				dp.Battery.Status = status
			}
			if isCharging, ok := batteryData["is_charging"].(bool); ok {
				dp.Battery.IsCharging = isCharging
			}
			if timeRemaining, ok := batteryData["time_remaining"].(string); ok {
				dp.Battery.TimeRemaining = timeRemaining
			}
		}
	}

	if d.GPUData != nil {
		if gpuData, ok := d.GPUData.([]interface{}); ok && len(gpuData) > 0 {
			for _, gpuInterface := range gpuData {
				if gpuMap, ok := gpuInterface.(map[string]interface{}); ok {
					var gpuInfo struct {
						Name        string  `json:"name"`
						Vendor      string  `json:"vendor"`
						MemoryTotal uint64  `json:"memory_total"`
						MemoryUsed  uint64  `json:"memory_used"`
						MemoryFree  uint64  `json:"memory_free"`
						Temperature float64 `json:"temperature"`
						Usage       float64 `json:"usage"`
						Available   bool    `json:"available"`
					}

					if name, ok := gpuMap["name"].(string); ok {
						gpuInfo.Name = name
					}
					if vendor, ok := gpuMap["vendor"].(string); ok {
						gpuInfo.Vendor = vendor
					}
					if memoryTotal, ok := gpuMap["memory_total"].(uint64); ok {
						gpuInfo.MemoryTotal = memoryTotal
					}
					if memoryUsed, ok := gpuMap["memory_used"].(uint64); ok {
						gpuInfo.MemoryUsed = memoryUsed
					}
					if memoryFree, ok := gpuMap["memory_free"].(uint64); ok {
						gpuInfo.MemoryFree = memoryFree
					}
					if temperature, ok := gpuMap["temperature"].(float64); ok {
						gpuInfo.Temperature = temperature
					}
					if usage, ok := gpuMap["usage"].(float64); ok {
						gpuInfo.Usage = usage
					}
					if available, ok := gpuMap["available"].(bool); ok {
						gpuInfo.Available = available
					}

					dp.GPU = append(dp.GPU, gpuInfo)
				}
			}
		}
	}

	return dp
}
