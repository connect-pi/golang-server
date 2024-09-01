package v2ray

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

// TestV2rayPing tests the ping time to Google through a SOCKS5 proxy with a time limit of 1.5 seconds
func TestV2rayPing(proxyPort int) float64 {
	// Set up the SOCKS5 proxy
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:"+strconv.Itoa(proxyPort), nil, proxy.Direct)
	if err != nil {
		return 0
	}

	// Set up the HTTP client with the SOCKS5 proxy
	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := &http.Client{
		Transport: transport,
	}

	// Create a context with a time limit of 1 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	// Create a request with the context
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.google.com", nil)
	if err != nil {
		return 0
	}

	// Record the start time
	start := time.Now()

	// Send the HTTP request for the ping test
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	// Calculate the ping time
	duration := time.Since(start).Seconds() * 1000 // Convert to milliseconds

	// Return the ping time in milliseconds
	return duration
}
