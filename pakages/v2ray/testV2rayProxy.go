package v2ray

import (
	"fmt"
	"time"

	"golang.org/x/net/proxy"
)

func TestV2rayProxy(proxyAddr string) bool {
	// Create a new SOCKS5 proxy dialer
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		fmt.Println("Error creating SOCKS5 dialer:", err)
		return false
	}

	// Attempt to connect to a well-known address via the SOCKS5 proxy
	conn, err := dialer.Dial("tcp", "example.com:80")
	if err != nil {
		fmt.Println("Error connecting through SOCKS5 proxy:", err)
		return false
	}
	defer conn.Close()

	// Set a timeout to ensure the connection is working
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	return true
}
