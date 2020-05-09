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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/statiko-dev/stkcli/utils"
)

func init() {
	var domain string
	var force bool

	c := &cobra.Command{
		Use:   "status",
		Short: "Shows the status of a node",
		Long: `Prints information about the status and health of the node.

The ` + "`" + `--domain` + "`" + ` flag allows selecting a specific site only.
`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Invoke the /status endpoint to get the status of the node
			// We're not using utils.RequestJSON here because we need to get the status code and parse the response regardless
			url := baseURL + "/status"
			if domain != "" {
				url += "/" + domain
			}
			if force {
				url += "?force=1"
			}
			// Build the request
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Error building request", err)
				return
			}
			// Authorization, if any
			if auth != "" {
				req.Header.Set("Authorization", auth)
			}
			resp, err := client.Do(req)
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
			defer resp.Body.Close()

			// If the status is 2xx or 503 (when the apps are down) parse the /status endpoint response
			if (resp.StatusCode >= 200 && resp.StatusCode <= 299) || resp.StatusCode == 503 {
				// The /status endpoint returns a 503 status code also when there's an issue with the apps, so let's still parse it but show an error
				if resp.StatusCode != http.StatusOK {
					fmt.Printf("\033[31mStatus endpoint returned a %d status code\033[0m\n", resp.StatusCode)
				}

				var r statusResponseModel
				if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
					utils.ExitWithError(utils.ErrorNode, "Invalid JSON response", err)
					return
				}

				fmt.Println(statusResponseModelFormat(&r))
			} else if domain != "" && resp.StatusCode == 404 {
				// While requesting a single domain, the status code was 404, meaning that the domain doesn't exist
				if resp.StatusCode != http.StatusOK {
					fmt.Printf("\033[31mStatus endpoint returned a %d status code\033[0m\n", resp.StatusCode)
					fmt.Println("The requested domain does not exist")
				}
			} else {
				b, _ := ioutil.ReadAll(resp.Body)
				err := fmt.Errorf("invalid response status code: %d; content: %s", resp.StatusCode, string(b))
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
		},
	}
	rootCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&domain, "domain", "d", "", "domain name")
	c.Flags().BoolVarP(&force, "force", "f", false, "force a recheck of all sites, ignoring status cache")

	// Add shared flags
	addSharedFlags(c)
}
