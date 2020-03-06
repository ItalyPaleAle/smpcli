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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ItalyPaleAle/stkcli/utils"
)

func init() {
	c := &cobra.Command{
		Use:               "list",
		Short:             "List sites",
		Long:              `Shows the list of all sites configured in the node.`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Invoke the /site endpoint and list sites
			var r siteListResponseModel
			err := utils.RequestJSON(utils.RequestOpts{
				Authorization: auth,
				Client:        client,
				Target:        &r,
				URL:           baseURL + "/site",
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}

			// Print the response
			fmt.Println(siteListResponseModelFormat(r))
		},
	}
	siteCmd.AddCommand(c)

	// Add shared flags
	addSharedFlags(c)
}
