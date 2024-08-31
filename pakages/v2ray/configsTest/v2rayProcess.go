package configsTest

import (
	"fmt"
	"path/filepath"
	"project/pakages/v2ray"
	"strconv"
	"sync"
)

func RunTestV2RayProcesses() int {
	fmt.Println("\n-------------\n")
	fmt.Println("Run test V2Ray processes...\n")

	// Initialize a slice with a length equal to the number of URIs
	testResult := make([]float64, len(v2ray.Uris))
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Number of processes to wait for
	wg.Add(len(v2ray.Uris))

	for i := range v2ray.Uris {
		go func(i int) {
			defer wg.Done()
			// Create path for each config
			path := filepath.Join(".v2ray", "testConfigs", strconv.Itoa(i))
			port := 3281 + i

			// Create a new V2RayProcess instance with the specified path
			v2Process := v2ray.NewV2RayProcess(path, port)

			// Start the V2Ray process
			if err := v2Process.Run(false); err != nil {
				// fmt.Printf("Failed to start V2Ray for config %d: %v\n", i, err)
				mu.Lock()
				testResult[i] = 0
				fmt.Printf("\nConfig %d: -Mb/s", i)
				mu.Unlock()
				return
			}

			// Test internet speed
			mu.Lock()
			speed := v2ray.TestV2raySpeed(port)
			testResult[i] = speed
			fmt.Printf("Config %d: %.2fMb/s\n", i, testResult[i])
			mu.Unlock()

			// Stop the V2Ray process
			v2Process.Stop(false)
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Determine the maximum speed and its index
	maxSpeed := testResult[0]
	maxIndex := 0
	for i, speed := range testResult {
		if speed > maxSpeed {
			maxSpeed = speed
			maxIndex = i
		}
	}

	fmt.Printf("\n💥 Fastest speed: %.2f Mb/s\n", maxSpeed)
	fmt.Printf("\n-------------\n\n")

	return maxIndex
}
