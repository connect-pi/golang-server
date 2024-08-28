package v2ray

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	// Create a pipe to capture the command's output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	commands := ""
	resultChan := make(chan bool)
	errorChan := make(chan error)

	// Read the command's output in a separate goroutine
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			commands += line

			// Check if the output contains the word "started"
			if strings.Contains(line, "started") {
				fmt.Println("✅ V2Ray started successfully.")
				IsRun = true
				resultChan <- true
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errorChan <- fmt.Errorf("error reading output: %v", err)
			return
		}

		// If we reach here, it means "started" was never found
		resultChan <- false
	}()

	select {
	case success := <-resultChan:
		if success {
			// V2Ray started successfully, return without waiting for the goroutine to finish
			return nil
		}
	case err := <-errorChan:
		// An error occurred while reading the output
		StopV2Ray()
		return err
	case <-time.After(1 * time.Second):
		// Timeout after 1 second, stop the goroutine and return an error
		StopV2Ray()
		return fmt.Errorf("V2Ray did not start within the expected time")
	}

	// If we reach here, return an error indicating failure to start
	return fmt.Errorf("V2Ray did not start")
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
