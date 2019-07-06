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
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"smpcli/utils"
)

var (
	nodeStore *utils.NodeStore

	httpClient         *http.Client
	httpClientInsecure *http.Client

	optAddress  string
	optPort     string
	optInsecure bool
	optNoTLS    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "smpcli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		// Load config
		if err := loadConfig(); err != nil {
			panic(err)
		}

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

	// Node address
	rootCmd.PersistentFlags().StringVarP(&optAddress, "node", "n", "", "node address or IP (required)")
	rootCmd.MarkPersistentFlagRequired("node")

	// Port the server is listening on
	// Default is 2265
	// TODO: SET DEFAULT TO 2265 or another better port
	rootCmd.PersistentFlags().StringVarP(&optPort, "port", "p", "2265", "port the node listens on")

	// Flags to control communication with the node
	// By default, we use TLS and validate the certificate
	rootCmd.PersistentFlags().BoolVarP(&optInsecure, "insecure", "k", false, "disable TLS certificate validation")
	rootCmd.PersistentFlags().BoolVarP(&optNoTLS, "http", "s", false, "use HTTP protocol (no TLS)")
}

func loadConfig() error {
	// Get the home directory
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	// Ensure the config folder exists
	folder := filepath.FromSlash(home + "/.smpcli")
	if err := utils.EnsureFolder(folder); err != nil {
		return err
	}

	// Load the config file
	viper.SetConfigFile(folder + "/config.yaml")

	// Read in the config file, ignoring errors
	_ = viper.ReadInConfig()

	return nil
}
