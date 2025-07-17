package plugins

import (
	"fmt"
	"syspulse/internal/utils"
	"time"

	"github.com/rivo/tview"
)

type ExamplePlugin struct {
	config PluginConfig
	data   map[string]interface{}
}

func NewExamplePlugin() *ExamplePlugin {
	return &ExamplePlugin{
		data: make(map[string]interface{}),
	}
}

func (p *ExamplePlugin) Name() string {
	return "Example Plugin"
}

func (p *ExamplePlugin) Version() string {
	return "1.0.0"
}

func (p *ExamplePlugin) Description() string {
	return "A simple example plugin that displays current time and custom data"
}

func (p *ExamplePlugin) Author() string {
	return "drclcomputers @ SysPulse"
}

func (p *ExamplePlugin) Initialize(config PluginConfig) error {
	p.config = config
	p.data["initialized"] = true
	p.data["init_time"] = time.Now().Format("15:04:05")
	return nil
}

func (p *ExamplePlugin) Shutdown() error {
	p.data["shutdown_time"] = time.Now().Format("15:04:05")
	return nil
}

func (p *ExamplePlugin) CreateWidget() (tview.Primitive, error) {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle(p.config.Layout.Title)
	textView.SetDynamicColors(true)

	if p.config.Layout.BorderColor != "" {
		textView.SetBorderColor(utils.GetColorFromName(p.config.Layout.BorderColor))
	}

	if p.config.Layout.ForegroundColor != "" {
		textView.SetTitleColor(utils.GetColorFromName(p.config.Layout.ForegroundColor))
	}

	if err := p.UpdateWidget(textView); err != nil {
		return nil, err
	}

	return textView, nil
}

func (p *ExamplePlugin) UpdateWidget(widget tview.Primitive) error {
	textView, ok := widget.(*tview.TextView)
	if !ok {
		return fmt.Errorf("widget is not a TextView")
	}

	currentTime := time.Now().Format("15:04:05")
	p.data["current_time"] = currentTime

	content := "[green]Example Plugin[white]\n"
	content += fmt.Sprintf("Current Time: %s\n", currentTime)
	content += fmt.Sprintf("Initialized: %v\n", p.data["initialized"])
	content += fmt.Sprintf("Init Time: %s\n", p.data["init_time"])
	content += fmt.Sprintf("Update Count: %d\n", p.getUpdateCount())

	if len(p.config.Settings) > 0 {
		content += "\n[yellow]Settings:[white]\n"
		for key, value := range p.config.Settings {
			content += fmt.Sprintf("  %s: %v\n", key, value)
		}
	}

	textView.SetText(content)
	return nil
}

func (p *ExamplePlugin) GetWidgetConfig() WidgetConfig {
	return p.config.Layout
}

func (p *ExamplePlugin) CollectData() (map[string]interface{}, error) {
	p.incrementUpdateCount()

	data := make(map[string]interface{})
	for k, v := range p.data {
		data[k] = v
	}

	return data, nil
}

func (p *ExamplePlugin) ExportData() map[string]interface{} {
	data, _ := p.CollectData()
	return data
}

func (p *ExamplePlugin) UpdateInterval() time.Duration {
	return 2 * time.Second
}

func (p *ExamplePlugin) getUpdateCount() int {
	if count, exists := p.data["update_count"]; exists {
		if c, ok := count.(int); ok {
			return c
		}
	}
	return 0
}

func (p *ExamplePlugin) incrementUpdateCount() {
	count := p.getUpdateCount()
	p.data["update_count"] = count + 1
}
