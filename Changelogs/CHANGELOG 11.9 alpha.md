# SysPulse Changelog

## Version 11.9 Alpha - Custom Update Intervals Release
*Release Date: July 18, 2025*

### ✨ New Features

#### Comprehensive Widget Information Modals
- **Complete modal system** - All widgets now feature detailed information modals accessible with the 'I' key
- **Memory widget modal** - Shows RAM/Swap usage, health status, and memory optimization tips
- **Disk widget modal** - Displays per-partition information, usage statistics, and disk health advice
- **Network widget modal** - Provides interface details, transfer rates, and network statistics
- **Enhanced CPU modal** - Now includes current CPU usage, per-core statistics, and average utilization
- **Consistent interface** - All widget modals follow the same design pattern and interaction model

#### Custom Widget Update Intervals
- **Individual widget update intervals** - Each widget now has its own configurable update interval instead of a global frequency category
- **Optimized default intervals** based on widget type and data volatility:
  - **CPU**: 1 second (high frequency - real-time monitoring)
  - **Memory**: 2 seconds (medium frequency - moderate changes)
  - **Disk**: 3 seconds (lower frequency - slower changes)
  - **Network**: 1 second (high frequency - real-time traffic)
  - **Process**: 2 seconds (medium frequency - process changes)
  - **GPU**: 3 seconds (lower frequency - stable metrics)
  - **Load**: 5 seconds (low frequency - system load averages)
  - **Temperature**: 5 seconds (low frequency - thermal changes)
  - **Network Connections**: 5 seconds (low frequency - connection changes)
  - **Disk I/O**: 3 seconds (lower frequency - I/O statistics)
  - **Process Tree**: 5 seconds (low frequency - process hierarchy)
  - **Battery**: 10 seconds (very low frequency - power status)

#### Configuration Updates
- **Theme files updated** - All built-in themes now include custom update intervals
- **JSON schema enhanced** - Added `update_interval` field to `WidgetConfig` structure
- **Backward compatibility** - Maintains compatibility with existing configurations (defaults to 1 second if not specified)

### 🔧 Technical Implementation

#### Modal Information System
- **Formatted info functions** - New `Get[Widget]FormattedInfo()` functions provide comprehensive system details
- **Memory information** - RAM/Swap usage, health status, performance tips, and memory optimization advice
- **Disk information** - Per-partition analysis, filesystem details, usage statistics, and health recommendations
- **Network information** - Interface statistics, transfer rates, packet counts, and network performance metrics
- **CPU information** - Hardware specifications, current usage calculation, per-core breakdown, and thermal status
- **Unified modal interface** - Consistent tview.NewModal() implementation across all widgets

#### Worker System Enhancement
- **Individual widget workers** - Each widget now runs in its own goroutine with custom timing
- **Efficient resource usage** - Reduces unnecessary updates for stable metrics
- **Thread-safe operations** - Proper synchronization between workers and UI updates

#### Performance Optimizations
- **Reduced CPU usage** - Less frequent updates for stable metrics like battery and temperature
- **Improved responsiveness** - Critical metrics (CPU, Network) maintain high update frequency
- **Smart resource allocation** - Different update frequencies based on data volatility

### 📋 Configuration Format

#### Widget Configuration Example
```json
{
  "layout": {
    "cpu": {
      "enabled": true,
      "row": 0,
      "column": 0,
      "rowSpan": 2,
      "colSpan": 1,
      "minWidth": 30,
      "weight": 1.0,
      "border_color": "cyan",
      "foreground_color": "cyan",
      "update_interval": 1
    },
    "memory": {
      "enabled": true,
      "row": 0,
      "column": 1,
      "rowSpan": 1,
      "colSpan": 1,
      "minWidth": 30,
      "weight": 1.0,
      "border_color": "magenta",
      "foreground_color": "magenta",
      "update_interval": 2
    }
  }
}
```

### 🎯 Benefits

#### Performance Improvements
- **Reduced system load** - Lower average CPU usage due to optimized update frequencies
- **Better battery life** - Especially beneficial for laptops with less frequent battery monitoring
- **Smoother operation** - More stable performance under varying system loads

#### User Experience
- **Responsive interface** - Critical metrics update quickly while maintaining smooth operation
- **Customizable timing** - Users can adjust individual widget update intervals in configuration files
- **Logical defaults** - Sensible default intervals based on typical monitoring needs

#### Enhanced Information Access
- **Detailed system insights** - Press 'I' on any widget to access comprehensive system information
- **Context-aware help** - Each widget provides specific tips and recommendations based on current system state
- **Consistent interaction** - Unified modal interface across all widgets for familiar user experience
- **Rich information display** - Color-coded status indicators, health assessments, and performance recommendations

### 📈 Migration from Previous Versions

#### Automatic Migration
- **Seamless upgrade** - Existing configurations automatically work with 1-second default intervals
- **Theme updates** - All built-in themes updated with optimized intervals
- **No breaking changes** - Maintains full backward compatibility

#### Customization Guide
1. **Edit theme files** - Modify `update_interval` values in theme JSON files
2. **Restart application** - Changes take effect on next startup
3. **Monitor performance** - Adjust intervals based on system performance and preferences

### 🔄 Update Frequency Rationale

#### High Frequency (1 second)
- **CPU usage** - Real-time monitoring crucial for performance analysis
- **Network traffic** - Important for bandwidth monitoring and troubleshooting

#### Medium Frequency (2 seconds)
- **Memory usage** - Moderate changes, good balance of accuracy and performance
- **Process list** - Process changes occur frequently enough to warrant regular updates

#### Lower Frequency (3 seconds)
- **Disk usage** - Disk space changes slowly, moderate frequency sufficient
- **GPU metrics** - Graphics performance relatively stable
- **Disk I/O** - I/O statistics don't require constant monitoring

#### Low Frequency (5 seconds)
- **System load** - Load averages change slowly by design
- **Temperature** - Thermal changes are gradual
- **Network connections** - Connection states change less frequently
- **Process tree** - Process hierarchy changes slowly

#### Very Low Frequency (10 seconds)
- **Battery status** - Power levels change very slowly

### 🔧 Developer Notes

#### Worker Implementation
```go
func startWidgetWorker(d *utils.Dashboard, quit chan struct{}, widgetName string, updateFunc func(), config utils.WidgetConfig) {
    if !config.Enabled {
        return
    }

    interval := config.UpdateInterval
    if interval <= 0 {
        interval = 1
    }

    go func() {
        ticker := time.NewTicker(time.Duration(interval) * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                updateFunc()
                d.App.QueueUpdateDraw(func() {})
            case <-quit:
                return
            }
        }
    }()
}
```

#### Configuration Structure
```go
type WidgetConfig struct {
    Enabled         bool    `json:"enabled"`
    Row             int     `json:"row"`
    Column          int     `json:"column"`
    RowSpan         int     `json:"rowSpan"`
    ColSpan         int     `json:"colSpan"`
    MinWidth        int     `json:"minWidth"`
    Weight          float64 `json:"weight"`
    BorderColor     string  `json:"border_color"`
    ForegroundColor string  `json:"foreground_color"`
    UpdateInterval  int     `json:"update_interval"` // Update interval in seconds
}
```

#### Modal Information Functions
```go
// Memory information with health status and tips
func GetMemoryFormattedInfo() string {
    // Returns comprehensive memory usage, swap details, and optimization advice
}

// Disk information with per-partition analysis
func GetDiskFormattedInfo() string {
    // Returns detailed disk usage, filesystem info, and health recommendations
}

// Network information with interface statistics
func GetNetworkFormattedInfo() string {
    // Returns interface details, transfer rates, and network metrics
}

// Enhanced CPU information with current usage
func GetCpuFormattedInfo() string {
    // Returns hardware specs, current usage, per-core stats, and thermal status
}
```

---

## Version 11.7 Alpha - First Public Release
*Release Date: July 15, 2025*

### 🎉 Initial Public Release

This marks the first public release of SysPulse, a powerful terminal-based system monitoring tool written in Go. This alpha version includes all core functionality and is ready for community testing and feedback.

### 🌟 Core Features Added

#### System Monitoring
- **Real-time CPU monitoring** with per-core statistics and usage bars
- **Memory tracking** including virtual memory, swap usage, and detailed statistics
- **Disk monitoring** with usage percentages, I/O statistics, and multiple disk support
- **Network interface monitoring** with bytes sent/received and packet statistics
- **GPU monitoring** with cross-platform support (NVIDIA, AMD, Intel)
- **Process management** with interactive process list, filtering, and termination

#### User Interface
- **Terminal-based UI** built with tview for rich terminal experiences
- **Grid-based layout system** with configurable rows, columns, and widget positioning
- **Smart widget cycling** that follows visual layout order (top-to-bottom, left-to-right)
- **Interactive navigation** with TAB/Shift+TAB cycling and direct widget shortcuts
- **Responsive design** that adapts to terminal size changes
- **Mouse support** for modern terminals

#### Advanced Process Management
- **Process filtering** by CPU/Memory usage with real-time search
- **Process sorting** by various metrics (CPU, Memory)
- **Detailed process information** modal with comprehensive process details
- **Safe process termination** with confirmation dialogs
- **Process tree visualization** (where supported)

#### Configuration & Theming
- **JSON-based configuration** with flexible theme support
- **8 Pre-built themes** including Cyberpunk, Matrix, Ocean, Sunset, Monochrome, Neon, Forest, and Fire
- **Customizable colors** for all UI elements and components
- **Layout customization** with precise widget positioning
- **Configurable update intervals** for performance optimization

#### Data Export & Analytics
- **Automatic data export** every 5 minutes during operation
- **Final export** on application shutdown with comprehensive metrics
- **Multiple export formats** (CSV for analysis, JSON for programmatic use)
- **Historical data tracking** with timestamps and comprehensive system metrics
- **Plugin data integration** in export files

#### Plugin System
- **Extensible plugin architecture** with well-defined interfaces
- **Built-in Example Plugin** demonstrating basic plugin functionality
- **Docker Plugin** for container and image monitoring
- **Plugin configuration** through dedicated JSON configuration files
- **Plugin lifecycle management** with initialization, updates, and cleanup
- **Plugin widget integration** with seamless UI integration

#### Enterprise Features
- **Advanced logging system** with automatic rotation and multiple severity levels
- **Performance self-monitoring** with resource usage tracking
- **Comprehensive error handling** with detailed context and recovery
- **Configuration validation** with fallback to default configurations
- **Cross-platform compatibility** (Windows, Linux, macOS)

### 🔧 Technical Implementation

#### Architecture
- **Modular design** with separate packages for each monitoring component
- **Concurrent processing** with goroutines for efficient resource monitoring
- **Memory-efficient** data structures and update mechanisms
- **Platform-specific optimizations** using build tags

#### Logging & Debugging
- **Multi-level logging** (DEBUG, INFO, WARN, ERROR, FATAL)
- **Daily log rotation** with automatic file management
- **Contextual logging** with file, line, function, and timestamp information
- **Dual output** supporting both console and file logging

#### Performance
- **Optimized update cycles** with configurable refresh rates
- **Resource-aware monitoring** that adapts to system load
- **Memory management** with proper cleanup and garbage collection
- **Efficient data structures** for high-frequency updates

#### Bug Fixes

### 🎮 User Experience

#### Keyboard Controls
- **Global navigation**: TAB/Shift+TAB for widget cycling, ESC to unfocus
- **Direct widget access**: C (CPU), M (Memory), D (Disk), N (Network), P (Process), G (GPU)
- **Process management**: K (kill), F (filter), Y (toggle sort), I (info)
- **Application controls**: Q (quit), H (help), I (system info)

#### Visual Features
- **Color-coded usage bars** with thresholds for different usage levels
- **Border highlighting** for focused widgets
- **Real-time updates** with smooth transitions
- **Comprehensive system information** modals

### 📁 File Structure
```
syspulse/
├── cmd/                    # Command-line interface
├── internal/              # Internal packages
│   ├── errors/           # Error handling
│   ├── export/           # Data export functionality
│   ├── logger/v2/        # Advanced logging system
│   ├── metrics/          # Performance monitoring
│   ├── plugins/          # Plugin system
│   └── services/         # Core monitoring services
├── themes/               # Pre-built theme configurations
├── config.json          # Main configuration file
├── plugins_config.json  # Plugin configuration
├── exports/             # Data export directory
└── logs/                # Log files directory
```

### 📋 Known Limitations (Alpha)
- GPU monitoring requires specific drivers for optimal functionality
- Some advanced features may require administrator privileges
- Plugin system is currently limited to built-in plugins

### 🚀 Future Roadmap
- Web interface for remote monitoring
- Alert system with notifications
- Historical data visualization
- Container monitoring enhancements
- Plugin marketplace and dynamic loading
- Network traffic analysis
- System benchmarking tools

### 🤝 Contributing
This is an open-source project welcoming contributions! See the [README.md](README.md) and [PLUGIN_USAGE_GUIDE.md](PLUGIN_USAGE_GUIDE.md) for development guidelines.

### 📝 License
MIT License - see [LICENSE](LICENSE) file for details.

---

**Thank you for using SysPulse!** 

Please report any issues, feature requests, or feedback through our GitHub repository. Your input helps make SysPulse better for everyone.

*Made with ❤️ by drclcomputers*
