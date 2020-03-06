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
	"errors"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/ItalyPaleAle/stkcli/utils"
)

func init() {
	c := &cobra.Command{
		Use:   "psk",
		Short: "Authenticate using a pre-shared key",
		Long: `Sets the pre-shared key used to authenticate API calls to a node.

The pre-shared key is defined in the node's configuration, and clients are authenticated if they send the same key in the header of API calls.
Note that the key is not hashed nor encrypted, so using TLS to connect to nodes is strongly recommended.
`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()

			// Invoke the /info endpoint to see what's the authentication method
			var rInfo infoResponseModel
			err := utils.RequestJSON(utils.RequestOpts{
				Client: client,
				Target: &rInfo,
				URL:    baseURL + "/info",
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}

			// Ensure the node supports pre-shared key authentication
			if !utils.SliceContainsString(rInfo.AuthMethods, "psk") {
				utils.ExitWithError(utils.ErrorUser, "This node does not support authenticating with a pre-shared key", nil)
				return
			}

			// Prompt the user for the shared key
			prompt := promptui.Prompt{
				Validate: func(input string) error {
					if len(input) < 1 {
						return errors.New("Pre-shared key must not be empty")
					}
					return nil
				},
				Label: "Pre-shared key",
				Mask:  '*',
			}

			sharedKey, err := prompt.Run()
			if err != nil {
				utils.ExitWithError(utils.ErrorUser, "Pre-shared key must not be empty", nil)
				return
			}

			// Test the shared key by requesting the node's site list, invoking the /site endpoint
			// We're not requesting anything from the response
			err = utils.RequestJSON(utils.RequestOpts{
				Authorization: sharedKey,
				Client:        client,
				URL:           baseURL + "/site",
			})
			if err != nil {
				// Check if the error is a 401
				if strings.HasPrefix(err.Error(), "invalid response status code: 401") {
					utils.ExitWithError(utils.ErrorUser, "Invalid pre-shared key", nil)
				} else {
					utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				}
				return
			}

			// Store the key in the node store
			if err := nodeStore.StoreSharedKey(optAddress, sharedKey); err != nil {
				utils.ExitWithError(utils.ErrorApp, "Error while storing the pre-shared key", err)
				return
			}
		},
	}

	authCmd.AddCommand(c)

	// Add shared flags
	addSharedFlags(c)
}
