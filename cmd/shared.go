package cmd

import (
	"fmt"
	"net/http"
)

func getURLClient() (baseURL string, client *http.Client) {
	// Output some warnings
	if optNoTLS {
		fmt.Println("\033[33mWARN: You are connecting to your node without using TLS. The connection (including the authorization token) is not encrypted.\033[0m")
	} else if optInsecure {
		fmt.Println("\033[33mWARN: TLS certificate validation is disabled. Your connection might not be secure.\033[0m")
	}

	// Get the URL
	protocol := "https"
	if optNoTLS {
		protocol = "http"
	}

	// Get the URL
	baseURL = fmt.Sprintf("%s://%s:%s", protocol, optAddress, optPort)

	// What client to use?
	client = httpClient
	if optInsecure {
		client = httpClientInsecure
	}

	return
}
