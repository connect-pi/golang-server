package proxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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

		fmt.Println(" ")

		// Check for open with VPN
		openWithVpn := true
		fmt.Println("vpn: ", openWithVpn)

		var conn net.Conn
		var err error

		if openWithVpn {
			// Forward traffic to the second proxy on port 2080
			if strings.Contains(addr, ":") {
				// Resolve only if addr is a domain name
				// Split addr into host and port
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, fmt.Errorf("failed to split host and port: %w", err)
				}

				// Attempt to resolve the address using Google DNS
				resolver := &net.Resolver{
					PreferGo: true,
					Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
						return net.Dial("udp", "8.8.8.8:53")
					},
				}

				// Resolve the host
				_, err = resolver.LookupHost(ctx, host)
				if err != nil {
					log.Printf("Failed to resolve address %s: %v", host, err)
					return nil, err
				}

				// Reconstruct addr
				addr = fmt.Sprintf("%s:%s", host, port)
			}

			conn, err = dialer.Dial(network, addr)
		} else {
			// Directly connect to the target address for all other traffic
			conn, err = net.Dial(network, addr)
		}

		// End Time
		elapsedTime := time.Since(startTime)
		fmt.Println("👾 Program run time: ", elapsedTime)
		fmt.Println(" ")

		return conn, err
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
