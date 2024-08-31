package v2ray

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"
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
	return &V2RayProcess{
		Cmd:    exec.Command("v2ray", "run"),
		StopCh: make(chan struct{}),
		Path:   path,
		Port:   port,
	}
}

// Run V2Ray
func (vp *V2RayProcess) Run() error {
	fmt.Println("Start V2Ray...")

	// Set the working directory to the specified path
	vp.Cmd.Dir = vp.Path

	// Start the command
	if err := vp.Cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	time.Sleep(1 * time.Second)

	if TestV2rayProxy(vp.Port) {
		fmt.Println("✅ V2Ray started successfully.")
		vp.IsRun = true
	} else {
		vp.Stop()
		return fmt.Errorf("V2Ray did not start within the expected time")
	}

	return nil
}

// Stop V2Ray
func (vp *V2RayProcess) Stop() error {
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
	fmt.Println("💤 V2Ray stopped.")
	return nil
}
