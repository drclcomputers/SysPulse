package utils

import (
	"fmt"
	"syspulse/internal/errors"
)

func Validate(t Theme) error {
	if t.UpdateTime <= 0 {
		return errors.NewAppError(errors.ValidationError,
			"Update interval must be greater than 0", nil)
	}

	if t.Layout.Rows < 1 {
		return errors.NewAppError(errors.ValidationError, "Layout rows must be greater than 0", nil)
	}
	if t.Layout.Columns < 1 {
		return errors.NewAppError(errors.ValidationError, "Layout columns must be greater than 0", nil)
	}
	if t.Layout.Spacing < 0 {
		return errors.NewAppError(errors.ValidationError, "Layout spacing cannot be negative", nil)
	}

	widgets := []struct {
		name   string
		config WidgetConfig
	}{
		{"CPU", t.Layout.CPU},
		{"Memory", t.Layout.Memory},
		{"Disk", t.Layout.Disk},
		{"Network", t.Layout.Network},
		{"Process", t.Layout.Process},
		{"GPU", t.Layout.GPU},
		{"Load", t.Layout.Load},
		{"Temperature", t.Layout.Temperature},
		{"NetworkConns", t.Layout.NetworkConns},
		{"DiskIO", t.Layout.DiskIO},
		{"ProcessTree", t.Layout.ProcessTree},
		{"Battery", t.Layout.Battery},
	}

	for _, w := range widgets {
		if err := validateWidgetConfig(w.name, w.config, t.Layout.Rows, t.Layout.Columns); err != nil {
			return err
		}
	}

	if err := validateExportConfig(t.Export); err != nil {
		return err
	}

	return nil
}

func validateWidgetConfig(name string, w WidgetConfig, maxRows, maxCols int) error {
	if !w.Enabled {
		return nil
	}

	if w.Row < 0 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("%s widget row cannot be negative", name), nil)
	}
	if w.Column < 0 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("%s widget column cannot be negative", name), nil)
	}
	if w.RowSpan < 1 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("%s widget row span must be greater than 0", name), nil)
	}
	if w.ColSpan < 1 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("%s widget column span must be greater than 0", name), nil)
	}

	if w.Row+w.RowSpan > maxRows {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("%s widget exceeds layout row bounds", name), nil)
	}
	if w.Column+w.ColSpan > maxCols {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("%s widget exceeds layout column bounds", name), nil)
	}

	if w.MinWidth < 1 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("%s widget minimum width must be greater than 0", name), nil)
	}

	if w.Weight <= 0 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("%s widget weight must be greater than 0", name), nil)
	}

	if w.UpdateInterval < 1 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("%s widget update interval must be at least 1 second", name), nil)
	}
	if w.UpdateInterval > 1800 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("%s widget update interval cannot exceed 1800 seconds (30 minutes)", name), nil)
	}

	return nil
}

func validateExportConfig(e ExportConfig) error {
	if e.Enabled {
		if e.Interval <= 0 {
			return errors.NewAppError(errors.ValidationError,
				"Export interval must be greater than 0", nil)
		}

		if len(e.Formats) == 0 {
			return errors.NewAppError(errors.ValidationError,
				"At least one export format must be specified", nil)
		}

		for _, format := range e.Formats {
			if format != "csv" && format != "json" {
				return errors.NewAppError(errors.ValidationError,
					fmt.Sprintf("Unsupported export format: %s (must be 'csv' or 'json')", format), nil)
			}
		}

		if e.Directory == "" {
			return errors.NewAppError(errors.ValidationError,
				"Export directory must be specified", nil)
		}

		if e.FilenamePrefix == "" {
			return errors.NewAppError(errors.ValidationError,
				"Export filename prefix must be specified", nil)
		}
	}

	return nil
}

func ValidatePluginWidget(name string, config interface{}, maxRows, maxCols int) error {
	type PluginWidgetConfig struct {
		Title           string `json:"title"`
		Row             int    `json:"row"`
		Column          int    `json:"column"`
		RowSpan         int    `json:"rowSpan"`
		ColSpan         int    `json:"colSpan"`
		MinWidth        int    `json:"minWidth"`
		Enabled         bool   `json:"enabled"`
		BorderColor     string `json:"border_color"`
		ForegroundColor string `json:"foreground_color"`
		UpdateInterval  int    `json:"update_interval"`
	}

	var w PluginWidgetConfig
	switch v := config.(type) {
	case PluginWidgetConfig:
		w = v
	case map[string]interface{}:
		if title, ok := v["title"].(string); ok {
			w.Title = title
		}
		if row, ok := v["row"].(int); ok {
			w.Row = row
		} else if rowFloat, ok := v["row"].(float64); ok {
			w.Row = int(rowFloat)
		}
		if column, ok := v["column"].(int); ok {
			w.Column = column
		} else if columnFloat, ok := v["column"].(float64); ok {
			w.Column = int(columnFloat)
		}
		if rowSpan, ok := v["rowSpan"].(int); ok {
			w.RowSpan = rowSpan
		} else if rowSpanFloat, ok := v["rowSpan"].(float64); ok {
			w.RowSpan = int(rowSpanFloat)
		}
		if colSpan, ok := v["colSpan"].(int); ok {
			w.ColSpan = colSpan
		} else if colSpanFloat, ok := v["colSpan"].(float64); ok {
			w.ColSpan = int(colSpanFloat)
		}
		if minWidth, ok := v["minWidth"].(int); ok {
			w.MinWidth = minWidth
		} else if minWidthFloat, ok := v["minWidth"].(float64); ok {
			w.MinWidth = int(minWidthFloat)
		}
		if enabled, ok := v["enabled"].(bool); ok {
			w.Enabled = enabled
		}
		if borderColor, ok := v["border_color"].(string); ok {
			w.BorderColor = borderColor
		}
		if foregroundColor, ok := v["foreground_color"].(string); ok {
			w.ForegroundColor = foregroundColor
		}
		if updateInterval, ok := v["update_interval"].(int); ok {
			w.UpdateInterval = updateInterval
		} else if updateIntervalFloat, ok := v["update_interval"].(float64); ok {
			w.UpdateInterval = int(updateIntervalFloat)
		}
	default:
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s has invalid widget configuration type", name), nil)
	}

	if !w.Enabled {
		return nil
	}

	if w.Title == "" {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget must have a title", name), nil)
	}
	if len(w.Title) > 50 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget title cannot exceed 50 characters", name), nil)
	}

	builtinWidgets := []string{"CPU", "Memory Usage", "Disk Usage", "Network Activity", "Processes", "GPU", "Load Average", "Temperature", "Network Connections", "DiskIO", "ProcessTree", "Battery"}
	for _, builtinWidget := range builtinWidgets {
		if w.Title == builtinWidget {
			return errors.NewAppError(errors.ValidationError,
				fmt.Sprintf("Plugin %s widget title cannot use built-in widget name: %s", name, w.Title), nil)
		}
	}

	if w.Row < 0 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget row cannot be negative", name), nil)
	}
	if w.Column < 0 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget column cannot be negative", name), nil)
	}
	if w.RowSpan < 1 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget row span must be greater than 0", name), nil)
	}
	if w.ColSpan < 1 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget column span must be greater than 0", name), nil)
	}

	if w.RowSpan > maxRows/2 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget row span cannot exceed half the layout rows (%d)", name, maxRows/2), nil)
	}
	if w.ColSpan > maxCols/2 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget column span cannot exceed half the layout columns (%d)", name, maxCols/2), nil)
	}

	if w.Row+w.RowSpan > maxRows {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget exceeds layout row bounds", name), nil)
	}
	if w.Column+w.ColSpan > maxCols {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget exceeds layout column bounds", name), nil)
	}

	if w.MinWidth < 10 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget minimum width must be at least 10 characters", name), nil)
	}
	if w.MinWidth > 200 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget minimum width cannot exceed 200 characters", name), nil)
	}

	if w.UpdateInterval < 1 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget update interval must be at least 1 second", name), nil)
	}
	if w.UpdateInterval > 300 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s widget update interval cannot exceed 300 seconds (5 minutes)", name), nil)
	}

	validColors := []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
		"gray", "grey", "darkred", "darkgreen", "darkyellow", "darkblue", "darkmagenta", "darkcyan",
		"silver", "maroon", "lime", "olive", "navy", "purple", "teal", "aqua", "orange", "pink"}

	if w.BorderColor != "" {
		isValid := false
		for _, validColor := range validColors {
			if w.BorderColor == validColor {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.NewAppError(errors.ValidationError,
				fmt.Sprintf("Plugin %s widget has invalid border color: %s", name, w.BorderColor), nil)
		}
	}

	if w.ForegroundColor != "" {
		isValid := false
		for _, validColor := range validColors {
			if w.ForegroundColor == validColor {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.NewAppError(errors.ValidationError,
				fmt.Sprintf("Plugin %s widget has invalid foreground color: %s", name, w.ForegroundColor), nil)
		}
	}

	return nil
}

func ValidatePluginConfig(name string, config interface{}, maxRows, maxCols int) error {
	type PluginConfig struct {
		Name     string                 `json:"name"`
		Enabled  bool                   `json:"enabled"`
		Settings map[string]interface{} `json:"settings"`
		Layout   interface{}            `json:"layout"`
	}

	var p PluginConfig
	switch v := config.(type) {
	case PluginConfig:
		p = v
	case map[string]interface{}:
		if nameVal, ok := v["name"].(string); ok {
			p.Name = nameVal
		}
		if enabled, ok := v["enabled"].(bool); ok {
			p.Enabled = enabled
		}
		if settings, ok := v["settings"].(map[string]interface{}); ok {
			p.Settings = settings
		}
		if layout, ok := v["layout"]; ok {
			p.Layout = layout
		}
	default:
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s has invalid configuration type", name), nil)
	}

	if !p.Enabled {
		return nil
	}

	if p.Name == "" {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s name cannot be empty", name), nil)
	}
	if len(p.Name) > 100 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s name cannot exceed 100 characters", name), nil)
	}

	if len(p.Settings) > 50 {
		return errors.NewAppError(errors.ValidationError,
			fmt.Sprintf("Plugin %s cannot have more than 50 settings", name), nil)
	}

	if p.Layout != nil {
		if err := ValidatePluginWidget(name, p.Layout, maxRows, maxCols); err != nil {
			return err
		}
	}

	return nil
}
