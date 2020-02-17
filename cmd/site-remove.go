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
	"net/http"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/ItalyPaleAle/stkcli/utils"
)

func init() {
	var (
		domain string
		yes    bool
	)

	c := &cobra.Command{
		Use:   "remove",
		Short: "Remove a site",
		Long: `Removes a site from the node, so the web server stops accepting requests for it.

You must specify the primary domain name (no aliases) in the '--domain' parameter to select the site to be removed.
`,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Ask for confirmatiom (unless we have `--yes`)
			if !yes {
				prompt := promptui.Prompt{
					Label:     "Remove the site",
					IsConfirm: true,
				}
				confirm, err := prompt.Run()
				if err != nil || strings.ToLower(confirm) != "y" {
					utils.ExitWithError(utils.ErrorUser, "Aborted", nil)
					return
				}
			}

			// Invoke the /site/:domain endpoint to delete the site
			err := utils.RequestJSON(utils.RequestOpts{
				Authorization: auth,
				Client:        client,
				Method:        utils.RequestDELETE,
				StatusCode:    http.StatusNoContent,
				URL:           baseURL + "/site/" + domain,
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
		},
	}
	siteCmd.AddCommand(c)

	// Flags
	c.Flags().BoolVarP(&yes, "yes", "", false, "do not ask for confirmation")
	c.Flags().StringVarP(&domain, "domain", "d", "", "primary domain name")
	c.MarkFlagRequired("domain")

	// Add shared flags
	addSharedFlags(c)
}
