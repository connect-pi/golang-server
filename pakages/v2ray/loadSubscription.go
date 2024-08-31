package v2ray

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"project/pakages/configs"
	"strings"
)

// Load
func LoadSubscription() error {
	fmt.Println("Load subs...")

	// Define the URL from which to fetch the response
	url := configs.Settings.SubscriptionLink

	// Make an HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("ERROR: Failed to fetch data from %s: %v\n", url, err)
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: Failed to read response body: %v\n", err)
		return err
	}

	// Process the response (assuming each character represents a subscription)
	response := string(body)
	decodedString, _ := base64.StdEncoding.DecodeString(response)
	uris := strings.Split(string(decodedString), "\n")
	finalUris := []string{}

	for _, uri := range uris {
		if uri != "False" {
			finalUris = append(finalUris, uri)
		}
	}

	fmt.Println(len(finalUris), "uris were loaded!")

	// Set uris
	Uris = finalUris

	return nil
}
