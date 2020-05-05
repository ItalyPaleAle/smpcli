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

	"github.com/spf13/cobra"

	"github.com/statiko-dev/stkcli/utils"
)

func init() {
	c := &cobra.Command{
		Use:   "sync",
		Short: "Triggers a sync of the node",
		Long: `Triggers a sync operation in the node, to ensure that the node's state matches the desired one.
`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Invoke the /sync endpoint and trigger a sync
			err := utils.RequestJSON(utils.RequestOpts{
				Authorization: auth,
				Client:        client,
				Method:        utils.RequestPOST,
				StatusCode:    http.StatusNoContent,
				URL:           baseURL + "/sync",
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
		},
	}
	stateCmd.AddCommand(c)

	// Add shared flags
	addSharedFlags(c)
}
