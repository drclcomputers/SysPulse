package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type PluginSystemConfig struct {
	Plugins        map[string]PluginConfig `json:"plugins"`
	PluginSettings PluginSettings          `json:"plugin_settings"`
}

type PluginSettings struct {
	AutoLoad            bool   `json:"auto_load"`
	UpdateInterval      int    `json:"update_interval"`
	EnableNotifications bool   `json:"enable_notifications"`
	PluginDirectory     string `json:"plugin_directory"`
}

func LoadPluginConfig(configPath string) (*PluginSystemConfig, error) {
	if configPath == "" {
		configPath = "plugins_config.json"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &PluginSystemConfig{
			Plugins: make(map[string]PluginConfig),
			PluginSettings: PluginSettings{
				AutoLoad:            false,
				UpdateInterval:      2,
				EnableNotifications: false,
				PluginDirectory:     "./plugins",
			},
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin config file: %w", err)
	}

	var config PluginSystemConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse plugin config file: %w", err)
	}

	if config.PluginSettings.UpdateInterval == 0 {
		config.PluginSettings.UpdateInterval = 2
	}
	if config.PluginSettings.PluginDirectory == "" {
		config.PluginSettings.PluginDirectory = "./plugins"
	}

	for pluginName, pluginConfig := range config.Plugins {
		if pluginConfig.Layout.UpdateInterval == 0 {
			pluginConfig.Layout.UpdateInterval = 5
			config.Plugins[pluginName] = pluginConfig
		}
	}

	return &config, nil
}

func SavePluginConfig(config *PluginSystemConfig, configPath string) error {
	if configPath == "" {
		configPath = "plugins_config.json"
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal plugin config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write plugin config file: %w", err)
	}

	return nil
}

func IsPluginEnabledInConfig(config *PluginSystemConfig, pluginName string) bool {
	if config == nil || config.Plugins == nil {
		return false
	}

	pluginConfig, exists := config.Plugins[pluginName]
	return exists && pluginConfig.Enabled
}

func GetPluginConfigFromFile(config *PluginSystemConfig, pluginName string) (PluginConfig, bool) {
	if config == nil || config.Plugins == nil {
		return PluginConfig{}, false
	}

	pluginConfig, exists := config.Plugins[pluginName]
	return pluginConfig, exists
}
