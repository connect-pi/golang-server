package proxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"project/pakages/v2ray"
	"time"

	"github.com/things-go/go-socks5"
	"golang.org/x/net/proxy"
)

var IsRun bool = false

func StartAppSocks5Proxy() {
	if IsRun {
		return
	}

	// Define the address of the downstream SOCKS5 proxy
	upstreamProxyAddress := "0.0.0.0:2088"

	// Create a SOCKS5 dialer to connect to the upstream proxy
	dialer, err := proxy.SOCKS5("tcp", upstreamProxyAddress, nil, proxy.Direct)
	if err != nil {
		log.Fatalf("Failed to create SOCKS5 dialer: %v", err)
	}

	// Custom dial function to selectively forward traffic
	dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		// Start Time
		startTime := time.Now()

		// log.Println("Address: ", addr)
		fmt.Println(" ")

		// Check for open with VPN
		v2rayIsRun := v2ray.MainV2RayProcess != nil && v2ray.MainV2RayProcess.IsRun
		openWithVpn := v2rayIsRun && OpenWithVpnOrNot(addr)
		fmt.Println("vpn: ", openWithVpn)

		if openWithVpn {
			// Forward traffic to the second proxy on port 2080
			// log.Printf("Forwarding traffic to %s via upstream proxy", addr)

			// End Time
			elapsedTime := time.Since(startTime)
			fmt.Println("👾 Program run time: ", elapsedTime)
			fmt.Println(" ")

			return dialer.Dial(network, addr)
		}

		// Directly connect to the target address for all other traffic
		// log.Printf("Connecting directly to %s", addr)

		// End Time
		elapsedTime := time.Since(startTime)
		fmt.Println("👾 Program run time: ", elapsedTime)
		fmt.Println(" ")

		return net.Dial(network, addr)
	}

	// Create a SOCKS5 server with a custom dial function
	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
		socks5.WithDial(dial),
	)

	// Start the SOCKS5 proxy on localhost port 8000
	go func() {
		if err := server.ListenAndServe("tcp", ":8000"); err != nil {
			log.Fatalf("Failed to start SOCKS5 proxy: %v", err)
		}
	}()

	// Log message after starting the proxy
	fmt.Println("✨ Starting SOCKS5 proxy server on :8000")

	IsRun = true

	// Keep the main function running
	select {}
}
