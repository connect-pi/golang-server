package proxy

import (
	"io"
	"log"
	"net"
	"net/url"
)

func handleConnection(clientConn net.Conn, proxyURL *url.URL) {
	defer clientConn.Close()

	// Connect to the HTTP proxy
	proxyConn, err := net.Dial("tcp", proxyURL.Host)
	if err != nil {
		log.Printf("Error connecting to proxy: %v", err)
		return
	}
	defer proxyConn.Close()

	// Relay data between client and proxy
	go io.Copy(proxyConn, clientConn)
	io.Copy(clientConn, proxyConn)
}

func startTCPServer(listenAddr string, proxyURL *url.URL) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("Listening on %s", listenAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(clientConn, proxyURL)
	}
}

func StartTcpProxy() {
	// Configure the proxy URL
	proxyURL, err := url.Parse("http://localhost:1080") // Change this to your proxy URL
	if err != nil {
		log.Fatalf("Error parsing proxy URL: %v", err)
	}

	// Start the TCP server
	log.Printf("Start the TCP...")

	startTCPServer(":9999", proxyURL) // Change the listening address and port as needed
}
