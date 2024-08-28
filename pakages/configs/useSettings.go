package configs

import (
	"encoding/json"
	"sync"
)

type SettingOptions struct {
	SubscriptionLink   string `json:"subscriptionLink"`
	UpdateSubscription bool   `json:"updateSubsInTurnOn"`
}

var (
	Settings SettingOptions
	loaded   bool
	mu       sync.Mutex
)

func LoadSettings() error {
	mu.Lock()
	defer mu.Unlock()

	// If the settings have already been loaded, return them immediately.
	if loaded {
		return nil
	}

	// Load the configuration from the file.
	configStr, err := UseConfig("settings")
	if err != nil {
		return nil
	}

	// Parse the JSON configuration data into the SettingOptions struct.
	var options SettingOptions
	err = json.Unmarshal([]byte(configStr), &options)
	if err != nil {
		return nil
	}

	// Cache the loaded settings and mark them as loaded.
	Settings = options
	loaded = true
	return nil
}
