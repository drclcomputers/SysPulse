# SysPulse Plugin System Usage Guide

## Overview

This guide explains how to use and develop plugins for SysPulse.

## Plugin System Features

### ✅ Implemented Features

1. **Plugin Manager**: Centralized management of all plugins
2. **Widget Integration**: Plugins can create custom widgets that integrate with the main dashboard
3. **Lifecycle Management**: Plugins can be enabled, disabled, and refreshed
4. **Configuration System**: JSON-based configuration for plugin settings and layout
5. **Data Collection**: Plugins can collect and export custom data
6. **Real-time Updates**: Plugin widgets update automatically
7. **Grid Layout Integration**: Plugin widgets are positioned using the grid layout system

### Current Plugins

#### 1. Example Plugin
- **Purpose**: Demonstrates basic plugin functionality
- **Features**: 
  - Shows current time
  - Displays update counter
  - Configurable message
- **Status**: Enabled by default
- **Location**: `row: 0, column: 0`

#### 2. Docker Plugin
- **Purpose**: Monitors Docker containers and images
- **Features**:
  - Container status monitoring
  - Image information
  - Resource usage statistics
- **Status**: Disabled by default (requires Docker)
- **Location**: `row: 1, column: 5`

## How to Use

### 1. Start SysPulse
```bash
./syspulse.exe
```

### 2. View Plugin Widgets
- Plugin widgets appear in the main dashboard grid
- Each plugin has its own bordered widget with a title
- Widgets update automatically every 2 seconds

### 3. Configure Plugins
Edit `plugins_config.json` to:
- Enable/disable plugins
- Change widget positioning
- Modify plugin settings

Example configuration:
```json
{
  "plugins": {
    "example": {
      "name": "Example Plugin",
      "enabled": true,
      "settings": {
        "show_time": true,
        "show_stats": true,
        "custom_message": "Hello from SysPulse!"
      },
      "layout": {
        "title": "Example",
        "row": 0,
        "column": 5,
        "rowSpan": 1,
        "colSpan": 2,
        "minWidth": 25,
        "enabled": true
      }
    }
  }
}
```

## Plugin Development

### Creating a New Plugin

1. **Implement the Plugin Interface**:
```go
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
```

2. **Example Plugin Implementation**:
```go
package plugins

import (
    "fmt"
    "time"
    "github.com/rivo/tview"
)

type MyPlugin struct {
    widget *tview.TextView
}

func NewMyPlugin() *MyPlugin {
    return &MyPlugin{}
}

func (p *MyPlugin) Name() string { return "my_plugin" }
func (p *MyPlugin) Version() string { return "1.0.0" }
func (p *MyPlugin) Description() string { return "My custom plugin" }
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
        tv.SetText(fmt.Sprintf("Current time: %s", time.Now().Format("15:04:05")))
    }
    return nil
}

func (p *MyPlugin) CollectData() (map[string]interface{}, error) {
    return map[string]interface{}{
        "timestamp": time.Now().Unix(),
        "status": "active",
    }, nil
}

func (p *MyPlugin) ExportData() interface{} {
    data, _ := p.CollectData()
    return data
}

func (p *MyPlugin) Cleanup() error {
    return nil
}
```

3. **Register the Plugin**:
Add your plugin to the `InitializePluginSystem` function in `internal/plugins/integration.go`:
```go
myPlugin := NewMyPlugin()
if err := pluginManager.LoadPlugin(myPlugin); err != nil {
    log.Printf("Failed to load my plugin: %v", err)
}
```

### Plugin Configuration

Add plugin configuration to `plugins_config.json`:
```json
{
  "plugins": {
    "my_plugin": {
      "name": "My Plugin",
      "enabled": true,
      "settings": {
        "custom_setting": "value"
      },
      "layout": {
        "title": "My Plugin",
        "row": 2,
        "column": 0,
        "rowSpan": 1,
        "colSpan": 1,
        "minWidth": 20,
        "enabled": true
      }
    }
  }
}
```

## Architecture

### Key Components

1. **Plugin Interface** (`internal/plugins/interface.go`):
   - Defines the contract for all plugins
   - Includes lifecycle methods and data collection

2. **Plugin Manager** (`internal/plugins/manager.go`):
   - Manages plugin loading, unloading, and updates
   - Thread-safe operations
   - Widget creation and management

3. **Integration Layer** (`internal/plugins/integration.go`):
   - Connects plugins to the dashboard
   - Handles initialization and layout integration
   - Provides utility functions for plugin management

4. **Example Plugins**:
   - `internal/plugins/example.go`: Basic example plugin
   - `internal/plugins/docker.go`: Docker monitoring plugin

### Integration Points

- **Dashboard Initialization**: Plugins are initialized when the dashboard starts
- **Layout Integration**: Plugin widgets are added to the main grid layout
- **Update Workers**: Plugin widgets are updated every 2 seconds
- **Data Export**: Plugin data is included in system exports

## Troubleshooting

### Common Issues

1. **Plugin Not Loading**:
   - Check if the plugin is enabled in `plugins_config.json`
   - Verify the plugin implementation is correct
   - Check logs for error messages

2. **Widget Not Appearing**:
   - Ensure the layout configuration is correct
   - Check for positioning conflicts with other widgets
   - Verify the widget is created successfully

3. **Build Errors**:
   - Ensure all imports are correct
   - Check for circular dependencies
   - Verify interface implementations

## Testing

To test the plugin system:

1. **Start SysPulse**:
```bash
./syspulse.exe
```

2. **Verify Example Plugin**:
   - Look for the "Example Plugin" widget in the dashboard
   - It should show the current time and update counter
   - The widget should update every 2 seconds

## Future Enhancements

Potential areas for future development:

1. **Dynamic Plugin Loading**: Load plugins from separate files/directories
2. **Plugin Marketplace**: Download and install plugins from a repository
3. **Plugin Configuration UI**: Manage plugins through the dashboard UI
4. **Plugin Themes**: Custom styling for plugin widgets
5. **Plugin Events**: Inter-plugin communication system
6. **Plugin Permissions**: Security and access control for plugins
