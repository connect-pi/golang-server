package configs

import (
	"encoding/json"
	"sync"
)

type SettingOptions struct {
	SubscriptionLink   string `json:"subscriptionLink"`
	UpdateSubsInTurnOn bool   `json:"updateSubsInTurnOn"`
}

var (
	settings SettingOptions
	loaded   bool
	mu       sync.Mutex
)

func UseSettings() (*SettingOptions, error) {
	mu.Lock()
	defer mu.Unlock()

	// If the settings have already been loaded, return them immediately.
	if loaded {
		return &settings, nil
	}

	// Load the configuration from the file.
	configStr, err := UseConfig("settings")
	if err != nil {
		return nil, err
	}

	// Parse the JSON configuration data into the SettingOptions struct.
	var options SettingOptions
	err = json.Unmarshal([]byte(configStr), &options)
	if err != nil {
		return nil, err
	}

	// Cache the loaded settings and mark them as loaded.
	settings = options
	loaded = true

	return &settings, nil
}
