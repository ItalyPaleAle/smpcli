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
		app string
		yes bool
	)

	c := &cobra.Command{
		Use:               "remove",
		Short:             "Remove an app from the node's repository",
		Long:              `Removes an app that is currently stored in the node's repository.`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Ask for confirmation (unless we have `--yes`)
			if !yes {
				prompt := promptui.Prompt{
					Label:     "Remove the app",
					IsConfirm: true,
				}
				confirm, err := prompt.Run()
				if err != nil || strings.ToLower(confirm) != "y" {
					utils.ExitWithError(utils.ErrorUser, "Aborted", nil)
					return
				}
			}

			// Invoke the /app/:name endpoint to delete the app
			err := utils.RequestJSON(utils.RequestOpts{
				Authorization: auth,
				Client:        client,
				Method:        utils.RequestDELETE,
				StatusCode:    http.StatusNoContent,
				URL:           baseURL + "/app/" + app,
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
		},
	}
	appCmd.AddCommand(c)

	// Flags
	c.Flags().BoolVarP(&yes, "yes", "", false, "do not ask for confirmation")
	c.Flags().StringVarP(&app, "app", "a", "", "name of the app to remove")
	c.MarkFlagRequired("app")

	// Add shared flags
	addSharedFlags(c)
}
