package export

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createTestData() []DataPoint {
	return []DataPoint{
		{
			Timestamp: time.Now(),
			CPU:       []float64{10.5, 20.3, 15.7},
			Memory: struct {
				Total     uint64
				Used      uint64
				SwapTotal uint64
				SwapUsed  uint64
			}{
				Total:     16000000000,
				Used:      8000000000,
				SwapTotal: 8000000000,
				SwapUsed:  1000000000,
			},
			Disk: struct {
				Path     string
				Total    uint64
				Used     uint64
				UsedPerc float64
				IOReads  uint64
				IOWrites uint64
			}{
				Path:     "/",
				Total:    500000000000,
				Used:     250000000000,
				UsedPerc: 50.0,
				IOReads:  1000,
				IOWrites: 500,
			},
			Network: struct {
				BytesSent     uint64
				BytesReceived uint64
				PacketsSent   uint64
				PacketsRecv   uint64
			}{
				BytesSent:     1000000,
				BytesReceived: 2000000,
				PacketsSent:   1000,
				PacketsRecv:   2000,
			},
		},
	}
}

func TestExportData(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "syspulse_test_export")
	defer os.RemoveAll(tmpDir)

	testData := createTestData()

	t.Run("CSV Export", func(t *testing.T) {
		csvPath := filepath.Join(tmpDir, "test.csv")
		err := ExportData(testData, csvPath, CSV)
		if err != nil {
			t.Fatalf("Failed to export CSV: %v", err)
		}

		file, err := os.Open(csvPath)
		if err != nil {
			t.Fatalf("Failed to open CSV file: %v", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read CSV: %v", err)
		}

		if len(records) != 2 {
			t.Errorf("Expected 2 CSV records, got %d", len(records))
		}
	})

	t.Run("JSON Export", func(t *testing.T) {
		jsonPath := filepath.Join(tmpDir, "test.json")
		err := ExportData(testData, jsonPath, JSON)
		if err != nil {
			t.Fatalf("Failed to export JSON: %v", err)
		}

		file, err := os.Open(jsonPath)
		if err != nil {
			t.Fatalf("Failed to open JSON file: %v", err)
		}
		defer file.Close()

		var decoded []DataPoint
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&decoded)
		if err != nil {
			t.Fatalf("Failed to decode JSON: %v", err)
		}

		if len(decoded) != len(testData) {
			t.Errorf("Expected %d data points, got %d", len(testData), len(decoded))
		}
	})

	t.Run("Invalid Format", func(t *testing.T) {
		err := ExportData(testData, "test.txt", ExportFormat(99))
		if err == nil {
			t.Error("Expected error for invalid format")
		}
	})
}
