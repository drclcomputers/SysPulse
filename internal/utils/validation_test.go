package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateWidgetConfig(t *testing.T) {
	tests := []struct {
		name        string
		widgetName  string
		config      WidgetConfig
		maxRows     int
		maxCols     int
		shouldError bool
		errorMsg    string
	}{
		{
			name:       "valid widget config",
			widgetName: "CPU",
			config: WidgetConfig{
				Enabled:         true,
				Row:             0,
				Column:          0,
				RowSpan:         1,
				ColSpan:         1,
				MinWidth:        30,
				Weight:          1.0,
				UpdateInterval:  1,
				BorderColor:     "blue",
				ForegroundColor: "white",
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: false,
		},
		{
			name:       "disabled widget should pass",
			widgetName: "Memory",
			config: WidgetConfig{
				Enabled: false,
				Row:     -1,
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: false,
		},
		{
			name:       "negative row",
			widgetName: "Disk",
			config: WidgetConfig{
				Enabled:        true,
				Row:            -1,
				Column:         0,
				RowSpan:        1,
				ColSpan:        1,
				MinWidth:       30,
				Weight:         1.0,
				UpdateInterval: 1,
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: true,
			errorMsg:    "row cannot be negative",
		},
		{
			name:       "negative column",
			widgetName: "Network",
			config: WidgetConfig{
				Enabled:        true,
				Row:            0,
				Column:         -1,
				RowSpan:        1,
				ColSpan:        1,
				MinWidth:       30,
				Weight:         1.0,
				UpdateInterval: 1,
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: true,
			errorMsg:    "column cannot be negative",
		},
		{
			name:       "zero row span",
			widgetName: "Process",
			config: WidgetConfig{
				Enabled:        true,
				Row:            0,
				Column:         0,
				RowSpan:        0,
				ColSpan:        1,
				MinWidth:       30,
				Weight:         1.0,
				UpdateInterval: 1,
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: true,
			errorMsg:    "row span must be greater than 0",
		},
		{
			name:       "zero column span",
			widgetName: "GPU",
			config: WidgetConfig{
				Enabled:        true,
				Row:            0,
				Column:         0,
				RowSpan:        1,
				ColSpan:        0,
				MinWidth:       30,
				Weight:         1.0,
				UpdateInterval: 1,
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: true,
			errorMsg:    "column span must be greater than 0",
		},
		{
			name:       "exceeds row bounds",
			widgetName: "Load",
			config: WidgetConfig{
				Enabled:        true,
				Row:            3,
				Column:         0,
				RowSpan:        2,
				ColSpan:        1,
				MinWidth:       30,
				Weight:         1.0,
				UpdateInterval: 1,
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: true,
			errorMsg:    "exceeds layout row bounds",
		},
		{
			name:       "exceeds column bounds",
			widgetName: "Temperature",
			config: WidgetConfig{
				Enabled:        true,
				Row:            0,
				Column:         1,
				RowSpan:        1,
				ColSpan:        2,
				MinWidth:       30,
				Weight:         1.0,
				UpdateInterval: 1,
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: true,
			errorMsg:    "exceeds layout column bounds",
		},
		{
			name:       "zero minimum width",
			widgetName: "Battery",
			config: WidgetConfig{
				Enabled:        true,
				Row:            0,
				Column:         0,
				RowSpan:        1,
				ColSpan:        1,
				MinWidth:       0,
				Weight:         1.0,
				UpdateInterval: 1,
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: true,
			errorMsg:    "minimum width must be greater than 0",
		},
		{
			name:       "zero weight",
			widgetName: "DiskIO",
			config: WidgetConfig{
				Enabled:        true,
				Row:            0,
				Column:         0,
				RowSpan:        1,
				ColSpan:        1,
				MinWidth:       30,
				Weight:         0,
				UpdateInterval: 1,
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: true,
			errorMsg:    "weight must be greater than 0",
		},
		{
			name:       "zero update interval",
			widgetName: "ProcessTree",
			config: WidgetConfig{
				Enabled:        true,
				Row:            0,
				Column:         0,
				RowSpan:        1,
				ColSpan:        1,
				MinWidth:       30,
				Weight:         1.0,
				UpdateInterval: 0,
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: true,
			errorMsg:    "update interval must be at least 1 second",
		},
		{
			name:       "excessive update interval",
			widgetName: "NetworkConns",
			config: WidgetConfig{
				Enabled:        true,
				Row:            0,
				Column:         0,
				RowSpan:        1,
				ColSpan:        1,
				MinWidth:       30,
				Weight:         1.0,
				UpdateInterval: 2000,
			},
			maxRows:     4,
			maxCols:     2,
			shouldError: true,
			errorMsg:    "update interval cannot exceed 1800 seconds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateWidgetConfig(tt.widgetName, tt.config, tt.maxRows, tt.maxCols)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for test case '%s', but got nil", tt.name)
				} else if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for test case '%s', but got: %v", tt.name, err)
				}
			}
		})
	}
}

func TestValidatePluginWidget(t *testing.T) {
	tests := []struct {
		name        string
		pluginName  string
		config      map[string]interface{}
		maxRows     int
		maxCols     int
		shouldError bool
		errorMsg    string
	}{
		{
			name:       "valid plugin widget config",
			pluginName: "example",
			config: map[string]interface{}{
				"title":            "Example Plugin",
				"enabled":          true,
				"row":              0,
				"column":           0,
				"rowSpan":          1,
				"colSpan":          1,
				"minWidth":         25,
				"border_color":     "aqua",
				"foreground_color": "white",
				"update_interval":  5,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: false,
		},
		{
			name:       "disabled plugin widget should pass",
			pluginName: "docker",
			config: map[string]interface{}{
				"title":   "Docker",
				"enabled": false,
				"row":     -1,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: false,
		},
		{
			name:       "missing title",
			pluginName: "example",
			config: map[string]interface{}{
				"enabled":         true,
				"row":             0,
				"column":          0,
				"rowSpan":         1,
				"colSpan":         1,
				"minWidth":        25,
				"update_interval": 5,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "must have a title",
		},
		{
			name:       "title too long",
			pluginName: "example",
			config: map[string]interface{}{
				"title":           "This is an extremely long title that exceeds the maximum allowed length of 50 characters",
				"enabled":         true,
				"row":             0,
				"column":          0,
				"rowSpan":         1,
				"colSpan":         1,
				"minWidth":        25,
				"update_interval": 5,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "title cannot exceed 50 characters",
		},
		{
			name:       "builtin widget name",
			pluginName: "example",
			config: map[string]interface{}{
				"title":           "CPU",
				"enabled":         true,
				"row":             0,
				"column":          0,
				"rowSpan":         1,
				"colSpan":         1,
				"minWidth":        25,
				"update_interval": 5,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "cannot use built-in widget name",
		},
		{
			name:       "excessive row span",
			pluginName: "docker",
			config: map[string]interface{}{
				"title":           "Docker",
				"enabled":         true,
				"row":             0,
				"column":          0,
				"rowSpan":         3,
				"colSpan":         1,
				"minWidth":        25,
				"update_interval": 5,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "row span cannot exceed half the layout rows",
		},
		{
			name:       "excessive column span",
			pluginName: "docker",
			config: map[string]interface{}{
				"title":           "Docker",
				"enabled":         true,
				"row":             0,
				"column":          0,
				"rowSpan":         1,
				"colSpan":         3,
				"minWidth":        25,
				"update_interval": 5,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "column span cannot exceed half the layout columns",
		},
		{
			name:       "minimum width too small",
			pluginName: "example",
			config: map[string]interface{}{
				"title":           "Example",
				"enabled":         true,
				"row":             0,
				"column":          0,
				"rowSpan":         1,
				"colSpan":         1,
				"minWidth":        5,
				"update_interval": 5,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "minimum width must be at least 10 characters",
		},
		{
			name:       "minimum width too large",
			pluginName: "example",
			config: map[string]interface{}{
				"title":           "Example",
				"enabled":         true,
				"row":             0,
				"column":          0,
				"rowSpan":         1,
				"colSpan":         1,
				"minWidth":        250,
				"update_interval": 5,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "minimum width cannot exceed 200 characters",
		},
		{
			name:       "update interval too small",
			pluginName: "docker",
			config: map[string]interface{}{
				"title":           "Docker",
				"enabled":         true,
				"row":             0,
				"column":          0,
				"rowSpan":         1,
				"colSpan":         1,
				"minWidth":        25,
				"update_interval": 0,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "update interval must be at least 1 second",
		},
		{
			name:       "update interval too large",
			pluginName: "docker",
			config: map[string]interface{}{
				"title":           "Docker",
				"enabled":         true,
				"row":             0,
				"column":          0,
				"rowSpan":         1,
				"colSpan":         1,
				"minWidth":        25,
				"update_interval": 350,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "update interval cannot exceed 300 seconds",
		},
		{
			name:       "invalid border color",
			pluginName: "example",
			config: map[string]interface{}{
				"title":           "Example",
				"enabled":         true,
				"row":             0,
				"column":          0,
				"rowSpan":         1,
				"colSpan":         1,
				"minWidth":        25,
				"border_color":    "invalidcolor",
				"update_interval": 5,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "invalid border color",
		},
		{
			name:       "invalid foreground color",
			pluginName: "example",
			config: map[string]interface{}{
				"title":            "Example",
				"enabled":          true,
				"row":              0,
				"column":           0,
				"rowSpan":          1,
				"colSpan":          1,
				"minWidth":         25,
				"foreground_color": "invalidcolor",
				"update_interval":  5,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "invalid foreground color",
		},
		{
			name:       "float values should be converted",
			pluginName: "example",
			config: map[string]interface{}{
				"title":           "Example",
				"enabled":         true,
				"row":             0.0,
				"column":          0.0,
				"rowSpan":         1.0,
				"colSpan":         1.0,
				"minWidth":        25.0,
				"update_interval": 5.0,
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePluginWidget(tt.pluginName, tt.config, tt.maxRows, tt.maxCols)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for test case '%s', but got nil", tt.name)
				} else if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for test case '%s', but got: %v", tt.name, err)
				}
			}
		})
	}
}

func TestValidatePluginConfig(t *testing.T) {
	tests := []struct {
		name        string
		pluginName  string
		config      map[string]interface{}
		maxRows     int
		maxCols     int
		shouldError bool
		errorMsg    string
	}{
		{
			name:       "valid plugin config",
			pluginName: "example",
			config: map[string]interface{}{
				"name":    "Example Plugin",
				"enabled": true,
				"settings": map[string]interface{}{
					"show_time":  true,
					"show_stats": true,
				},
				"layout": map[string]interface{}{
					"title":           "Example",
					"enabled":         true,
					"row":             0,
					"column":          0,
					"rowSpan":         1,
					"colSpan":         1,
					"minWidth":        25,
					"update_interval": 5,
				},
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: false,
		},
		{
			name:       "disabled plugin should pass",
			pluginName: "docker",
			config: map[string]interface{}{
				"name":    "Docker Plugin",
				"enabled": false,
				"settings": map[string]interface{}{
					"invalid": "config",
				},
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: false,
		},
		{
			name:       "missing name",
			pluginName: "example",
			config: map[string]interface{}{
				"enabled": true,
				"settings": map[string]interface{}{
					"show_time": true,
				},
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "name cannot be empty",
		},
		{
			name:       "name too long",
			pluginName: "example",
			config: map[string]interface{}{
				"name":    "This is an extremely long plugin name that exceeds the maximum allowed length of 100 characters and should trigger a validation error",
				"enabled": true,
				"settings": map[string]interface{}{
					"show_time": true,
				},
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "name cannot exceed 100 characters",
		},
		{
			name:       "too many settings",
			pluginName: "example",
			config: map[string]interface{}{
				"name":     "Example Plugin",
				"enabled":  true,
				"settings": generateManySettings(60),
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "cannot have more than 50 settings",
		},
		{
			name:       "invalid widget layout",
			pluginName: "example",
			config: map[string]interface{}{
				"name":    "Example Plugin",
				"enabled": true,
				"settings": map[string]interface{}{
					"show_time": true,
				},
				"layout": map[string]interface{}{
					"title":           "Example",
					"enabled":         true,
					"row":             0,
					"column":          0,
					"rowSpan":         1,
					"colSpan":         1,
					"minWidth":        5,
					"update_interval": 5,
				},
			},
			maxRows:     4,
			maxCols:     4,
			shouldError: true,
			errorMsg:    "minimum width must be at least 10 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePluginConfig(tt.pluginName, tt.config, tt.maxRows, tt.maxCols)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for test case '%s', but got nil", tt.name)
				} else if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for test case '%s', but got: %v", tt.name, err)
				}
			}
		})
	}
}

func TestValidateTheme(t *testing.T) {
	tests := []struct {
		name        string
		theme       Theme
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid theme",
			theme: Theme{
				UpdateTime: 1,
				Layout: LayoutConfig{
					Rows:    4,
					Columns: 2,
					Spacing: 0,
					CPU: WidgetConfig{
						Enabled:        true,
						Row:            0,
						Column:         0,
						RowSpan:        1,
						ColSpan:        1,
						MinWidth:       30,
						Weight:         1.0,
						UpdateInterval: 1,
					},
					Memory: WidgetConfig{
						Enabled:        true,
						Row:            0,
						Column:         1,
						RowSpan:        1,
						ColSpan:        1,
						MinWidth:       30,
						Weight:         1.0,
						UpdateInterval: 2,
					},
				},
				Export: ExportConfig{
					Enabled:        true,
					Interval:       300,
					Formats:        []string{"json", "csv"},
					Directory:      "exports",
					FilenamePrefix: "syspulse",
				},
			},
			shouldError: false,
		},
		{
			name: "invalid update time",
			theme: Theme{
				UpdateTime: 0,
				Layout: LayoutConfig{
					Rows:    4,
					Columns: 2,
					Spacing: 0,
				},
			},
			shouldError: true,
			errorMsg:    "Update interval must be greater than 0",
		},
		{
			name: "invalid layout rows",
			theme: Theme{
				UpdateTime: 1,
				Layout: LayoutConfig{
					Rows:    0,
					Columns: 2,
					Spacing: 0,
				},
			},
			shouldError: true,
			errorMsg:    "Layout rows must be greater than 0",
		},
		{
			name: "invalid layout columns",
			theme: Theme{
				UpdateTime: 1,
				Layout: LayoutConfig{
					Rows:    4,
					Columns: 0,
					Spacing: 0,
				},
			},
			shouldError: true,
			errorMsg:    "Layout columns must be greater than 0",
		},
		{
			name: "negative spacing",
			theme: Theme{
				UpdateTime: 1,
				Layout: LayoutConfig{
					Rows:    4,
					Columns: 2,
					Spacing: -1,
				},
			},
			shouldError: true,
			errorMsg:    "Layout spacing cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.theme)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for test case '%s', but got nil", tt.name)
				} else if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for test case '%s', but got: %v", tt.name, err)
				}
			}
		})
	}
}

func TestValidateExportConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      ExportConfig
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid export config",
			config: ExportConfig{
				Enabled:        true,
				Interval:       300,
				Formats:        []string{"json", "csv"},
				Directory:      "exports",
				FilenamePrefix: "syspulse",
			},
			shouldError: false,
		},
		{
			name: "disabled export should pass",
			config: ExportConfig{
				Enabled:  false,
				Interval: 0,
			},
			shouldError: false,
		},
		{
			name: "zero interval",
			config: ExportConfig{
				Enabled:        true,
				Interval:       0,
				Formats:        []string{"json"},
				Directory:      "exports",
				FilenamePrefix: "syspulse",
			},
			shouldError: true,
			errorMsg:    "Export interval must be greater than 0",
		},
		{
			name: "no formats",
			config: ExportConfig{
				Enabled:        true,
				Interval:       300,
				Formats:        []string{},
				Directory:      "exports",
				FilenamePrefix: "syspulse",
			},
			shouldError: true,
			errorMsg:    "At least one export format must be specified",
		},
		{
			name: "invalid format",
			config: ExportConfig{
				Enabled:        true,
				Interval:       300,
				Formats:        []string{"invalid"},
				Directory:      "exports",
				FilenamePrefix: "syspulse",
			},
			shouldError: true,
			errorMsg:    "Unsupported export format",
		},
		{
			name: "empty directory",
			config: ExportConfig{
				Enabled:        true,
				Interval:       300,
				Formats:        []string{"json"},
				Directory:      "",
				FilenamePrefix: "syspulse",
			},
			shouldError: true,
			errorMsg:    "Export directory must be specified",
		},
		{
			name: "empty filename prefix",
			config: ExportConfig{
				Enabled:        true,
				Interval:       300,
				Formats:        []string{"json"},
				Directory:      "exports",
				FilenamePrefix: "",
			},
			shouldError: true,
			errorMsg:    "Export filename prefix must be specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExportConfig(tt.config)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for test case '%s', but got nil", tt.name)
				} else if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for test case '%s', but got: %v", tt.name, err)
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && s[:len(substr)] == substr) ||
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
		strings.Contains(s, substr))
}

func generateManySettings(count int) map[string]interface{} {
	settings := make(map[string]interface{})
	for i := 0; i < count; i++ {
		settings[fmt.Sprintf("setting_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	return settings
}
