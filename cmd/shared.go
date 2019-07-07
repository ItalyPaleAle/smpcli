/*
Copyright © 2019 Alessandro Segala (@ItalyPaleAle)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	optAddress  string
	optPort     string
	optInsecure bool
	optNoTLS    bool
)

func addSharedFlags(cmd *cobra.Command) {
	// Node address
	cmd.Flags().StringVarP(&optAddress, "node", "n", "", "node address or IP (required)")
	cmd.MarkFlagRequired("node")

	// Port the server is listening on
	// Default is 2265
	// TODO: SET DEFAULT TO 2265 or another better port
	cmd.Flags().StringVarP(&optPort, "port", "p", "2265", "port the node listens on")

	// Flags to control communication with the node
	// By default, we use TLS and validate the certificate
	cmd.Flags().BoolVarP(&optInsecure, "insecure", "k", false, "disable TLS certificate validation")
	cmd.Flags().BoolVarP(&optNoTLS, "http", "s", false, "use HTTP protocol (no TLS)")
}

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
