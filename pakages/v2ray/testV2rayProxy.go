package v2ray

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

func TestV2rayProxy(proxyPort int) bool {
	// Create a new SOCKS5 proxy dialer
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:"+strconv.Itoa(proxyPort), nil, proxy.Direct)
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
