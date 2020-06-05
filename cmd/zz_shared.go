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
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	optAddress  string
	optPort     string
	optInsecure bool
	optHTTP     bool
)

func addSharedFlags(cmd *cobra.Command) {
	// Get deafaults
	defaultNode := viper.GetString("node")
	defaultPort := viper.GetInt("port")
	defaultPortString := ""
	if defaultPort > 0 {
		defaultPortString = strconv.Itoa(defaultPort)
	}

	// Insecure and HTTP are read from the config only
	optInsecure = viper.GetBool("insecure")
	optHTTP = viper.GetBool("http")

	// Node address
	cmd.Flags().StringVarP(&optAddress, "node", "N", defaultNode, "node address or IP")

	// Port the server is listening on
	// Default is 2265
	cmd.Flags().StringVarP(&optPort, "port", "P", defaultPortString, "port the node listens on")
}

func getURLClient() (baseURL string, client *http.Client) {
	// Output some warnings
	if optHTTP {
		fmt.Fprintln(os.Stderr, "\033[33mWARN: You are connecting to your node without using TLS. The connection (including the authorization token) is not encrypted.\033[0m")
	} else if optInsecure {
		fmt.Fprintln(os.Stderr, "\033[33mWARN: TLS certificate validation is disabled. Your connection might not be secure.\033[0m")
	}

	// Get the URL
	protocol := "https"
	if optHTTP {
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

// Accepts a PEM-encoded key or the path to a key
func loadRSAPrivateKey(key string) *rsa.PrivateKey {
	// Check if we have a key, then parse it
	if key == "" {
		return nil
	}

	// Check if we have a key or the path to a file
	if !strings.HasPrefix(key, "-----BEGIN") {
		read, err := ioutil.ReadFile(key)
		if err == nil && len(read) > 0 {
			key = string(read)
		} else {
			return nil
		}
	}

	// Parse the pem file
	block, _ := pem.Decode([]byte(key))
	if block == nil || len(block.Bytes) == 0 {
		return nil
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		// PKCS#1
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil || key == nil {
			return nil
		}
		return key
	case "PRIVATE KEY":
		// PKCS#8 (un-encrypted)
		pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil || pk == nil {
			return nil
		}
		key, ok := pk.(*rsa.PrivateKey)
		if !ok {
			return nil
		}
		return key
	default:
		return nil
	}
}
