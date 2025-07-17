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
					fmt.Sprintf("Invalid export format: %s (must be 'csv' or 'json')", format), nil)
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
