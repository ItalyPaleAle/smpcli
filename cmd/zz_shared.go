/*
Copyright © 2020 Alessandro Segala (@ItalyPaleAle)

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
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/viper"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/spf13/cobra"

	"github.com/statiko-dev/stkcli/utils"
)

var (
	optAddress  string
	optPort     string
	optInsecure bool
	optNoTLS    bool
)

func addSharedFlags(cmd *cobra.Command) {
	// Get deafaults
	defaultNode := viper.GetString("node")
	defaultPort := viper.GetInt("port")
	defaultPortString := ""
	if defaultPort > 0 {
		defaultPortString = strconv.Itoa(defaultPort)
	}
	defaultInsecure := viper.GetBool("insecure")
	defaultHTTP := viper.GetBool("http")

	// Node address
	// Mark as required if we don't have a default value
	if defaultNode == "" {
		cmd.Flags().StringVarP(&optAddress, "node", "n", defaultNode, "node address or IP (required)")
		cmd.MarkFlagRequired("node")
	} else {
		cmd.Flags().StringVarP(&optAddress, "node", "n", defaultNode, "node address or IP")
	}

	// Port the server is listening on
	// Default is 2265
	cmd.Flags().StringVarP(&optPort, "port", "P", defaultPortString, "port the node listens on")

	// Flags to control communication with the node
	// By default, we use TLS and validate the certificate
	cmd.Flags().BoolVarP(&optInsecure, "insecure", "k", defaultInsecure, "disable TLS certificate validation for node connections")
	cmd.Flags().BoolVarP(&optNoTLS, "http", "S", defaultHTTP, "use HTTP protocol, without TLS, for node connections")
}

func getURLClient() (baseURL string, client *http.Client) {
	// Output some warnings
	if optNoTLS {
		fmt.Fprintln(os.Stderr, "\033[33mWARN: You are connecting to your node without using TLS. The connection (including the authorization token) is not encrypted.\033[0m")
	} else if optInsecure {
		fmt.Fprintln(os.Stderr, "\033[33mWARN: TLS certificate validation is disabled. Your connection might not be secure.\033[0m")
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

// This function gets a client authenticated with Azure Key Vault
func getKeyVault() *keyvault.BaseClient {
	// Create a new client
	akvClient := keyvault.New()

	// Authorize from the Azure CLI
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		utils.ExitWithError(utils.ErrorApp, "Error while authorizing the Azure Key Vault client", err)
		return nil
	}
	akvClient.Authorizer = authorizer

	return &akvClient
}

// This function requests the name of the Azure Key Vault from the node
func getKeyVaultInfo() (keyVaultURL string, codesignKeyName string, codesignKeyVersion string, err error) {
	baseURL, client := getURLClient()
	auth := nodeStore.GetAuthToken(optAddress)

	// Invoke the /keyvaultinfo endpoint to get the name and URL of the key vault
	var r map[string]string
	err = utils.RequestJSON(utils.RequestOpts{
		Authorization: auth,
		Client:        client,
		Target:        &r,
		URL:           baseURL + "/keyvaultinfo",
	})
	if err != nil {
		return
	}

	// Check the response
	var ok bool
	keyVaultURL, ok = r["url"]
	if !ok || keyVaultURL == "" {
		err = errors.New("invalid response: empty url")
		return
	}
	codesignKeyName, ok = r["codesignKeyName"]
	if !ok || codesignKeyName == "" {
		err = errors.New("invalid response: empty codesignKeyName")
		return
	}
	codesignKeyVersion, ok = r["codesignKeyVersion"]
	if !ok || codesignKeyVersion == "" {
		err = errors.New("invalid response: empty codesignKeyVersion")
		return
	}

	return
}
