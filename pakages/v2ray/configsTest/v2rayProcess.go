package configsTest

import (
	"fmt"
	"path/filepath"
	"project/pakages/v2ray"
	"strconv"
	"sync"
)

func RunTestV2RayProcesses() int {
	fmt.Println("\n-------------")
	fmt.Println("Run test V2Ray processes...")

	var testResult []float64
	var mu sync.Mutex
	var wg sync.WaitGroup

	// تعداد پردازش‌ها برای انتظار
	wg.Add(len(v2ray.Uris))

	for i := range v2ray.Uris {
		go func(i int) {
			defer wg.Done()
			// ایجاد مسیر برای هر config
			path := filepath.Join(".v2ray", "testConfigs", strconv.Itoa(i))
			port := 3281 + i

			// ایجاد یک نمونه جدید از V2RayProcess با مسیر مورد نظر
			v2Process := v2ray.NewV2RayProcess(path, port)

			// اجرای پراسس V2Ray
			if err := v2Process.Run(); err != nil {
				fmt.Printf("Failed to start V2Ray for config %d: %v\n", i, err)
				mu.Lock()
				testResult = append(testResult, 0)
				mu.Unlock()
				return
			}

			// تست سرعت اینترنت
			bytesRead, err := v2ray.TestV2raySpeed(port)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				mu.Lock()
				testResult = append(testResult, 0)
				mu.Unlock()
			} else {
				mu.Lock()
				testResult = append(testResult, float64(bytesRead))
				mu.Unlock()
			}

			// متوقف کردن پراسس V2Ray
			if err := v2Process.Stop(); err != nil {
				fmt.Printf("Failed to stop V2Ray for config %d: %v\n", i, err)
			}
		}(i)
	}

	// منتظر ماندن تا همه goroutine ها به پایان برسند
	wg.Wait()

	maxSpeed := testResult[0]
	maxIndex := 0
	for i, speed := range testResult {
		if speed > maxSpeed {
			maxSpeed = speed
			maxIndex = i
		}
	}

	fmt.Printf("\nFastest speed: %.2f bytes\n", maxSpeed)
	fmt.Printf("Index of fastest speed: %d\n", maxIndex)
	fmt.Printf("\n-------------\n")

	return maxIndex
}
