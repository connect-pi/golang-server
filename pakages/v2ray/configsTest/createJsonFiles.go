package configsTest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"project/pakages/v2ray"
)

// Directory path where the files will be placed
var DirPath = ".v2ray/testConfigs"

// Create JSON files based on the provided data
func CreateJsonFiles() {
	// Remove the existing test directory if it exists
	RemoveTestsDir()

	// Array of strings as input
	data := v2ray.Uris

	// Recreate the directory
	if err := os.MkdirAll(DirPath, os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// For each item in the array, create a folder and a JSON file
	for i, item := range data {
		// Create a subdirectory with the index as the name
		subDirPath := filepath.Join(DirPath, fmt.Sprintf("%d", i))
		if err := os.MkdirAll(subDirPath, os.ModePerm); err != nil {
			fmt.Println("Error creating subdirectory:", err)
			continue
		}

		// JSON file name
		fileName := filepath.Join(subDirPath, "config.json")

		// Generate JSON content
		jsonData := v2ray.UriToJson(item, 3281+i)

		// Open or create the file
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating file:", err)
			continue
		}

		// Write JSON content to the file
		if _, err := io.WriteString(file, jsonData); err != nil {
			fmt.Println("Error writing to file:", err)
			file.Close()
			continue
		}

		// Close the file after writing
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
			continue
		}
	}

	fmt.Println("Test JSON files created!")
}

// Remove the testing directory if it exists
func RemoveTestsDir() {
	if err := os.RemoveAll(DirPath); err != nil {
		fmt.Println("Error removing directory:", err)
		return
	}
}
