package utils

import (
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type CPUModel struct {
	BarLow  string `json:"bar_low"`
	BarHigh string `json:"bar_high"`
}

type MEMModel struct {
	VMemGauge string `json:"vmem_gauge"`
	SMemGauge string `json:"smem_gauge"`
}

type NETModel CPUModel

type DISKModel struct {
	BarLow    string `json:"bar_low"`
	BarMedium string `json:"bar_medium"`
	BarHigh   string `json:"bar_high"`
	BarEmpty  string `json:"bar_empty"`
}

type GPUModel struct {
	BarLow  string `json:"bar_low"`
	BarHigh string `json:"bar_high"`
}

type ExportConfig struct {
	Enabled        bool     `json:"enabled"`
	Interval       int      `json:"interval"`
	Formats        []string `json:"formats"`
	Directory      string   `json:"directory"`
	FilenamePrefix string   `json:"filename_prefix"`
}

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
}

type Widget struct {
	Enabled bool `json:"enabled"`
	Row     int  `json:"row"`
	Column  int  `json:"column"`
	RowSpan int  `json:"rowSpan"`
	ColSpan int  `json:"colSpan"`
}

type LayoutConfig struct {
	CPU          WidgetConfig `json:"cpu"`
	Memory       WidgetConfig `json:"memory"`
	Disk         WidgetConfig `json:"disk"`
	Network      WidgetConfig `json:"network"`
	Process      WidgetConfig `json:"process"`
	GPU          WidgetConfig `json:"gpu"`
	Load         WidgetConfig `json:"load"`
	Temperature  WidgetConfig `json:"temperature"`
	NetworkConns WidgetConfig `json:"network_connections"`
	DiskIO       WidgetConfig `json:"disk_io"`
	ProcessTree  WidgetConfig `json:"process_tree"`
	Battery      WidgetConfig `json:"battery"`
	Rows         int          `json:"rows"`
	Columns      int          `json:"columns"`
	Spacing      int          `json:"spacing"`
}

type Theme struct {
	Background    string       `json:"background"`
	Foreground    string       `json:"foreground"`
	Altforeground string       `json:"altforeground"`
	CPU           CPUModel     `json:"cpu"`
	Memory        MEMModel     `json:"memory"`
	Network       NETModel     `json:"network"`
	Disk          DISKModel    `json:"disk"`
	GPU           GPUModel     `json:"gpu"`
	Layout        LayoutConfig `json:"layout"`
	Sorting       string       `json:"processsort"`
	UpdateTime    int          `json:"updatetime"`
	Export        ExportConfig `json:"export"`
}

type Dashboard struct {
	App                *tview.Application
	CpuWidget          *tview.Box
	MemWidget          *tview.Box
	DiskWidget         *tview.Box
	NetWidget          *tview.Box
	ProcessWidget      *tview.List
	GPUWidget          *tview.Box
	LoadWidget         *tview.Box
	TemperatureWidget  *tview.Box
	NetworkConnsWidget *tview.Box
	DiskIOWidget       *tview.Box
	ProcessTreeWidget  *tview.Box
	BatteryWidget      *tview.Box
	MainWidget         *tview.Flex
	Theme              Theme
	CpuData            []float64
	VMemData           *mem.VirtualMemoryStat
	SMemData           *mem.SwapMemoryStat
	DiskData           []*disk.UsageStat
	NetData            *net.IOCountersStat
	Processes          []process.Process
	HeaderWidget       *tview.Box
	FooterWidget       *tview.TextView
	LoadData           interface{}
	TemperatureData    interface{}
	NetworkConnsData   interface{}
	DiskIOData         interface{}
	ProcessTreeData    interface{}
	BatteryData        interface{}
	GPUData            interface{}

	ProcessFilterActive bool
	ProcessFilterTerm   string
	ProcessFilterType   string

	PluginManager interface{}
	PluginWidgets map[string]tview.Primitive
}
