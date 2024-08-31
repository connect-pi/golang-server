package rules

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Has checks if a value exists in the On or Off slice
func (oo OnOff) Has(slice string, value string) bool {
	if slice == "on" {
		for _, v := range oo.On {
			if v == value {
				return true
			}
		}
	} else if slice == "off" {
		for _, v := range oo.Off {
			if v == value {
				return true
			}
		}
	}
	return false
}

// get custom rules on init
func LoadCustomRules() error {
	// Get the directory of the executable
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	// Construct the full path to the JSON file
	configPath := filepath.Join(filepath.Dir(path), "/golang-server/.configs/custom-rules.json")

	// Open the JSON file
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error opening file '%s': %v", configPath, err)
		return err
	}
	defer file.Close()

	// Decode JSON data into rules struct
	var rules Rules
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&rules); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
		return err
	}

	// Cache the rules data
	CustomRules = rules

	// Combine the Domain and IP rules from both default and custom rules
	CombinedRules = Rules{
		Domain: OnOff{
			On:  append(CustomRules.Domain.On, DefaultRules.Domain.On...),
			Off: append(CustomRules.Domain.Off, DefaultRules.Domain.Off...),
		},
		IP: OnOff{
			On:  append(CustomRules.IP.On, DefaultRules.IP.On...),
			Off: append(CustomRules.IP.Off, DefaultRules.IP.Off...),
		},
	}

	return nil
}
