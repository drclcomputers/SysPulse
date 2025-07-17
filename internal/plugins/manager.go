package plugins

import (
	"fmt"
	"sync"
	"time"

	"github.com/rivo/tview"
)

type PluginManager struct {
	plugins map[string]Plugin
	configs map[string]PluginConfig
	widgets map[string]tview.Primitive
	mutex   sync.RWMutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
		configs: make(map[string]PluginConfig),
		widgets: make(map[string]tview.Primitive),
	}
}

func (pm *PluginManager) LoadPlugin(plugin Plugin) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	name := plugin.Name()

	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("plugin %s already loaded", name)
	}

	config := PluginConfig{
		Name:     name,
		Enabled:  true,
		Settings: make(map[string]interface{}),
		Layout: WidgetConfig{
			Title:    name,
			Row:      0,
			Column:   0,
			RowSpan:  1,
			ColSpan:  1,
			MinWidth: 20,
			Enabled:  true,
		},
	}

	if err := plugin.Initialize(config); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %w", name, err)
	}

	pm.plugins[name] = plugin
	pm.configs[name] = config

	return nil
}

func (pm *PluginManager) LoadPluginWithConfig(plugin Plugin, config PluginConfig) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	name := plugin.Name()

	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("plugin %s already loaded", name)
	}

	config.Name = name

	if err := plugin.Initialize(config); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %w", name, err)
	}

	pm.plugins[name] = plugin
	pm.configs[name] = config

	return nil
}

func (pm *PluginManager) UnloadPlugin(name string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	if err := plugin.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown plugin %s: %w", name, err)
	}

	delete(pm.plugins, name)
	delete(pm.configs, name)
	delete(pm.widgets, name)

	return nil
}

func (pm *PluginManager) GetPlugin(name string) (Plugin, bool) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	plugin, exists := pm.plugins[name]
	return plugin, exists
}

func (pm *PluginManager) GetAllPlugins() []Plugin {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	plugins := make([]Plugin, 0, len(pm.plugins))
	for _, plugin := range pm.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

func (pm *PluginManager) CreateWidgets() map[string]tview.Primitive {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	widgets := make(map[string]tview.Primitive)

	for name, plugin := range pm.plugins {
		config := pm.configs[name]
		if config.Enabled {
			if widget, err := plugin.CreateWidget(); err == nil {
				widgets[name] = widget
				pm.widgets[name] = widget
			}
		}
	}

	return widgets
}

func (pm *PluginManager) UpdatePlugins() {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	for name, plugin := range pm.plugins {
		config := pm.configs[name]
		if config.Enabled {
			if widget, exists := pm.widgets[name]; exists {
				go plugin.UpdateWidget(widget)
			}
		}
	}
}

func (pm *PluginManager) EnablePlugin(name string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if _, exists := pm.plugins[name]; !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	config := pm.configs[name]
	config.Enabled = true
	pm.configs[name] = config

	return nil
}

func (pm *PluginManager) DisablePlugin(name string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if _, exists := pm.plugins[name]; !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	config := pm.configs[name]
	config.Enabled = false
	pm.configs[name] = config

	return nil
}

func (pm *PluginManager) UpdatePluginConfig(name string, config PluginConfig) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	if err := plugin.Initialize(config); err != nil {
		return fmt.Errorf("failed to update plugin %s config: %w", name, err)
	}

	pm.configs[name] = config
	return nil
}

func (pm *PluginManager) GetPluginConfig(name string) (PluginConfig, bool) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	config, exists := pm.configs[name]
	return config, exists
}

func (pm *PluginManager) GetPluginInfo() []PluginInfo {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	info := make([]PluginInfo, 0, len(pm.plugins))
	for name, plugin := range pm.plugins {
		config := pm.configs[name]
		widget, hasWidget := pm.widgets[name]

		var data map[string]interface{}
		if pluginData, err := plugin.CollectData(); err == nil {
			data = pluginData
		}

		pluginInfo := PluginInfo{
			Name:        plugin.Name(),
			Version:     plugin.Version(),
			Description: plugin.Description(),
			Author:      plugin.Author(),
			Config:      config,
			Widget:      widget,
			Data:        data,
			LastUpdate:  time.Now(),
		}

		if !hasWidget {
			pluginInfo.Widget = nil
		}

		info = append(info, pluginInfo)
	}
	return info
}

func (pm *PluginManager) IsPluginEnabled(name string) bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	config, exists := pm.configs[name]
	return exists && config.Enabled
}

func (pm *PluginManager) GetWidget(name string) (tview.Primitive, bool) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	widget, exists := pm.widgets[name]
	return widget, exists
}

func (pm *PluginManager) RefreshPlugin(name string) error {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	config := pm.configs[name]
	if !config.Enabled {
		return fmt.Errorf("plugin %s is disabled", name)
	}

	widget, hasWidget := pm.widgets[name]
	if !hasWidget {
		return fmt.Errorf("plugin %s has no widget", name)
	}

	return plugin.UpdateWidget(widget)
}
