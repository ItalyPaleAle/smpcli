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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	var (
		domain  string
		app     string
		version string
	)

	// This funtion sends the request to the node to deploy the app
	// Returns true in case of success, and false if there's an error
	var sendRequest = func(sharedKey string) bool {
		baseURL, client := getURLClient()

		// Request body
		reqBody := &deployRequestModel{
			App:     app,
			Version: version,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(reqBody)

		// Invoke the /site endpoint and add the site
		req, err := http.NewRequest("POST", baseURL+"/site/"+domain+"/deploy", buf)
		if err != nil {
			fmt.Println("[Fatal error]\nCould not build the request:", err)
			return false
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", sharedKey)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("[Fatal error]\nRequest failed:", err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			b, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("[Server error]\n%d: %s\n", resp.StatusCode, string(b))
			return false
		}

		// Parse the response
		var r deployResponseModel
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			fmt.Println("[Fatal error]\nInvalid JSON response:", err)
			return false
		}

		fmt.Println(deployResponseModelFormat(&r))
		return true
	}

	c := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an app",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			// Get the shared key
			sharedKey, found, err := nodeStore.GetSharedKey(optAddress)
			if err != nil {
				fmt.Println("[Fatal error]\nError while reading node store:", err)
				return
			}
			if !found {
				fmt.Printf("[Error]\nNo authentication data for the domain %s; please make sure you've executed the 'auth' command.\n", optAddress)
				return
			}

			// Returns true if succeeded
			_ = sendRequest(sharedKey)
		},
	}
	rootCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&domain, "domain", "d", "", "Primary domain name (required)")
	c.MarkFlagRequired("domain")
	c.Flags().StringVarP(&app, "app", "a", "", "App's bundle name (required)")
	c.MarkFlagRequired("app")
	c.Flags().StringVarP(&version, "version", "v", "", "App's bundle version (required)")
	c.MarkFlagRequired("version")

	// Add shared flags
	addSharedFlags(c)
}
