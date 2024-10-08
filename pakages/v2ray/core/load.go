package v2raycore

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var downloadLinks = map[string]string{
	"android-arm64-v8a.zip": "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/android-arm64-v8a.zip",
	"dragonfly-64.zip":      "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/dragonfly-64.zip",
	"freebsd-32.zip":        "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/freebsd-32.zip",
	"freebsd-64.zip":        "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/freebsd-64.zip",
	"freebsd-arm32-v6.zip":  "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/freebsd-arm32-v6.zip",
	"freebsd-arm32-v7a.zip": "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/freebsd-arm32-v7a.zip",
	"freebsd-arm64-v8a.zip": "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/freebsd-arm64-v8a.zip",
	"linux-32.zip":          "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-32.zip",
	"linux-64.zip":          "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-64.zip",
	"linux-arm32-v5.zip":    "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-arm32-v5.zip",
	"linux-arm32-v6.zip":    "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-arm32-v6.zip",
	"linux-arm32-v7a.zip":   "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-arm32-v7a.zip",
	"linux-arm64-v8a.zip":   "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-arm64-v8a.zip",
	"linux-loong64.zip":     "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-loong64.zip",
	"linux-mips32.zip":      "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-mips32.zip",
	"linux-mips32le.zip":    "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-mips32le.zip",
	"linux-mips64.zip":      "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-mips64.zip",
	"linux-mips64le.zip":    "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-mips64le.zip",
	"linux-riscv64.zip":     "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/linux-riscv64.zip",
	"macos-64.zip":          "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/macos-64.zip",
	"macos-arm64-v8a.zip":   "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/macos-arm64-v8a.zip",
	"openbsd-32.zip":        "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/openbsd-32.zip",
	"openbsd-64.zip":        "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/openbsd-64.zip",
	"openbsd-arm32-v6.zip":  "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/openbsd-arm32-v6.zip",
	"openbsd-arm32-v7a.zip": "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/openbsd-arm32-v7a.zip",
	"openbsd-arm64-v8a.zip": "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/openbsd-arm64-v8a.zip",
	"windows-32.zip":        "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/windows-32.zip",
	"windows-64.zip":        "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/windows-64.zip",
	"windows-arm32-v7a.zip": "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/windows-arm32-v7a.zip",
	"windows-arm64-v8a.zip": "https://hardl.s3.ir-thr-at1.arvanstorage.ir/VRY/v5.19.0/windows-arm64-v8a.zip",
}

func Load() {
	// Identify OS and architecture
	osType := runtime.GOOS
	archType := runtime.GOARCH

	var fileName string
	switch osType {
	case "linux":
		if archType == "amd64" {
			fileName = "linux-64.zip"
		} else if archType == "386" {
			fileName = "linux-32.zip"
		} else if archType == "arm" {
			fileName = "linux-arm32-v6.zip" // or choose another type
		} else if archType == "arm64" {
			fileName = "linux-arm64-v8a.zip"
		}
	case "windows":
		if archType == "amd64" {
			fileName = "windows-64.zip"
		} else if archType == "386" {
			fileName = "windows-32.zip"
		} else if archType == "arm" {
			fileName = "windows-arm32-v7a.zip"
		} else if archType == "arm64" {
			fileName = "windows-arm64-v8a.zip"
		}
	case "darwin":
		if archType == "amd64" {
			fileName = "macos-64.zip"
		} else if archType == "arm64" {
			fileName = "macos-arm64-v8a.zip"
		}
	default:
		fmt.Printf("Unsupported OS: %s\n", osType)
		return
	}

	// Check if the v2ray-core directory exists
	coreDir := "v2ray-core"
	if _, err := os.Stat(coreDir); os.IsNotExist(err) {
		fmt.Printf("\nStarting download V2ray core process...\n")
		err := os.MkdirAll(coreDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}
	} else {
		return
	}

	// Download the file
	url, ok := downloadLinks[fileName]
	if !ok {
		fmt.Printf("No suitable version found for %s/%s.\n", osType, archType)
		return
	}

	err := downloadFile(fileName, url)
	if err != nil {
		fmt.Printf("Download error: %v\n", err)
		return
	}

	fmt.Printf("Download completed successfully: %s\n\n", fileName)

	// Extract and delete the downloaded file
	if err := extractAndDelete(fileName); err != nil {
		fmt.Printf("Error extracting and deleting file: %v\n", err)
		return
	}

	// fmt.Printf("File extracted and downloaded file deleted: %s\n", fileName)
}

func downloadFile(fileName, url string) error {
	// Create v2ray-core directory if it doesn't exist
	err := os.MkdirAll("v2ray-core", os.ModePerm)
	if err != nil {
		return err
	}

	// Update the file name to store in the v2ray-core directory
	fileName = filepath.Join("v2ray-core", fileName)

	// Create the file for storing
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the contents of the URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the server response is successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: %s", resp.Status)
	}

	// Get the total size of the file
	totalSize := resp.ContentLength
	fmt.Printf("File: %s (size: %.1f MB)\n", fileName, float64(totalSize)/1024/1024)

	// Create a buffer to hold the downloaded data
	buffer := make([]byte, 1024)
	var downloaded int64
	startTime := time.Now() // Start time of download

	// Create a ticker to update every second
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	// Download the file in chunks
	go func() {
		for range ticker.C {
			elapsedTime := time.Since(startTime).Seconds()
			speed := float64(downloaded) / (1024 * 1024) / elapsedTime // Download speed in MB/s
			remainingTime := (float64(totalSize) - float64(downloaded)) / (1024 * 1024 * speed)

			progress := float64(downloaded) / float64(totalSize) * 100
			downloadedMB := float64(downloaded) / 1024 / 1024
			totalMB := float64(totalSize) / 1024 / 1024

			// Print progress bar and statistics
			fmt.Printf("\rDownloading... [")
			percent := int(progress / 2) // Scale progress to bar length
			for i := 0; i < 50; i++ {
				if i < percent {
					fmt.Print("=")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Printf("] %.1f%% (%.1f/%.1f MB) - %.0fs ", progress, downloadedMB, totalMB, remainingTime)
		}
	}()

	// Download in chunks
	for {
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break // Finished downloading
		}

		// Write the chunk to the file
		_, err = out.Write(buffer[:n])
		if err != nil {
			return err
		}

		// Update the number of downloaded bytes
		downloaded += int64(n)
	}

	fmt.Println() // New line after download completes
	return nil
}

func extractAndDelete(fileName string) error {
	zipFile := filepath.Join("v2ray-core", fileName)

	// Open the zip file
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create each file in the zip archive
	for _, f := range r.File {
		fpath := filepath.Join("v2ray-core", f.Name)

		if f.FileInfo().IsDir() {
			// Create directories
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		// Create the necessary directories
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Open the zip file for reading
		rc, err := f.Open()
		if err != nil {
			return err
		}

		// Create a new file for writing
		dstFile, err := os.Create(fpath)
		if err != nil {
			return err
		}

		// Copy the content to the new file
		_, err = io.Copy(dstFile, rc)
		if err != nil {
			return err
		}

		dstFile.Close()
		rc.Close()
	}

	// Delete the downloaded zip file
	err = os.Remove(zipFile)
	if err != nil {
		return err
	}

	// Set executable permissions on the v2ray file
	v2rayPath := filepath.Join("v2ray-core", "v2ray") // Make sure the filename matches
	err = os.Chmod(v2rayPath, 0755)                   // Set permissions to rwxr-xr-x
	if err != nil {
		return fmt.Errorf("failed to set permissions: %v", err)
	}

	return nil
}
