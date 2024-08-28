package configs

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFiles() {
	// Path to the configs directory
	configDir := ".configs"

	// Check if the directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// If the directory does not exist, create it
		err := os.Mkdir(configDir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	} else {
		return
	}

	// Path to the custom-rules.json file
	configFile := filepath.Join(configDir, "custom-rules.json")

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
	err := os.WriteFile(configFile, []byte(jsonData), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("The custom-rules.json file was successfully created.")
}
