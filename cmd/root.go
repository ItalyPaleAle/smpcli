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
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/ItalyPaleAle/stkcli/utils"
)

var (
	nodeStore *utils.NodeStore

	httpClient         *http.Client
	httpClientInsecure *http.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stkcli",
	Short: "Manage a Statiko node",
	Long: `This CLI is part of the Statiko project: https://statiko.dev

stkcli allows managing Statiko nodes conveniently, by offering a command-line experience that interacts with nodes' REST APIs.
Additionally, stkcli offers commands that simplify uploading and signing app bundles, and uploading TLS certificates.

stkcli is released under a GNU General Public License v3.0 license. Source code is available on GitHub: https://github.com/ItalyPaleAle/stkcli
`,
	DisableAutoGenTag: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(10)
	}
}

func init() {
	cobra.OnInitialize(func() {
		// Init the node store
		nodeStore = &utils.NodeStore{}
		if err := nodeStore.Init(); err != nil {
			panic(err)
		}

		// Initialize the HTTP clients
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}

		// The "insecure" client doesn't validate TLS certificates
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		httpClientInsecure = &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		}
	})
}
