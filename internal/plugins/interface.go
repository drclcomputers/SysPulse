package plugins

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Plugin interface {
	Name() string
	Version() string
	Description() string
	Author() string

	Initialize(config PluginConfig) error
	Shutdown() error

	CreateWidget() (tview.Primitive, error)
	UpdateWidget(widget tview.Primitive) error
	GetWidgetConfig() WidgetConfig

	CollectData() (map[string]interface{}, error)
	ExportData() map[string]interface{}
}

type PluginConfig struct {
	Name     string                 `json:"name"`
	Enabled  bool                   `json:"enabled"`
	Settings map[string]interface{} `json:"settings"`
	Layout   WidgetConfig           `json:"layout"`
}

type WidgetConfig struct {
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

type PluginInfo struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Author      string                 `json:"author"`
	Config      PluginConfig           `json:"config"`
	Widget      tview.Primitive        `json:"-"`
	Data        map[string]interface{} `json:"-"`
	LastUpdate  time.Time              `json:"last_update"`
}

type InputHandler interface {
	HandleInput(event *tcell.EventKey) *tcell.EventKey
}

type InfoProvider interface {
	GetDetailedInfo() string
	ShowInfoModal(app *tview.Application, returnWidget tview.Primitive) tview.Primitive
}

type FilterProvider interface {
	ApplyFilter(filterTerm string) error
	ClearFilter() error
}
