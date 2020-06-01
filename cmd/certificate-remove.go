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
	"net/http"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/statiko-dev/stkcli/utils"
)

func init() {
	var (
		name string
		yes  bool
	)

	c := &cobra.Command{
		Use:               "remove",
		Short:             "Remove a TLS certificate",
		Long:              `Removes an imported TLS certificate that is stored in the node's state.`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Ask for confirmation (unless we have `--yes`)
			if !yes {
				prompt := promptui.Prompt{
					Label:     "Remove the certificate",
					IsConfirm: true,
				}
				confirm, err := prompt.Run()
				if err != nil || strings.ToLower(confirm) != "y" {
					utils.ExitWithError(utils.ErrorUser, "Aborted", nil)
					return
				}
			}

			// Invoke the /certificate/:name endpoint to delete the certificate
			err := utils.RequestJSON(utils.RequestOpts{
				Authorization: auth,
				Client:        client,
				Method:        utils.RequestDELETE,
				StatusCode:    http.StatusNoContent,
				URL:           baseURL + "/certificate/" + name,
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
		},
	}
	certificateCmd.AddCommand(c)

	// Flags
	c.Flags().BoolVarP(&yes, "yes", "", false, "do not ask for confirmation")
	c.Flags().StringVarP(&name, "name", "i", "", "name for the certificate")
	c.MarkFlagRequired("name")

	// Add shared flags
	addSharedFlags(c)
}
