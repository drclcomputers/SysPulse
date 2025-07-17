package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"syspulse/internal/export"
	ui "syspulse/internal/services/UI"
	"syspulse/internal/services/battery"
	"syspulse/internal/services/disk"
	"syspulse/internal/services/gpu"
	"syspulse/internal/services/load"
	"syspulse/internal/services/memory"
	"syspulse/internal/services/network"
	"syspulse/internal/services/sysinfo"
	"syspulse/internal/services/temperature"

	"github.com/spf13/cobra"
)

var (
	exportFormat    string
	exportOutput    string
	exportDirectory string
	exportDuration  int
	exportInterval  int
	exportSamples   int
	exportAll       bool
	exportQuiet     bool
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export system metrics to CSV or JSON format",
	Long: `Export system metrics to CSV or JSON format without opening the UI.
This command allows you to collect system metrics and export them directly
to files for analysis or integration with other tools.

Examples:
  syspulse export --format csv --output metrics.csv
  syspulse export --format json --output metrics.json --samples 10
  syspulse export --format csv --directory exports --duration 60
  syspulse export --format json --samples 5 --interval 2`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runExport(); err != nil {
			fmt.Fprintf(os.Stderr, "Export failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "csv", "Export format (csv, json)")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output filename (default: auto-generated)")
	exportCmd.Flags().StringVarP(&exportDirectory, "directory", "d", "exports", "Export directory")
	exportCmd.Flags().IntVar(&exportDuration, "duration", 0, "Collection duration in seconds (0 = single snapshot)")
	exportCmd.Flags().IntVar(&exportInterval, "interval", 1, "Collection interval in seconds")
	exportCmd.Flags().IntVar(&exportSamples, "samples", 1, "Number of samples to collect")
	exportCmd.Flags().BoolVar(&exportAll, "all", false, "Export both CSV and JSON formats")
	exportCmd.Flags().BoolVarP(&exportQuiet, "quiet", "q", false, "Quiet mode - minimal output")
}

func runExport() error {
	if !exportAll && exportFormat != "csv" && exportFormat != "json" {
		return fmt.Errorf("invalid format: %s (must be 'csv' or 'json')", exportFormat)
	}

	dashboard := ui.NewDashboard()

	if !exportQuiet {
		fmt.Printf("Collecting system metrics...\n")
	}

	var dataPoints []export.DataPoint

	if exportDuration > 0 {
		if !exportQuiet {
			fmt.Printf("Collection duration: %d seconds\n", exportDuration)
			fmt.Printf("Collection interval: %d seconds\n", exportInterval)
		}

		endTime := time.Now().Add(time.Duration(exportDuration) * time.Second)
		sampleCount := 0

		for time.Now().Before(endTime) {
			sysinfo.UpdateCPU(dashboard)
			memory.UpdateVMem(dashboard)
			disk.UpdateDisk(dashboard)
			network.UpdateNetwork(dashboard)
			battery.UpdateBatteryStatus(dashboard)
			temperature.UpdateTemperatures(dashboard)
			gpu.UpdateGPU(dashboard)
			load.UpdateLoadAverage(dashboard)

			snapshot := export.CreateSnapshot(dashboard)
			dataPoints = append(dataPoints, snapshot)
			sampleCount++

			if !exportQuiet {
				fmt.Printf("Collected sample %d\n", sampleCount)
			}

			time.Sleep(time.Duration(exportInterval) * time.Second)
		}
	} else {
		if !exportQuiet {
			fmt.Printf("Collecting %d samples with %d second intervals\n", exportSamples, exportInterval)
		}

		for i := 0; i < exportSamples; i++ {
			sysinfo.UpdateCPU(dashboard)
			memory.UpdateVMem(dashboard)
			disk.UpdateDisk(dashboard)
			network.UpdateNetwork(dashboard)
			battery.UpdateBatteryStatus(dashboard)
			temperature.UpdateTemperatures(dashboard)
			gpu.UpdateGPU(dashboard)
			load.UpdateLoadAverage(dashboard)

			snapshot := export.CreateSnapshot(dashboard)
			dataPoints = append(dataPoints, snapshot)

			if !exportQuiet {
				fmt.Printf("Collected sample %d/%d\n", i+1, exportSamples)
			}

			if i < exportSamples-1 {
				time.Sleep(time.Duration(exportInterval) * time.Second)
			}
		}
	}

	formats := []string{}
	if exportAll {
		formats = []string{"csv", "json"}
	} else {
		formats = []string{exportFormat}
	}

	for _, format := range formats {
		outputFile := exportOutput
		if outputFile == "" {
			timestamp := time.Now().Format("2006-01-02_15-04-05")
			outputFile = fmt.Sprintf("syspulse_export_%s.%s", timestamp, format)
		}

		if !strings.HasSuffix(outputFile, "."+format) {
			if idx := strings.LastIndex(outputFile, "."); idx != -1 {
				outputFile = outputFile[:idx]
			}
			outputFile += "." + format
		}

		outputPath := outputFile
		if !strings.Contains(outputFile, "/") && !strings.Contains(outputFile, "\\") {
			outputPath = fmt.Sprintf("%s/%s", exportDirectory, outputFile)
		}

		var exportFormatEnum export.ExportFormat
		switch format {
		case "csv":
			exportFormatEnum = export.CSV
		case "json":
			exportFormatEnum = export.JSON
		}

		if !exportQuiet {
			fmt.Printf("Exporting %d data points to %s...\n", len(dataPoints), outputPath)
		}

		if err := export.ExportData(dataPoints, outputPath, exportFormatEnum); err != nil {
			return fmt.Errorf("failed to export %s data: %v", format, err)
		}

		if !exportQuiet {
			fmt.Printf("✓ Successfully exported %d data points to %s\n", len(dataPoints), outputPath)
		}
	}

	if !exportQuiet && len(dataPoints) > 0 {
		fmt.Printf("\nData Summary:\n")
		fmt.Printf("  First sample: %s\n", dataPoints[0].Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Last sample:  %s\n", dataPoints[len(dataPoints)-1].Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Total samples: %d\n", len(dataPoints))

		if len(dataPoints) > 1 {
			duration := dataPoints[len(dataPoints)-1].Timestamp.Sub(dataPoints[0].Timestamp)
			fmt.Printf("  Collection duration: %s\n", duration.String())
		}
	}

	return nil
}
