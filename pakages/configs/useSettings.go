package configs

import (
	"encoding/json"
)

type SettingOptions struct {
	SubscriptionLink string `json:"SubscriptionLink"`
}

var (
	Settings SettingOptions
)

func LoadSettings() error {
	// Load the configuration from the file.
	configStr, UseConfigErr := UseConfig("settings")
	if UseConfigErr != nil {
		return UseConfigErr
	}

	// Parse the JSON configuration data into the SettingOptions struct.
	var options SettingOptions
	UnmarshalErr := json.Unmarshal([]byte(configStr), &options)
	if UnmarshalErr != nil {
		return UnmarshalErr
	}

	// Assign parsed settings to the global Settings variable.
	Settings = options
	return nil
}
