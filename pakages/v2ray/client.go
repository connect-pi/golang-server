package v2ray

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

var cmd *exec.Cmd
var stopChan = make(chan struct{}) // Channel to signal stopping

func Connect(uriIndex int) error {
	if IsRun {
		// Stop
		if stopErr := StopV2Ray(); stopErr != nil {
			return stopErr
		}
	}

	// Select
	if selectErr := selectUri(uriIndex); selectErr != nil {
		return selectErr
	}

	// Run
	if runErr := runV2Ray(); runErr != nil {
		return runErr
	}

	return nil
}

// Create uri config file
func selectUri(uriIndex int) error {
	// Define the directory and file path
	dir := ".v2rayConfig"
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
	_, err = file.WriteString(UriToJson(Uris[uriIndex]))
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	ActiveUri = uriIndex
	return nil
}

func runV2Ray() error {
	fmt.Println("Start V2Ray...")

	// Create the command
	cmd := exec.Command("v2ray", "run")
	cmd.Dir = ".v2rayConfig"

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	time.Sleep(1 * time.Second)

	if TestSocks5Proxy("127.0.0.1:2085") {
		fmt.Println("✅ V2Ray started successfully.")
		IsRun = true
	} else {
		StopV2Ray()
		return fmt.Errorf("V2Ray did not start within the expected time")
	}

	return nil
}

// Stop v2ray
func StopV2Ray() error {
	if cmd == nil || cmd.Process == nil {
		return fmt.Errorf("no process to stop")
	}

	// Send a signal to stop the process
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		return fmt.Errorf("error sending stop signal: %v", err)
	}

	// Wait for the process to finish
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-stopChan:
		// Stop signal received
	case err := <-done:
		if err != nil {
			return fmt.Errorf("process stopped with error: %v", err)
		}
	}

	IsRun = false
	fmt.Println("💤 V2Ray stopped.")
	return nil
}
