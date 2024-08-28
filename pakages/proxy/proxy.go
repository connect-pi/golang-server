package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"project/pakages/proxy/checkIp"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// SOCKS5 proxy settings
var socks5ProxyAddr = "socks5://0.0.0.0:2080"

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
			// Time
			startTime := time.Now()

			var domain string = ""
			var ip string = ""

			// Get domain
			domain = r.RequestURI
			domain = strings.Replace(domain, "http://", "", -1)
			domain = strings.Replace(domain, "https://", "", -1)
			domain = strings.Replace(domain, "/", "", -1)
			domain = strings.Split(domain, ":")[0]
			fmt.Print("\ndomain: ", domain)

			// Get IP
			if domain == "localhost" {
				ip = "127.0.0.1"
			} else if IsValidDomain(domain) {
				getedIp, ipErr := GetDomainIp(domain)

				if ipErr != nil {
					fmt.Print("\nGet ip error: ", ipErr)
				} else {
					ip = getedIp
				}
			} else { // For ips
				ip = domain
			}

			fmt.Print("\nip: ", ip, "\n")

			isIranIp := checkIp.IsIranIp(ip)
			fmt.Print("isIranIp: ", isIranIp)

			// Time
			elapsedTime := time.Since(startTime)
			fmt.Printf("\n-- Program run time: %s\n", elapsedTime)

			// Send
			if isIranIp {
				// Internet
				if r.Method == http.MethodConnect {
					handleInternetTunnel(w, r)
				} else {
					handleInternetHTTP(w, r)
				}
			} else {
				// Vpn
				if r.Method == http.MethodConnect {
					handleVpnTunnel(w, r)
				} else {
					handleVpnHTTP(w, r)
				}
			}

		}),
	}

	log.Printf("Starting proxy server on %s\n", address)
	if err := proxy.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not start proxy: %v\n", err)
	}
}
