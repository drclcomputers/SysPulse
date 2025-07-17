/*
Package syspulse provides a comprehensive terminal-based system monitoring tool with
real-time metrics visualization, process management capabilities, advanced logging,
cross-platform GPU monitoring, and an extensible plugin system.

# Overview

SysPulse is designed to offer detailed system monitoring through an intuitive
terminal user interface. It provides real-time information about:

  - CPU usage and statistics (per-core and total)
  - Memory utilization (RAM and swap)
  - Disk usage and I/O metrics
  - Network activity and interface details
  - GPU monitoring (cross-platform support for NVIDIA, AMD, Intel)
  - Process management with filtering and search
  - Performance metrics and self-monitoring
  - Custom plugin widgets for extensible monitoring

# Core Features

  - Real-time Monitoring:
    CPU, memory, disk, network, and GPU metrics updated in real-time

  - Plugin System:
    Extensible architecture allowing custom widgets and monitoring capabilities

  - Cross-platform GPU Support:
    Windows (WMI), Linux (nvidia-smi, sysfs), macOS (system_profiler)

  - Process Management:
    List, search, filter, kill, and manage system processes

  - Interactive UI:
    Keyboard-driven interface with fully customizable layouts and intelligent widget cycling

  - Performance Metrics:
    Built-in performance monitoring of all update operations with self-optimization

  - Data Export:
    Automatic periodic export of monitoring data in CSV and JSON formats

  - Advanced Logging:
    Comprehensive logging system with rotation, multiple severity levels, and dual output

  - Enterprise Features:
    Configuration validation, error recovery, and graceful degradation

  - Themeable:
    Fully configurable colors, layouts, and widget positioning

  - Resource Efficient:
    Minimal system overhead while monitoring with adaptive performance tuning

# Architecture

The application is structured into several key packages:

  - services/: Core monitoring functionality

  - cpu/: CPU monitoring and statistics

  - memory/: Memory usage tracking

  - disk/: Disk usage and I/O monitoring

  - network/: Network interface monitoring

  - gpu/: Cross-platform GPU monitoring

  - processes/: Process management and monitoring

  - sysinfo/: System information gathering

  - ui/: Terminal user interface components with intelligent widget cycling

  - plugins/: Plugin system for extensible monitoring

  - internal/: Internal utilities and shared components

  - logger/v2/: Advanced logging with rotation

  - metrics/: Performance monitoring and optimization

  - export/: Data export functionality

  - errors/: Error handling and types

  - utils/: Shared utilities and helper functions

# Plugin System

The plugin system allows extending SysPulse with custom monitoring capabilities:

  - Plugin Interface: Defines the contract for all plugins with lifecycle methods

  - Plugin Manager: Manages loading, unloading, and updating of plugins

  - Widget Integration: Plugins create custom tview widgets for the dashboard

  - Smart Focus Integration: Plugin widgets participate in the intelligent focus cycling system

  - Configuration System: JSON-based configuration for plugin settings and layout

  - Data Collection: Plugins can collect custom metrics and export data

  - Real-time Updates: Plugin widgets update automatically with configurable intervals

  - Built-in Plugins: Example and Docker monitoring plugins included

Creating a plugin involves implementing the Plugin interface:

	type Plugin interface {
	    Name() string
	    Version() string
	    Description() string
	    Author() string
	    Initialize() error
	    CreateWidget() (tview.Primitive, error)
	    UpdateWidget(widget tview.Primitive) error
	    CollectData() (map[string]interface{}, error)
	    ExportData() interface{}
	    Cleanup() error
	}

# Usage

Basic usage:

	import "github.com/drclcomputers/syspulse"

	func main() {
	    dashboard := syspulse.NewDashboard()
	    dashboard.Run()
	}

Plugin usage:

	import "github.com/drclcomputers/syspulse/internal/plugins"

	func main() {
	    dashboard := syspulse.NewDashboard()

	    // Initialize plugin system
	    plugins.InitializePluginSystem(dashboard)

	    dashboard.Run()
	}

Creating a custom plugin:

	type MyPlugin struct {
	    widget *tview.TextView
	}

	func (p *MyPlugin) Name() string { return "My Plugin" }
	func (p *MyPlugin) Version() string { return "1.0.0" }
	func (p *MyPlugin) Description() string { return "Custom monitoring plugin" }
	func (p *MyPlugin) Author() string { return "Your Name" }

	func (p *MyPlugin) Initialize() error {
	    return nil
	}

	func (p *MyPlugin) CreateWidget() (tview.Primitive, error) {
	    p.widget = tview.NewTextView()
	    p.widget.SetBorder(true)
	    p.widget.SetTitle("My Plugin")
	    return p.widget, nil
	}

	func (p *MyPlugin) UpdateWidget(widget tview.Primitive) error {
	    if tv, ok := widget.(*tview.TextView); ok {
	        tv.SetText("Hello from my plugin!")
	    }
	    return nil
	}

	func (p *MyPlugin) CollectData() (map[string]interface{}, error) {
	    return map[string]interface{}{
	        "status": "active",
	        "timestamp": time.Now().Unix(),
	    }, nil
	}

	func (p *MyPlugin) ExportData() interface{} {
	    data, _ := p.CollectData()
	    return data
	}

	func (p *MyPlugin) Cleanup() error {
	    return nil
	}

# Keyboard Navigation

The application features intelligent keyboard navigation:

  - TAB/Shift+TAB: Cycle through widgets in screen order (top-to-bottom, left-to-right)
  - Widget focus follows the visual layout positioning
  - Plugin widgets are seamlessly integrated into the cycling order
  - Focus cycling respects enabled/disabled widget states
  - Adaptive navigation based on current theme and layout configuration

The cycling system automatically sorts widgets by their grid position:
  - Primary sort: Row position (top to bottom)
  - Secondary sort: Column position (left to right)
  - Header widget always appears first in the cycle
  - Disabled widgets are excluded from the cycling order

# Error Handling

The application implements a robust error handling system:

  - Critical errors that prevent core functionality are logged and may trigger
    application shutdown
  - Non-critical errors are logged and reported to the user via the UI
  - All errors include context and stack traces
  - Custom error types for different categories of errors
  - Graceful degradation when components fail

Example error handling:

	if err := dashboard.LoadTheme(); err != nil {
	    if errors.IsType(err, errors.ConfigError) {
	        // Handle configuration error
	    }
	    return err
	}

# Configuration

The application can be configured via a config.json file:

	{
	    "background": "black",
	    "foreground": "white",
	    "cpu": {
	        "bar_low": "green",
	        "bar_high": "red"
	    },
	    "memory": {
	        "vmem_gauge": "blue",
	        "smem_gauge": "cyan"
	    },
	    "gpu": {
	        "bar_low": "green",
	        "bar_high": "red"
	    },
	    "layout": {
	        "rows": 4,
	        "columns": 2,
	        "spacing": 0,
	        "cpu": {
	            "enabled": true,
	            "row": 0,
	            "column": 0,
	            "rowSpan": 2,
	            "colSpan": 1
	        },
	        "gpu": {
	            "enabled": false,
	            "row": 3,
	            "column": 0,
	            "rowSpan": 1,
	            "colSpan": 1
	        }
	    },
	    "processsort": "cpu",
	    "updatetime": 1
	}

# Layout System

The application uses a flexible grid-based layout system:

  - Grid positioning: Widgets are positioned using row/column coordinates
  - Intelligent focus cycling: TAB navigation follows visual layout order
  - Plugin integration: Plugin widgets participate in the focus cycling system
  - Responsive design: Layout adapts to terminal size while maintaining proportions
  - Configuration-driven: All positioning and sizing defined in JSON configuration

Layout positioning example:

	// Widget at row 0, column 0 will be focused first
	// Widget at row 0, column 1 will be focused second
	// Widget at row 1, column 0 will be focused third
	// And so on...

The system automatically sorts widgets for focus cycling:
 1. Header widget (always first)
 2. Core widgets sorted by row, then column
 3. Plugin widgets integrated based on their configured positions

# GPU Monitoring

Cross-platform GPU monitoring with support for:

  - Windows: WMI-based GPU detection and monitoring
  - Linux: nvidia-smi, AMD sysfs, Intel graphics support
  - macOS: system_profiler and Metal API integration

Example GPU usage:

	gpus, err := gpu.GetGPUInfo()
	if err != nil {
	    log.Printf("GPU monitoring unavailable: %v", err)
	    return
	}

	for _, gpu := range gpus {
	    fmt.Printf("GPU: %s (%s)\n", gpu.Name, gpu.Vendor)
	    fmt.Printf("Memory: %s / %s\n",
	        formatBytes(gpu.MemoryUsed), formatBytes(gpu.MemoryTotal))
	    fmt.Printf("Temperature: %.1f°C\n", gpu.Temperature)
	    fmt.Printf("Usage: %.1f%%\n", gpu.Usage)
	}

# Input Validation

The application implements comprehensive input validation:

  - Configuration file validation with detailed error messages
  - Theme validation for color schemes and layout parameters
  - Layout validation for grid positioning and sizing
  - Process management input validation for safety
  - GPU configuration validation for cross-platform compatibility

Example validation:

	func (t *Theme) Validate() error {
	    if t.Layout.Rows < 1 {
	        return errors.NewAppError(errors.ValidationError,
	            "Layout rows must be greater than 0")
	    }
	    if t.UpdateTime <= 0 {
	        return errors.NewAppError(errors.ValidationError,
	            "Update interval must be greater than 0")
	    }
	    return nil
	}

# Keyboard Controls

The application supports various keyboard shortcuts:

  - Global Controls:

  - TAB/Shift+TAB: Cycle through widgets in screen order

  - ESC: Unfocus current widget

  - Q: Quit application

  - H: Show help screen

  - I: Show detailed system information

  - Navigation:

  - C: Focus CPU widget

  - M: Focus Memory widget

  - D: Focus Disk widget

  - N: Focus Network widget

  - P: Focus Process widget

  - G: Focus GPU widget (if enabled)

  - Process Management:

  - K: Kill selected process (with confirmation)

  - F: Search/filter processes

  - S/Y: Toggle process sorting (CPU/Memory)

  - Up/Down or W/S: Navigate process list

  - Enter: View process details

# Performance Monitoring

The application includes built-in performance monitoring:

  - Update duration tracking for all metrics
  - Error counting and tracking
  - Performance statistics export
  - Automatic cleanup of old performance data

Example usage:

	metrics := metrics.New(time.Hour)

	// Record update duration
	start := time.Now()
	UpdateCPU()
	metrics.RecordUpdateDuration(metrics.CPUUpdate, time.Since(start))

	// Get performance stats
	stats := metrics.GetStats()

# Data Export

SysPulse automatically exports monitoring data:

  - Supports CSV and JSON formats
  - Periodic export (every 5 minutes)
  - Final export on shutdown
  - Custom export paths and intervals

Example export:

	dataPoints := []export.DataPoint{...}

	// Export as CSV
	err := export.ExportData(dataPoints, "metrics.csv", export.CSV)

	// Export as JSON
	err := export.ExportData(dataPoints, "metrics.json", export.JSON)

# Advanced Logging

The logging system provides:

  - Multiple severity levels (DEBUG, INFO, WARN, ERROR)
  - Automatic log rotation
  - Contextual information (timestamp, file, line number)
  - Console and file output

Example logging:

	logger, err := logger.New("logs", logger.INFO)
	if err != nil {
	    panic(err)
	}
	defer logger.Close()

	logger.Info("Application started")
	logger.Error("An error occurred")

# Examples

Basic monitoring setup:

	dashboard := syspulse.NewDashboard()
	if err := dashboard.LoadTheme("config.json"); err != nil {
	    log.Fatal(err)
	}
	dashboard.Run()

Custom theme configuration:

	theme := &utils.Theme{
	    Background: "black",
	    Foreground: "white",
	    CPU: utils.CPUModel{
	        BarLow:  "green",
	        BarHigh: "red",
	    },
	}
	dashboard.SetTheme(theme)

Process management:

	processes.UpdateProcesses(dashboard)
	dashboard.ProcessWidget.SetTitle("Processes")

For more examples and detailed documentation, visit the project repository:
https://github.com/drclcomputers/syspulse
*/
package main
