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
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/ItalyPaleAle/stkcli/utils"
)

func init() {
	var (
		stateFile string
	)

	c := &cobra.Command{
		Use:   "set",
		Short: "Restores the state of a node",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Read the file if we have one
			var stateBuf io.Reader
			if len(stateFile) != 0 {
				exists, err := utils.PathExists(stateFile)
				if err != nil {
					utils.ExitWithError(utils.ErrorApp, "Error while checking file", err)
					return
				}
				if !exists {
					utils.ExitWithError(utils.ErrorUser, "Files does not exist", nil)
					return
				}
				state, err := ioutil.ReadFile(stateFile)
				if err != nil {
					utils.ExitWithError(utils.ErrorApp, "Error while reading file", err)
					return
				}
				if state == nil || len(state) == 0 {
					utils.ExitWithError(utils.ErrorUser, "Files is empty", nil)
					return
				}
				stateBuf = bytes.NewBuffer(state)
			} else {
				// Read from stdin
				stateBuf = os.Stdin
			}

			// Invoke the /state endpoint
			err := utils.RequestJSON(utils.RequestOpts{
				Authorization:   auth,
				Body:            stateBuf,
				BodyContentType: "application/json",
				Client:          client,
				Method:          utils.RequestPOST,
				StatusCode:      http.StatusNoContent,
				URL:             baseURL + "/state",
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
		},
	}

	stateCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&stateFile, "file", "f", "", "File containing the desired state (if not set, read from)")

	// Add shared flags
	addSharedFlags(c)
}
