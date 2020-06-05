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
	"regexp"

	"github.com/spf13/cobra"

	"github.com/statiko-dev/stkcli/utils"
)

func init() {
	var (
		name        string
		certificate string
		key         string
		force       bool
	)

	c := &cobra.Command{
		Use:   "add",
		Short: "Import a new TLS certificate",
		Long: `Imports a new TLS certificate and stores it in the cluster's state.

You must provide a path to a PEM-encoded certificate and key using the ` + "`" + `--certificate` + "`" + ` and ` + "`" + `--key` + "`" + ` flags respectively.

The ` + "`" + `--name` + "`" + ` flag is the name of the TLS certificate used as identifier only.
`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Name
			certNameRegEx := regexp.MustCompile("^([a-z][a-z0-9\\.\\-]*)$")
			if !certNameRegEx.MatchString(name) {
				utils.ExitWithError(utils.ErrorUser, "Certificate name must contain letters, numbers, dots and dashes only, and it must begin with a letter", nil)
				return
			}

			// Certificate and key
			certData, err := ioutil.ReadFile(certificate)
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Error while reading TLS certificate file", err)
				return
			}
			keyData, err := ioutil.ReadFile(key)
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Error while reading TLS key file", err)
				return
			}

			// Request body
			reqBody := &certificateAddRequestModel{
				Name:        name,
				Certificate: string(certData),
				Key:         string(keyData),
				Force:       force,
			}
			buf := new(bytes.Buffer)
			err = json.NewEncoder(buf).Encode(reqBody)
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Error while encoding to JSON", err)
				return
			}

			// Invoke the /certificate endpoint and add the certificate
			err = utils.RequestJSON(utils.RequestOpts{
				Authorization:   auth,
				Body:            buf,
				BodyContentType: "application/json",
				Client:          client,
				Method:          utils.RequestPOST,
				StatusCode:      http.StatusNoContent,
				URL:             baseURL + "/certificate",
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
		},
	}
	certificateCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&name, "name", "n", "", "name for the certificate")
	c.MarkFlagRequired("name")
	c.Flags().StringVarP(&certificate, "certificate", "c", "", "path to TLS certificate file")
	c.MarkFlagRequired("certificate")
	c.Flags().StringVarP(&key, "key", "k", "", "path to TLS key file")
	c.MarkFlagRequired("key")
	c.Flags().BoolVarP(&force, "force", "f", false, "force adding invalid/expired certificates")

	// Add shared flags
	addSharedFlags(c)
}
