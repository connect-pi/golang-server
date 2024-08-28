package configs

import (
	"os"
	"path/filepath"
)

func UseConfig(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(dir, ".configs", filename+".json")
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
