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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/ItalyPaleAle/stkcli/utils"
)

func init() {
	var domain string

	c := &cobra.Command{
		Use:   "status",
		Short: "Shows the status of a node",
		Long: `Prints information about the status and health of the node.

The '--domain' flag allows selecting a specific site only.
`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()

			// Invoke the /status endpoint to get the status of the node
			// We're not using utils.RequestJSON here because we need to get the status code and parse the response regardless
			url := baseURL + "/status"
			if domain != "" {
				url += "/" + domain
			}
			resp, err := client.Get(url)
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
			defer resp.Body.Close()

			// The /status endpoint returns a non-200 status code also when there's an issue with the apps, so let's still parse it but show an error
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("\033[31mStatus endpoint returned a %d status code\033[0m\n", resp.StatusCode)
			}

			// Parse the response
			var r statusResponseModel
			if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
				utils.ExitWithError(utils.ErrorNode, "Invalid JSON response", err)
				return
			}

			fmt.Println(statusResponseModelFormat(&r))
		},
	}
	rootCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&domain, "domain", "d", "", "domain name")

	// Add shared flags
	addSharedFlags(c)
}
