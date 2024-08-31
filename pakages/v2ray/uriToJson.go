package v2ray

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

func UriToJson(uri string, proxyPort int) string {
	if proxyPort == 0 {
		proxyPort = ProxyPort
	}

	if strings.Contains(uri, "vmess://") {
		return vmessToJson(uri, proxyPort)
	}

	if strings.Contains(uri, "vless://") {
		return vlessToJson(uri, proxyPort)
	}

	return ""
}

// Vmess
func vmessToJson(uri string, proxyPort int) string {
	// Step 1: Remove the URI scheme (e.g., "vmess://")
	base64EncodedString := strings.Split(uri, "://")[1]

	// Step 2: Decode the Base64 string
	decodedBytes, err := base64.StdEncoding.DecodeString(base64EncodedString)
	if err != nil {
		return ""
	}
	decodedString := string(decodedBytes)

	// Step 3: Parse the decoded string into a JSON object
	var v2rayDecodedConfig map[string]interface{}
	err = json.Unmarshal([]byte(decodedString), &v2rayDecodedConfig)
	if err != nil {
		return ""
	}

	// Step 4: Create the full V2Ray JSON configuration object
	v2rayJsonConfig := map[string]interface{}{
		"inbounds": []map[string]interface{}{
			{
				"listen":   "0.0.0.0",
				"port":     proxyPort,
				"protocol": "socks",
				"settings": map[string]interface{}{
					"udp": true,
				},
				"sniffing": map[string]interface{}{
					"destOverride": []string{"http", "tls"},
					"enabled":      true,
					"metadataOnly": false,
					"routeOnly":    true,
				},
				"tag": "socks-in",
			},
		},
		"outbounds": []map[string]interface{}{
			{
				"mux": map[string]interface{}{
					"enabled":     false,
					"concurrency": 50,
				},
				"tag":      "proxy",
				"protocol": "vmess",
				"streamSettings": map[string]interface{}{
					"tcpSettings": map[string]interface{}{
						"header": map[string]interface{}{
							"type": "http",
							"request": map[string]interface{}{
								"method":  "GET",
								"version": "1.1",
								"path":    []string{"/"},
								"headers": map[string]interface{}{
									"Accept-Encoding": []string{"gzip, deflate"},
									"Connection":      []string{"keep-alive"},
									"Pragma":          []string{"no-cache"},
									"Host":            []string{"google.com"},
									"User-Agent": []string{
										"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36",
										"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0_2 like Mac OS X) AppleWebKit/601.1 (KHTML, like Gecko) CriOS/53.0.2785.109 Mobile/14A456 Safari/601.1.46",
									},
								},
							},
						},
					},
					"network": "tcp",
				},
				"settings": map[string]interface{}{
					"vnext": []map[string]interface{}{
						{
							"users": []map[string]interface{}{
								{
									"id":       v2rayDecodedConfig["id"],
									"security": v2rayDecodedConfig["scy"],
									"alterId":  0, // v2rayDecodedConfig["aid"]
									"level":    8,
									"email":    "",
								},
							},
							"port":    v2rayDecodedConfig["port"],
							"address": v2rayDecodedConfig["add"],
						},
					},
				},
			},
		},
		"api": map[string]interface{}{
			"tag":      "api",
			"services": []string{"StatsService"},
		},
		"dns": map[string]interface{}{
			"disableFallback":        true,
			"disableFallbackIfMatch": true,
			"queryStrategy":          "UseIP",
			"servers": []map[string]interface{}{
				{
					"address":      "8.8.8.8",
					"skipFallback": false,
				},
			},
			"disableCache": true,
			"tag":          "dnsQuery",
		},
		"stats": map[string]interface{}{},
		"routing": map[string]interface{}{
			"balancers": []interface{}{},
			"rules": []map[string]interface{}{
				{
					"inboundTag":  []string{"api"},
					"type":        "field",
					"outboundTag": "api",
				},
				{
					"type":        "field",
					"inboundTag":  []string{"inDns"},
					"outboundTag": "outDns",
				},
				{
					"type":        "field",
					"inboundTag":  []string{"dnsQuery"},
					"outboundTag": "proxy",
				},
			},
			"domainStrategy": "AsIs",
		},
		"policy": map[string]interface{}{
			"levels": map[string]interface{}{
				"8": map[string]interface{}{
					"uplinkOnly":        1,
					"handshake":         4,
					"downlinkOnly":      1,
					"statsUserDownlink": false,
					"bufferSize":        0,
					"connIdle":          30,
					"statsUserUplink":   false,
				},
			},
			"system": map[string]interface{}{
				"statsOutboundUplink":   true,
				"statsOutboundDownlink": true,
				"statsInboundUplink":    true,
				"statsInboundDownlink":  true,
			},
		},
		"transport": map[string]interface{}{},
	}

	jsonBytes, err := json.MarshalIndent(v2rayJsonConfig, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// Vless
func vlessToJson(uri string, proxyPort int) string {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	queryParams := parsedURL.Query()

	data := map[string]interface{}{
		"dns": map[string]interface{}{
			"hosts": map[string]interface{}{
				"domain:googleapis.cn": "googleapis.com",
			},
			"servers": []string{"1.1.1.1"},
		},
		"inbounds": []map[string]interface{}{
			{
				"listen":   "0.0.0.0",
				"port":     proxyPort,
				"protocol": "socks",
				"settings": map[string]interface{}{
					"udp": true,
				},
				"sniffing": map[string]interface{}{
					"destOverride": []string{"http", "tls"},
					"enabled":      true,
					"metadataOnly": false,
					"routeOnly":    true,
				},
				"tag": "socks-in",
			},
		},
		"log": map[string]interface{}{
			"loglevel": "warning",
		},
		"outbounds": []map[string]interface{}{
			{
				"mux": map[string]interface{}{
					"concurrency":     -1,
					"enabled":         false,
					"xudpConcurrency": 8,
					"xudpProxyUDP443": "",
				},
				"protocol": "vless",
				"settings": map[string]interface{}{
					"vnext": []map[string]interface{}{
						{
							"address": parsedURL.Hostname(),
							"port":    parsePort(parsedURL.Port()),
							"users": []map[string]interface{}{
								{
									"encryption": getQueryParam(queryParams, "encryption", "none"),
									"flow":       "",
									"id":         parsedURL.User.Username(),
									"level":      8,
									"security":   "auto",
								},
							},
						},
					},
				},
				"streamSettings": map[string]interface{}{
					"network":  getQueryParam(queryParams, "type", "tcp"),
					"security": "none",
					"tcpSettings": map[string]interface{}{
						"header": map[string]interface{}{
							"request": map[string]interface{}{
								"headers": map[string]interface{}{
									"Connection":      []string{"keep-alive"},
									"Host":            []string{"speedtest.net"},
									"Pragma":          "no-cache",
									"Accept-Encoding": []string{"gzip, deflate"},
									"User-Agent": []string{
										"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
										"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0_2 like Mac OS X) AppleWebKit/601.1 (KHTML, like Gecko) CriOS/53.0.2785.109 Mobile/14A456 Safari/601.1.46",
									},
								},
								"path":    []string{"/"},
								"method":  "GET",
								"version": "1.1",
							},
							"type": getQueryParam(queryParams, "headerType", "http"),
						},
					},
				},
				"tag": "proxy",
			},
			{
				"protocol": "freedom",
				"settings": map[string]interface{}{},
				"tag":      "direct",
			},
			{
				"protocol": "blackhole",
				"settings": map[string]interface{}{
					"response": map[string]interface{}{
						"type": "http",
					},
				},
				"tag": "block",
			},
		},
		"remarks": "",
		"routing": map[string]interface{}{
			"domainStrategy": "IPIfNonMatch",
			"rules": []map[string]interface{}{
				{
					"ip":          []string{"1.1.1.1"},
					"outboundTag": "proxy",
					"port":        "53",
					"type":        "field",
				},
			},
		},
	}

	jsonBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func getQueryParam(params url.Values, key, defaultValue string) string {
	if value := params.Get(key); value != "" {
		return value
	}
	return defaultValue
}

func parsePort(portStr string) int {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0
	}
	return port
}
