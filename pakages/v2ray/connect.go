package v2ray

import (
	"fmt"
	"os"
	"path/filepath"
)

// Connect
func Connect(uriIndex int) error {
	// Set core dir
	rootPath, _ := os.Getwd()
	CoreDir = filepath.Join(rootPath, ".v2ray-core")

	// Select
	if selectErr := selectUri(uriIndex); selectErr != nil {
		return selectErr
	}

	// Run
	MainV2RayProcess = NewV2RayProcess(".v2ray", V2rayProxyPort)
	if runErr := MainV2RayProcess.Run(true); runErr != nil {
		return runErr
	}

	return nil
}

// Create uri config file
func selectUri(uriIndex int) error {
	// Define the directory and file path
	dir := ".v2ray"
	filePath := filepath.Join(dir, "config.json")

	// Create the directory if it does not exist
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Open the file for writing (creates the file if it does not exist or truncates it if it does)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(UriToJson(Uris[uriIndex], V2rayProxyPort))
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	ActiveUri = uriIndex
	return nil
}
