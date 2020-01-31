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
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	var (
		domain string
		yes    bool
	)

	c := &cobra.Command{
		Use:   "remove",
		Short: "Remove a site",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()

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

			// Ask for confirmatiom (unless we have `--yes`)
			if !yes {
				prompt := promptui.Prompt{
					Label:     "Remove the site",
					IsConfirm: true,
				}
				confirm, err := prompt.Run()
				if err != nil || strings.ToLower(confirm) != "y" {
					fmt.Println("Aborted")
					return
				}
			}

			// Invoke the /site/:domain endpoint to delete the site
			req, err := http.NewRequest("DELETE", baseURL+"/site/"+domain, nil)
			if err != nil {
				fmt.Println("[Fatal error]\nCould not build the request:", err)
				return
			}
			req.Header.Set("Authorization", sharedKey)
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("[Fatal error]\nRequest failed:", err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusNoContent {
				b, _ := ioutil.ReadAll(resp.Body)
				fmt.Printf("[Server error]\n%d: %s\n", resp.StatusCode, string(b))
				return
			}
		},
	}
	siteCmd.AddCommand(c)

	// Flags
	c.Flags().BoolVarP(&yes, "yes", "", false, "do not ask for confirmation")
	c.Flags().StringVarP(&domain, "domain", "d", "", "Primary domain name")
	c.MarkFlagRequired("domain")

	// Add shared flags
	addSharedFlags(c)
}
