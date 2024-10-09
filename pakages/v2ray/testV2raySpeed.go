package v2ray

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

// TestV2raySpeed tests the internet speed through a SOCKS5 proxy with a time limit
func TestV2raySpeed(proxyPort int) float64 {
	// fmt.Println("Connect to: SOCKS5", "127.0.0.1:"+strconv.Itoa(proxyPort))

	// Set up the SOCKS5 proxy
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:"+strconv.Itoa(proxyPort), nil, proxy.Direct)
	if err != nil {
		// return 0, fmt.Errorf("failed to create SOCKS5 dialer: %v", err)
		return 0
	}

	// Set up the HTTP client with the SOCKS5 proxy
	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := &http.Client{
		Transport: transport,
	}

	// Create a context with a time limit
	ctx, cancel := context.WithTimeout(context.Background(), 1300*time.Millisecond)
	defer cancel()

	// Create a request with the context
	req, err := http.NewRequestWithContext(ctx, "GET", "https://raw.githubusercontent.com/BitDoctor/speed-test-file/master/5mb.txt", nil)
	if err != nil {
		// return 0, fmt.Errorf("failed to create HTTP request: %v", err)
		return 0
	}

	// Send the HTTP request for the speed test
	resp, err := client.Do(req)
	if err != nil {
		// return 0, fmt.Errorf("failed to send HTTP request: %v", err)
		return 0
	}
	defer resp.Body.Close()

	// Read the data within the 1-second time frame
	start := time.Now()
	var totalBytesRead int64
	buf := make([]byte, 8192) // 8 KB buffer

	for {
		select {
		case <-ctx.Done(): // If the time is up
			duration := time.Since(start)
			// Convert bytes to bits and then to megabits
			speedMbps := float64(totalBytesRead*8) / duration.Seconds() / 1_000_000 // speed in Mbps
			// return speedMbps, nil
			return speedMbps
		default:
			n, err := resp.Body.Read(buf)
			if err != nil && err != io.EOF {
				// return 0, fmt.Errorf("failed to read response body: %v", err)
				return 0
			}
			totalBytesRead += int64(n)
			if err == io.EOF {
				break
			}
		}
	}

	// Calculate speed if the loop ends naturally
	// duration := time.Since(start)
	// speedMbps := float64(totalBytesRead*8) / duration.Seconds() / 1_000_000 // speed in Mbps

	// return speedMbps, nil
	// return speedMbps
}
