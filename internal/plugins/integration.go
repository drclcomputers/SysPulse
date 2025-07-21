package plugins

import (
	"fmt"
	"log"
	"syspulse/internal/utils"
	"time"

	"github.com/rivo/tview"
)

func InitializePluginSystem(dashboard *utils.Dashboard) error {
	dashboard.PluginManager = NewPluginManager()
	dashboard.PluginWidgets = make(map[string]tview.Primitive)

	pluginManager, ok := dashboard.PluginManager.(*PluginManager)
	if !ok {
		return fmt.Errorf("failed to create plugin manager")
	}

	pluginConfig, err := LoadPluginConfig("plugins_config.json")
	if err != nil {
		log.Printf("Failed to load plugin config: %v", err)
		pluginConfig = &PluginSystemConfig{
			Plugins: make(map[string]PluginConfig),
			PluginSettings: PluginSettings{
				AutoLoad:            false,
				UpdateInterval:      2,
				EnableNotifications: false,
				PluginDirectory:     "./plugins",
			},
		}
	}

	if IsPluginEnabledInConfig(pluginConfig, "example") {
		examplePlugin := NewExamplePlugin()
		if config, exists := GetPluginConfigFromFile(pluginConfig, "example"); exists {
			if err := pluginManager.LoadPluginWithConfig(examplePlugin, config); err != nil {
				log.Printf("Failed to load example plugin: %v", err)
			}
		} else {
			if err := pluginManager.LoadPlugin(examplePlugin); err != nil {
				log.Printf("Failed to load example plugin: %v", err)
			}
		}
	}

	if IsPluginEnabledInConfig(pluginConfig, "docker") {
		dockerPlugin := NewDockerPlugin()
		if config, exists := GetPluginConfigFromFile(pluginConfig, "docker"); exists {
			if err := pluginManager.LoadPluginWithConfig(dockerPlugin, config); err != nil {
				log.Printf("Failed to load Docker plugin: %v", err)
			}
		} else {
			if err := pluginManager.LoadPlugin(dockerPlugin); err != nil {
				log.Printf("Failed to load Docker plugin: %v", err)
			}
		}
	}

	widgets := pluginManager.CreateWidgets()
	for name, widget := range widgets {
		dashboard.PluginWidgets[name] = widget
	}

	return nil
}

func StartPluginUpdateWorker(dashboard *utils.Dashboard) {
	if dashboard.PluginManager == nil {
		return
	}

	pluginManager, ok := dashboard.PluginManager.(*PluginManager)
	if !ok {
		return
	}

	pluginInfo := pluginManager.GetPluginInfo()
	for _, info := range pluginInfo {
		if info.Config.Enabled {
			plugin, exists := pluginManager.GetPlugin(info.Name)
			if exists {
				startPluginWorker(dashboard, plugin, info.Config)
			}
		}
	}
}

func startPluginWorker(dashboard *utils.Dashboard, plugin Plugin, config PluginConfig) {
	interval := config.Layout.UpdateInterval
	if interval <= 0 {
		interval = 5
	}

	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if dashboard.PluginManager != nil {
				if pluginManager, ok := dashboard.PluginManager.(*PluginManager); ok {
					if widget, exists := pluginManager.GetWidget(plugin.Name()); exists {
						plugin.UpdateWidget(widget)
					}
				}
			}
		}
	}()
}

func AddPluginWidgetsToLayout(dashboard *utils.Dashboard, mainFlex *tview.Flex) {
	if dashboard.PluginManager == nil {
		return
	}

	pluginManager, ok := dashboard.PluginManager.(*PluginManager)
	if !ok {
		return
	}

	pluginInfo := pluginManager.GetPluginInfo()

	for _, info := range pluginInfo {
		if info.Config.Enabled && info.Widget != nil {
			container := tview.NewFlex().SetDirection(tview.FlexRow)
			container.AddItem(info.Widget, 0, 1, false)
			utils.SetBorderStyle(container.Box)
			container.SetTitle(info.Name + " v" + info.Version)

			mainFlex.AddItem(container, 0, 1, false)
		}
	}
}

func AddPluginWidgetsToGrid(dashboard *utils.Dashboard, grid *tview.Grid) {
	if dashboard.PluginManager == nil {
		return
	}

	pluginManager, ok := dashboard.PluginManager.(*PluginManager)
	if !ok {
		return
	}

	pluginInfo := pluginManager.GetPluginInfo()

	for _, info := range pluginInfo {
		if info.Config.Enabled && info.Widget != nil {
			config := info.Config.Layout
			grid.AddItem(info.Widget,
				config.Row, config.Column,
				config.RowSpan, config.ColSpan,
				config.MinWidth, 0, false)
		}
	}
}

func ShowPluginManagerModal(dashboard *utils.Dashboard) {
	if dashboard.PluginManager == nil {
		return
	}

	pluginManager, ok := dashboard.PluginManager.(*PluginManager)
	if !ok {
		return
	}

	modal := tview.NewModal()
	modal.SetTitle("Plugin Manager")
	modal.SetBorder(true)

	pluginInfo := pluginManager.GetPluginInfo()

	text := "Loaded Plugins:\n\n"
	for _, info := range pluginInfo {
		status := "Disabled"
		if info.Config.Enabled {
			status = "Enabled"
		}
		text += fmt.Sprintf("• %s v%s - %s\n", info.Name, info.Version, status)
		text += fmt.Sprintf("  %s\n", info.Description)
		text += fmt.Sprintf("  Author: %s\n\n", info.Author)
	}

	modal.SetText(text)
	modal.AddButtons([]string{"Close"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		dashboard.App.SetRoot(dashboard.MainWidget, true)
	})

	dashboard.App.SetRoot(modal, true)
}

func EnablePlugin(dashboard *utils.Dashboard, pluginName string) error {
	if dashboard.PluginManager == nil {
		return fmt.Errorf("plugin manager not initialized")
	}

	pluginManager, ok := dashboard.PluginManager.(*PluginManager)
	if !ok {
		return fmt.Errorf("invalid plugin manager type")
	}

	return pluginManager.EnablePlugin(pluginName)
}

func DisablePlugin(dashboard *utils.Dashboard, pluginName string) error {
	if dashboard.PluginManager == nil {
		return fmt.Errorf("plugin manager not initialized")
	}

	pluginManager, ok := dashboard.PluginManager.(*PluginManager)
	if !ok {
		return fmt.Errorf("invalid plugin manager type")
	}

	return pluginManager.DisablePlugin(pluginName)
}

func RefreshPlugin(dashboard *utils.Dashboard, pluginName string) error {
	if dashboard.PluginManager == nil {
		return fmt.Errorf("plugin manager not initialized")
	}

	pluginManager, ok := dashboard.PluginManager.(*PluginManager)
	if !ok {
		return fmt.Errorf("invalid plugin manager type")
	}

	return pluginManager.RefreshPlugin(pluginName)
}

func GetPluginData(dashboard *utils.Dashboard, pluginName string) (map[string]interface{}, error) {
	if dashboard.PluginManager == nil {
		return nil, fmt.Errorf("plugin manager not initialized")
	}

	pluginManager, ok := dashboard.PluginManager.(*PluginManager)
	if !ok {
		return nil, fmt.Errorf("invalid plugin manager type")
	}

	plugin, exists := pluginManager.GetPlugin(pluginName)
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginName)
	}

	return plugin.CollectData()
}

func ExportPluginData(dashboard *utils.Dashboard) map[string]interface{} {
	if dashboard.PluginManager == nil {
		return nil
	}

	pluginManager, ok := dashboard.PluginManager.(*PluginManager)
	if !ok {
		return nil
	}

	data := make(map[string]interface{})
	pluginInfo := pluginManager.GetPluginInfo()

	for _, info := range pluginInfo {
		if info.Config.Enabled {
			plugin, exists := pluginManager.GetPlugin(info.Name)
			if exists {
				data[info.Name] = plugin.ExportData()
			}
		}
	}

	return data
}
