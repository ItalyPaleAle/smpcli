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
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/ItalyPaleAle/stkcli/utils"
)

func init() {
	var (
		outFile string
	)

	c := &cobra.Command{
		Use:   "get",
		Short: "Retrieve state and save to file",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Invoke the /state endpoint and get the state
			body, err := utils.RequestRaw(utils.RequestOpts{
				Authorization: auth,
				Client:        client,
				URL:           baseURL + "/state",
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
			defer body.Close()

			// If we have a file, write the response to disk
			if len(outFile) != 0 {
				out, err := os.Create(outFile)
				if err != nil {
					utils.ExitWithError(utils.ErrorApp, "Cannot create file", err)
					return
				}
				defer out.Close()
				io.Copy(out, body)
			} else {
				// Write to stdout
				io.Copy(os.Stdout, body)
			}
		},
	}
	stateCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&outFile, "out", "o", "", "Output file where to store state (if not set, print to stdout)")

	// Add shared flags
	addSharedFlags(c)
}
