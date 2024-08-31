package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"project/pakages/v2ray"
	"time"

	"golang.org/x/net/proxy"
)

// SOCKS5 proxy settings
var socks5ProxyAddr = "socks5://0.0.0.0:2086"

// getSocks5Dialer returns a dialer that connects through the SOCKS5 proxy
func getSocks5Dialer() (proxy.Dialer, error) {
	u, err := url.Parse(socks5ProxyAddr)
	if err != nil {
		return nil, err
	}
	return proxy.FromURL(u, proxy.Direct)
}

// HandleTunnel handles HTTPS connections by tunneling them through the proxy.
func handleVpnTunnel(w http.ResponseWriter, r *http.Request) {
	dialer, err := getSocks5Dialer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	destConn, err := dialer.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	go io.Copy(destConn, clientConn)
	io.Copy(clientConn, destConn)
}

// HandleHTTP handles regular HTTP requests by forwarding them through the proxy.
func handleVpnHTTP(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(r.RequestURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(r.Method, url.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	dialer, err := getSocks5Dialer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}
	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// HandleTunnel handles HTTPS connections by tunneling them through the proxy.
func handleInternetTunnel(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	go io.Copy(destConn, clientConn)
	io.Copy(clientConn, destConn)
}

// HandleHTTP handles regular HTTP requests by forwarding them through the proxy.
func handleInternetHTTP(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(r.RequestURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(r.Method, url.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// StartProxyServer starts the proxy server on a specified address.
func Start(address string) {

	proxy := &http.Server{
		Addr: address,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			fmt.Print("\nV2ray Status: ", v2ray.IsRun)

			// Start Time
			startTime := time.Now()

			// Check for open with vpn
			openWithVpn := OpenWithVpnOrNot(r.RequestURI)
			fmt.Print("\nvpn: ", openWithVpn)

			// End Time
			elapsedTime := time.Since(startTime)
			fmt.Printf("\n-- Program run time: %s\n", elapsedTime)

			// Send
			if openWithVpn {
				// Vpn
				if r.Method == http.MethodConnect {
					handleVpnTunnel(w, r)
				} else {
					handleVpnHTTP(w, r)
				}
			} else {
				// Internet
				if r.Method == http.MethodConnect {
					handleInternetTunnel(w, r)
				} else {
					handleInternetHTTP(w, r)
				}
			}

		}),
	}

	fmt.Println("✅ Starting proxy server on ", address)
	if err := proxy.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("❌ Could not start proxy: ", err)
	}
}
