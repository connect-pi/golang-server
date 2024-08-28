package configs

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFiles() error {
	// Path to the configs directory
	configDir := ".configs"

	// Check if the directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// If the directory does not exist, create it
		err := os.Mkdir(configDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	}

	// Path to the custom-rules.json file
	customRulesFile := filepath.Join(configDir, "custom-rules.json")
	if _, err := os.Stat(customRulesFile); os.IsNotExist(err) {
		// JSON content to be written to the file
		jsonData := `
{
      "domain": {
            "on": [],
            "off": []
      },
      "ip": {
            "on": [],
            "off": []
      }
}`

		// Write the JSON content to the file
		err := os.WriteFile(customRulesFile, []byte(jsonData), 0644)
		if err != nil {
			return fmt.Errorf("error writing custom-rules.json: %v", err)
		}

		fmt.Println("The custom-rules.json file was successfully created.")
	} else {
		fmt.Println("The custom-rules.json file already exists.")
	}

	// Path to the settings.json file
	settingsFile := filepath.Join(configDir, "settings.json")
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		// JSON content to be written to the file
		settingjsonData := `
{
      "SubscriptionLink": "",
      "UpdateSubscription": true
}`

		// Write the JSON content to the file
		settingsErr := os.WriteFile(settingsFile, []byte(settingjsonData), 0644)
		if settingsErr != nil {
			return fmt.Errorf("error writing settings.json: %v", settingsErr)
		}

		fmt.Println("The settings.json file was successfully created.")
	} else {
		fmt.Println("The settings.json file already exists.")
	}

	return nil
}
