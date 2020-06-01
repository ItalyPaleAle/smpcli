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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/statiko-dev/stkcli/utils"
)

func init() {
	var (
		file string
	)

	c := &cobra.Command{
		Use:   "set",
		Short: "Sets new DH parameters",
		Long: `Sets new Diffie-Hellman parameters for the cluster. If the cluster is currently re-generating them, this interrupts the operation.

The --` + "`" + `file` + "`" + ` flag is the path to a PEM-encoded file containing DH parameters.
`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Read file
			pemData, err := ioutil.ReadFile(file)
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Error while reading DH parameters file", err)
				return
			}

			// Request body
			reqBody := &dhParamsSetRequestModel{
				DHParams: string(pemData),
			}
			buf := new(bytes.Buffer)
			err = json.NewEncoder(buf).Encode(reqBody)
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Error while encoding to JSON", err)
				return
			}

			// Invoke the /dhparams endpoint and set the DH parameters
			err = utils.RequestJSON(utils.RequestOpts{
				Authorization:   auth,
				Body:            buf,
				BodyContentType: "application/json",
				Client:          client,
				Method:          utils.RequestPOST,
				StatusCode:      http.StatusNoContent,
				URL:             baseURL + "/dhparams",
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
		},
	}
	dhParamsCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&file, "file", "f", "", "path to DH parameters file")
	c.MarkFlagRequired("file")

	// Add shared flags
	addSharedFlags(c)
}
