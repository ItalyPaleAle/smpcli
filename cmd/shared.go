package cmd

import (
	"fmt"
	"net/http"
)

func getURLClient() (baseURL string, client *http.Client) {
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
