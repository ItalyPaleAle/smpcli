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
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/ItalyPaleAle/stkcli/utils"
)

func init() {
	var (
		domain         string
		aliases        []string
		tlsCertificate string
	)

	c := &cobra.Command{
		Use:   "add",
		Short: "Add a new site",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Request body
			reqBody := &siteAddRequestModel{
				Domain:         domain,
				Aliases:        aliases,
				TLSCertificate: tlsCertificate,
			}
			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(reqBody)
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Error while encoding to JSON", err)
				return
			}

			// Invoke the /site endpoint and add the site
			err = utils.RequestJSON(utils.RequestOpts{
				Authorization:   auth,
				Body:            buf,
				BodyContentType: "application/json",
				Client:          client,
				Method:          utils.RequestPOST,
				StatusCode:      http.StatusNoContent,
				URL:             baseURL + "/site",
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}
		},
	}
	siteCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&domain, "domain", "d", "", "Primary domain name")
	c.MarkFlagRequired("domain")
	c.Flags().StringArrayVarP(&aliases, "alias", "a", []string{}, "Alias domain (can be used multiple times)")
	c.Flags().StringVarP(&tlsCertificate, "certificate", "c", "", "Name of the TLS certificate")

	// Add shared flags
	addSharedFlags(c)
}
