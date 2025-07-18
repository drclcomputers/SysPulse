# SysPulse

A powerful terminal-based system monitoring tool written in Go, featuring real-time monitoring of CPU, memory, disk, network, and processes.

## 🌟 Features

- **Real-time System Monitoring**
  - CPU usage per core with load visualization
  - Memory (RAM and Swap) usage tracking
  - Disk usage and I/O statistics
  - Network activity monitoring
  - **GPU monitoring (cross-platform)** - NVIDIA, AMD, Intel support
  - Process management with search and filtering
  - Performance metrics tracking and self-monitoring
  - Automatic data export (CSV/JSON) with scheduled exports
  - Advanced logging system with rotation and multiple severity levels
  - **Comprehensive widget information modals** - Detailed system insights accessible with 'I'/'Enter' keys

- **Plugin System**
  - **Extensible Architecture** - Add custom widgets and monitoring capabilities
  - **Widget Integration** - Plugins create custom widgets that integrate with the main dashboard
  - **Real-time Updates** - Plugin widgets update automatically with configurable intervals
  - **Configuration Management** - JSON-based configuration for plugin settings and layout
  - **Data Collection** - Plugins can collect and export custom monitoring data
  - **Built-in Plugins** - Example and Docker monitoring plugins included

- **Beautiful Terminal UI**
  - Intuitive keyboard-driven interface
  - Fully customizable themes and layouts via JSON configuration
  - Real-time updates with configurable refresh rates
  - Interactive process management with kill functionality
  - Responsive design that adapts to terminal size
  - Smart widget cycling that follows screen layout order
  - Mouse support for modern terminals

- **Advanced Process Management**
  - Process filtering by CPU/Memory usage
  - Quick search functionality with real-time filtering
  - Detailed process information view
  - Safe process termination with confirmation
  - Process sorting by various metrics

- **Data Export & Analytics**
  - Automatic periodic data export (every 5 minutes)
  - Final export on application shutdown
  - CSV and JSON format support
  - Historical data tracking
  - Performance metrics for system optimization
  - Plugin data integration in exports

- **Enterprise-Ready Features**
  - Comprehensive logging with automatic rotation
  - Performance self-monitoring
  - Error handling with detailed context
  - Configuration validation
  - Cross-platform compatibility
  - Plugin lifecycle management

## 📦 Installation

### Prerequisites

- Go 1.24 or higher
- Terminal with Unicode and color support

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/syspulse.git

# Navigate to the project directory
cd syspulse

# Build the project
go build

# Run SysPulse
./syspulse
```

### Using Go Install

```bash
go install github.com/yourusername/syspulse@latest
```

## 🎮 Usage

### Keyboard Shortcuts

#### Global Controls
- `TAB`/`Shift+TAB` - Cycle through widgets in screen order (top-to-bottom, left-to-right)
- `ESC` - Unfocus current widget
- `Q` - Quit application
- `H` - Show help screen
- `I` - Show system information (detailed modal)

#### Widget Navigation
- `C` - Focus CPU widget
- `M` - Focus Memory widget
- `D` - Focus Disk widget
- `N` - Focus Network widget
- `P` - Focus Process widget
- `G` - Focus GPU widget

#### Process Management
- `K` - Kill selected process (platform-specific methods with confirmation)
- `F` - Search/filter processes
- `Y` - Toggle process sorting (CPU/Memory)
- `Up/Down` or `W/S` - Navigate process list
- `I` - View detailed process information

#### Process Kill Methods
- **Windows**: Graceful termination → Taskkill → Windows API
- **Linux/Unix**: SIGTERM → SIGKILL with signal handling
- **Cross-platform**: Fallback to basic kill method

#### System Information
- `I` (on any widget) - Show detailed information for that component
- `I` (on CPU widget) - Show CPU specifications, current usage, and per-core statistics
- `I` (on Memory widget) - Show RAM/Swap usage, health status, and optimization tips
- `I` (on Disk widget) - Show per-partition information, usage statistics, and health advice
- `I` (on Network widget) - Show interface details, transfer rates, and network statistics
- `I` (on GPU widget) - Show GPU details and driver information

## ⚙️ Configuration

SysPulse uses a flexible JSON configuration system. Configuration files are loaded in this order:
1. `config.json` in the root directory
2. `internal/services/UI/default.json` (fallback)
3. `plugins_config.json` (plugins configuration)

### Example Configuration

```json
{
  "background": "black",
  "foreground": "white",
  "altforeground": "grey",
  "cpu": {
    "bar_low": "green",
    "bar_high": "red"
  },
  "memory": {
    "vmem_gauge": "blue",
    "smem_gauge": "cyan"
  },
  "network": {
    "bar_low": "yellow",
    "bar_high": "purple"
  },
  "disk": {
    "bar_low": "blue",
    "bar_medium": "yellow",
    "bar_high": "red",
    "bar_empty": "white"
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
      "colSpan": 1,
      "minWidth": 30,
      "weight": 1.0
    },
    "memory": {
      "enabled": true,
      "row": 0,
      "column": 1,
      "rowSpan": 1,
      "colSpan": 1,
      "minWidth": 30,
      "weight": 1.0
    },
    "disk": {
      "enabled": true,
      "row": 1,
      "column": 1,
      "rowSpan": 1,
      "colSpan": 1,
      "minWidth": 30,
      "weight": 1.0
    },
    "network": {
      "enabled": true,
      "row": 2,
      "column": 0,
      "rowSpan": 1,
      "colSpan": 1,
      "minWidth": 30,
      "weight": 1.0
    },
    "process": {
      "enabled": true,
      "row": 2,
      "column": 1,
      "rowSpan": 2,
      "colSpan": 1,
      "minWidth": 30,
      "weight": 1.0
    },
    "gpu": {
      "enabled": true,
      "row": 3,
      "column": 0,
      "rowSpan": 1,
      "colSpan": 1,
      "minWidth": 30,
      "weight": 1.0
    }
  },
  "processsort": "cpu",
  "updatetime": 1,
	"export": {
		"enabled": true,
		"interval": 300,
		"formats": ["csv", "json"],
		"directory": "exports",
		"filename_prefix": "syspulse"
	}
}
```

### Configuration Options

#### Theme Colors
- **Background/Foreground**: Basic terminal colors
- **Component Colors**: Each widget has customizable color schemes
- **Bar Colors**: Different colors for usage thresholds (low, medium, high)

#### Layout System
- **Grid-based**: Configurable rows and columns
- **Widget Positioning**: Precise control over widget placement
- **Responsive**: Automatic sizing based on terminal dimensions
- **Smart Focus Cycling**: Widget cycling follows visual layout order
- **Enable/Disable**: Toggle individual widgets on/off

#### Update Settings
- **Refresh Rate**: Configurable update interval (in seconds)
- **Process Sorting**: Default sort method (cpu/memory)
- **Data Export**: Automatic export scheduling

#### GPU Configuration
- **Cross-platform**: Works on Windows, Linux, and macOS
- **Auto-detection**: Automatically detects NVIDIA, AMD, and Intel GPUs

## 🛠️ Development

### Project Structure

```
syspulse/
├── cmd/                     # Command-line interface
├── internal/                # Internal packages
│   ├── errors/             # Error handling and types
│   ├── export/             # Data export functionality (CSV/JSON)
│   ├── logger/             # Logging system
│   │   └── v2/            # Advanced logging with rotation
│   ├── metrics/            # Performance monitoring
│   ├── plugins/            # Plugin system
│   │   ├── interface.go   # Plugin interface definition
│   │   ├── manager.go     # Plugin manager
│   │   ├── integration.go # Dashboard integration
│   │   ├── example.go     # Example plugin
│   │   └── docker.go      # Docker monitoring plugin
│   └── services/           # Core monitoring services
│       ├── cpu/           # CPU monitoring and statistics
│       ├── disk/          # Disk usage and I/O monitoring
│       ├── gpu/           # GPU monitoring (cross-platform)
│       ├── memory/        # Memory usage tracking
│       ├── network/       # Network interface monitoring
│       ├── processes/     # Process management
│       ├── sysinfo/       # System information gathering
│       └── ui/            # Terminal user interface
│       └── utils/              # Utility functions and models
├── config.json            # Configuration file
├── plugins_config.json    # Plugin configuration
├── exports/               # Exported data directory
├── logs/                  # Log files directory
├── PLUGIN_USAGE_GUIDE.md  # Detailed plugin development guide
└── README.md              # This file
```

### Key Features in Detail

#### Widget Information Modals
- **Comprehensive system insights**: Each widget provides detailed information accessible with the 'I' key
- **Memory modal**: Shows RAM/Swap usage, health status, and memory optimization tips
- **Disk modal**: Displays per-partition information, usage statistics, and health recommendations
- **Network modal**: Provides interface details, transfer rates, and network performance metrics
- **CPU modal**: Hardware specifications, current usage, per-core breakdown, and thermal status
- **Consistent interface**: Unified modal design across all widgets for familiar user experience
- **Context-aware help**: Each modal provides specific tips and recommendations based on current system state

#### Plugin System
- **Extensible Architecture**: Add custom widgets and monitoring capabilities without modifying core code
- **Widget Integration**: Plugins create custom tview widgets that integrate seamlessly with the main dashboard
- **Real-time Updates**: Plugin widgets update automatically with configurable intervals
- **Configuration Management**: JSON-based configuration for plugin settings, layout, and positioning
- **Smart Focus Integration**: Plugin widgets participate in the intelligent focus cycling system
- **Data Collection**: Plugins can collect custom metrics and export data for analysis
- **Built-in Examples**: Example and Docker monitoring plugins included as templates
- **Lifecycle Management**: Proper initialization, update, and cleanup methods for plugins

#### GPU Monitoring
- **Cross-platform support**: Windows (WMI), Linux (nvidia-smi, sysfs), macOS (system_profiler)
- **Multi-vendor**: NVIDIA, AMD, Intel GPU detection
- **Comprehensive metrics**: Temperature, memory usage, utilization, driver info
- **Graceful fallback**: Continues working even if GPU monitoring fails

#### Data Export System
- **Multiple formats**: CSV for analysis, JSON for programmatic use
- **Historical tracking**: Maintains comprehensive system metrics over time
- **Final export**: Saves data on application shutdown

#### Advanced Logging
- **Multiple severity levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Automatic rotation**: Daily log file rotation
- **Contextual information**: File, line, function, and timestamp
- **Dual output**: Console and file logging simultaneously

### Building from Source

1. **Clone the repository**
   ```bash
   git clone https://github.com/drclcomputers/syspulse.git
   cd syspulse
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Build the project**
   ```bash
   go build -o syspulse
   ```

4. **Run SysPulse**
   ```bash
   ./syspulse
   ```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/export/
go test ./internal/logger/v2/
go test ./internal/metrics/
```

### Development Guidelines

#### Code Organization
- Use build tags for platform-specific code (Windows, Linux, macOS)
- Implement comprehensive error handling with context
- Follow Go best practices and idioms
- Write tests for all new functionality
- Maintain consistent widget positioning and cycling behavior
- Ensure plugin widgets integrate properly with the layout system

#### Testing Strategy
- Unit tests for individual components
- Integration tests for data export and logging
- Platform-specific tests for GPU monitoring
- Performance tests for resource usage

#### Performance Considerations
- Monitor memory usage and optimize hot paths
- Use goroutines efficiently for concurrent updates
- Implement proper error handling to prevent crashes
- Profile the application under various system loads

## 📈 Data Export

SysPulse automatically exports monitoring data for analysis and archival:

### Export Features
- **Automatic scheduling**: Exports every 5 minutes during operation
- **Final export**: Comprehensive data dump on application shutdown
- **Multiple formats**: CSV for spreadsheet analysis, JSON for programmatic use
- **Comprehensive metrics**: CPU, memory, disk, network, and process data
- **Plugin data integration**: Plugin-collected data included in exports

### Export Location
- **Directory**: `exports/` in the project root
- **Naming convention**: `syspulse_YYYY-MM-DD_HH-MM-SS.csv/json`
- **Final exports**: `syspulse_final_YYYY-MM-DD_HH-MM-SS.csv/json`

### CSV Format
```csv
Timestamp,CPU_Total,Memory_Total,Memory_Used,Swap_Total,Swap_Used,
Disk_Path,Disk_Total,Disk_Used,Disk_UsedPerc,Disk_IOReads,Disk_IOWrites,
Net_BytesSent,Net_BytesReceived,Net_PacketsSent,Net_PacketsReceived
```

### JSON Format
```json
[
  {
    "Timestamp": "2025-07-15T12:30:00Z",
    "CPU": [15.2, 12.8, 18.5, 10.1],
    "Memory": {
      "Total": 16777216000,
      "Used": 8388608000,
      "SwapTotal": 2147483648,
      "SwapUsed": 0
    },
    "Disk": {
      "Path": "/",
      "Total": 1000000000000,
      "Used": 500000000000,
      "UsedPerc": 50.0,
      "IOReads": 12345,
      "IOWrites": 67890
    },
    "Network": {
      "BytesSent": 1024000,
      "BytesReceived": 2048000,
      "PacketsSent": 1000,
      "PacketsRecv": 1500
    },
    "Plugins": {
      "example_plugin": {
        "timestamp": 1642248600,
        "status": "active",
        "update_count": 150
      },
      "docker_plugin": {
        "containers_running": 5,
        "containers_stopped": 2,
        "images_count": 12
      }
    }
  }
]
```

## 🔧 Advanced Usage

### Custom Themes
Create custom theme files and load them at runtime:

```bash
# Use custom configuration
./syspulse --config custom-theme.json

# Export current configuration
./syspulse --export-config theme-backup.json
```

### Monitoring Specific Components
Enable only the components you need:

```json
{
  "layout": {
    "cpu": {"enabled": true},
    "memory": {"enabled": true},
    "disk": {"enabled": false},
    "network": {"enabled": false},
    "process": {"enabled": true},
    "gpu": {"enabled": true}
  }
}
```

### Performance Tuning
Adjust update intervals based on your needs:

```json
{
  "updatetime": 2,     // Update every 2 seconds (default: 1)
  "processsort": "mem" // Sort processes by memory usage
}
```

## 🔌 Plugin System

SysPulse features a powerful plugin system that allows you to extend functionality by adding custom widgets and monitoring capabilities.

### Plugin Architecture

The plugin system is built around a simple interface that allows plugins to:

1. **Initialize** with custom configuration
2. **Create widgets** for the UI
3. **Export data** for analysis
4. **Clean up** resources when unloaded

### Built-in Plugins

#### Example Plugin
- **Purpose**: Demonstrates basic plugin functionality
- **Features**: Shows current time, initialization status, and update count
- **Configuration**: Customizable message and display options
- **Status**: Enabled by default

#### Docker Plugin
- **Purpose**: Monitors Docker containers and images
- **Features**: Shows running/stopped containers, recent images, and system stats
- **Requirements**: Docker must be installed and running
- **Configuration**: Configurable limits for containers and images shown
- **Status**: Disabled by default

### Plugin Configuration

Plugins can be configured through the `plugins_config.json` file:

```json
{
  "plugins": {
    "example": {
      "name": "Example Plugin",
      "enabled": true,
      "settings": {
        "show_time": true,
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
    },
    "docker": {
      "name": "Docker Monitor",
      "enabled": false,
      "settings": {
        "show_containers": true,
        "show_images": true,
        "container_limit": 5
      },
      "layout": {
        "title": "Docker",
        "row": 1,
        "column": 5,
        "rowSpan": 2,
        "colSpan": 2,
        "minWidth": 30,
        "enabled": true
      }
    }
  }
}
```

### Creating a Custom Plugin

To create a new plugin, implement the `Plugin` interface:

```go
type Plugin interface {
    // Metadata
    Name() string
    Version() string
    Description() string
    Author() string
    
    // Plugin lifecycle
    Initialize() error
    Cleanup() error
    
    // Widget management
    CreateWidget() (tview.Primitive, error)
    UpdateWidget(widget tview.Primitive) error
    
    // Data collection and export
    CollectData() (map[string]interface{}, error)
    ExportData() interface{}
}
```

### Example Plugin Implementation

```go
package plugins

import (
    "fmt"
    "time"
    "github.com/rivo/tview"
)

type MyPlugin struct {
    widget *tview.TextView
    data   map[string]interface{}
}

func NewMyPlugin() *MyPlugin {
    return &MyPlugin{
        data: make(map[string]interface{}),
    }
}

func (p *MyPlugin) Name() string { return "My Custom Plugin" }
func (p *MyPlugin) Version() string { return "1.0.0" }
func (p *MyPlugin) Description() string { return "A custom plugin for monitoring" }
func (p *MyPlugin) Author() string { return "Your Name" }

func (p *MyPlugin) Initialize() error {
    return nil
}

func (p *MyPlugin) CreateWidget() (tview.Primitive, error) {
    p.widget = tview.NewTextView()
    p.widget.SetBorder(true)
    p.widget.SetTitle("My Plugin")
    p.widget.SetDynamicColors(true)
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

### Plugin Integration

To integrate your plugin with SysPulse:

1. **Add your plugin** to the `InitializePluginSystem` function in `internal/plugins/integration.go`
2. **Configure the plugin** in `plugins_config.json` with layout positioning and settings
3. **Build and run** SysPulse to see your plugin in action

For detailed plugin development instructions, see the [Plugin Usage Guide](PLUGIN_USAGE_GUIDE.md).

### Plugin Ideas

Here are some ideas for plugins that could be created (if you're interested in developing a plugin):

- **System Services Monitor** - Monitor systemd services, Windows services, etc.
- **Database Monitor** - Monitor database connections, queries, performance
- **Web Server Monitor** - Monitor Apache, Nginx, or other web servers
- **Log Monitor** - Tail and analyze log files
- **Weather Widget** - Display weather information
- **Stock Monitor** - Monitor stock prices and financial data
- **Kubernetes Monitor** - Monitor Kubernetes pods, services, and nodes
- **Git Monitor** - Monitor Git repositories for changes
- **Package Manager** - Monitor system packages and updates
- **Custom Sensors** - Monitor custom hardware sensors or APIs

## 🤝 Contributing

I welcome contributions! Here's how to get started:

### Development Setup
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Set up your development environment:
   ```bash
   go mod download
   go install golang.org/x/tools/cmd/goimports@latest
   ```

### Making Changes
1. Write comprehensive tests for new features
2. Ensure all tests pass (`go test ./...`)
3. Format your code (`goimports -w .`)
4. Update documentation as needed
5. Commit your changes (`git commit -m 'Add some amazing feature'`)

### Pull Request Process
1. Push to your feature branch (`git push origin feature/amazing-feature`)
2. Open a Pull Request with:
   - Clear description of changes
   - Screenshots for UI changes
   - Test results
   - Performance impact analysis

### Code Style Guidelines
- Follow Go best practices and idioms
- Use meaningful variable and function names
- Add comments for complex logic
- Include build tags for platform-specific code
- Handle errors gracefully with proper context
- Write comprehensive tests with good coverage

### Areas for Contribution
- **Platform support**: Improve GPU monitoring for different platforms
- **Performance**: Optimize resource usage and update speeds
- **Features**: Add new monitoring capabilities
- **Plugin Development**: Create new plugins for different monitoring needs
- **UI/UX**: Enhance the terminal interface and widget cycling
- **Documentation**: Improve examples and guides
- **Testing**: Expand test coverage and add benchmarks
- **Layout System**: Enhance the grid-based layout capabilities

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## � Recent Updates

### v1.2.0 - Widget Cycling Improvements
- **Smart Widget Cycling**: TAB navigation now follows visual layout order (top-to-bottom, left-to-right)
- **Plugin Integration**: Plugin widgets seamlessly integrate into the focus cycling system
- **Layout-Aware Navigation**: Widget cycling respects grid positioning and enabled states
- **Improved User Experience**: More intuitive navigation that matches visual expectations

### v1.1.0 - Export System Fixes
- **GPU Export Issues**: Fixed memory overflow and missing GPU data in exports
- **Data Validation**: Improved bounds checking for GPU memory values
- **Export Consistency**: Better handling of multiple GPU configurations

## �🙏 Acknowledgments

![MIT License](https://img.shields.io/badge/License-MIT-green.svg)
![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)
![Platform](https://img.shields.io/badge/platform-windows%20%7C%20linux%20%7C%20macos-lightgrey.svg)

### Libraries and Dependencies
- [tview](https://github.com/rivo/tview) - Advanced terminal UI library
- [gopsutil](https://github.com/shirou/gopsutil) - Cross-platform system monitoring
- [tcell](https://github.com/gdamore/tcell) - Terminal handling and input
- [cobra](https://github.com/spf13/cobra) - Modern CLI framework
- [wmi](https://github.com/StackExchange/wmi) - Windows Management Instrumentation

### Special Thanks
- **Go Team**: For creating an excellent systems programming language
- **Terminal UI Community**: For inspiration and best practices
- **System Monitoring Tools**: htop, btop++, glances for setting the standard

### Features Roadmap
- [ ] Web interface for remote monitoring
- [x] **Plugin system for custom metrics** ✅
- [x] **Smart widget cycling based on layout** ✅
- [ ] Alert system with notifications
- [ ] Historical data visualization
- [ ] Container monitoring support
- [ ] Network traffic analysis
- [ ] System benchmarking tools
- [ ] Plugin marketplace and repository
- [ ] Dynamic plugin loading from external files
- [ ] Plugin configuration UI within the dashboard

## Actively Known Bugs / Features to include

### These are bugs / features that have been detected by SysPulse's developer or by the users. Currently, there is active development in fixing them / adding them.

- [x] **Stricter plugin widgets verification** ✅
- [x] **Platform-specific process killing** ✅
- [x] **No GPU utilisation in Windows** !!! Probably. I don't know exactly. Getting GPU info without DX or other libraries is such a hassle !!!
- [x] **When the terminal sizes are too small (even excesive zooming in could cause this), no widgets would be shown. -> Fix: make your terminal bigger by zooming out, making it fullscreen or make the min width in config.json smaller** ✅
- [x] **Temperature monitoring rarely works, even in Linux** - Still don't know exactly if it completely works. It's a hassle to get info about temperature in Windows.

If you find any bugs, open an ISSUE and I'll do my best in fixing them as quickly as possible.

## Changelog

- 11.9 alpha - Added custom update intervals for widgets + Bug Fixes
- 11.7 alpha - Initial Release

---

**Made with ❤️ by the drclcomputers**

*For issues, feature requests, and plugin contributions, visit our [GitHub repository](https://github.com/drclcomputers/syspulse)*