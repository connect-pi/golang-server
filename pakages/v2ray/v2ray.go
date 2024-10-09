package v2ray

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"
)

// V2RayProcess holds information about a V2Ray process
type V2RayProcess struct {
	Cmd    *exec.Cmd
	IsRun  bool
	StopCh chan struct{}
	Path   string
	Port   int
}

// NewV2RayProcess creates a new V2Ray process with a given path
func NewV2RayProcess(path string, port int) *V2RayProcess {
	command := filepath.Join(CoreDir, "/v2ray")

	// fmt.Println(exec.Command(command, "run"))

	return &V2RayProcess{
		Cmd:    exec.Command(command, "run"),
		IsRun:  false,
		StopCh: make(chan struct{}),
		Path:   path,
		Port:   port,
	}
}

// Run V2Ray
func (vp *V2RayProcess) Run(prints bool) error {
	if prints {
		fmt.Println("Start V2Ray...")

	}

	// Set the working directory to the specified path
	vp.Cmd.Dir = vp.Path

	// Start the command
	if err := vp.Cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	// Set IsRun to true after starting the command
	vp.IsRun = true

	if prints {
		fmt.Println("✅ V2Ray started successfully.")
	}

	return nil
}

// Stop V2Ray
func (vp *V2RayProcess) Stop(prints bool) error {
	fmt.Println("Stop V2RayProcess")
	if vp.Cmd == nil || vp.Cmd.Process == nil {
		return fmt.Errorf("no process to stop")
	}

	// Send a signal to stop the process
	if err := vp.Cmd.Process.Signal(syscall.SIGINT); err != nil {
		return fmt.Errorf("error sending stop signal: %v", err)
	}

	// Wait for the process to finish
	done := make(chan error)
	go func() {
		done <- vp.Cmd.Wait()
	}()

	select {
	case <-vp.StopCh:
		// Stop signal received
	case err := <-done:
		if err != nil {
			return fmt.Errorf("process stopped with error: %v", err)
		}
	}

	vp.IsRun = false
	if prints {
		fmt.Println("💤 V2Ray stopped.")
	}
	return nil
}
