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
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	var (
		domain          string
		aliases         []string
		tlsCertificate  string
		noClientCaching bool
	)

	c := &cobra.Command{
		Use:   "add [name]",
		Short: "Add a new site",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()

			// Get the shared key
			sharedKey, found, err := nodeStore.GetSharedKey(optAddress)
			if err != nil {
				fmt.Println("[Fatal error]\nError while reading node store:", err)
				return
			}
			if !found {
				fmt.Printf("[Error]\nNo authentication data for the domain %s; please make sure you've executed the 'auth' command.\n", optAddress)
				return
			}

			// Request body
			reqBody := &siteAddRequestModel{
				Domain:         domain,
				Aliases:        aliases,
				TLSCertificate: tlsCertificate,
				ClientCaching:  !noClientCaching,
			}
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(reqBody)

			// Invoke the /site endpoint and add the site
			req, err := http.NewRequest("POST", baseURL+"/site", buf)
			if err != nil {
				fmt.Println("[Fatal error]\nCould not build the request:", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", sharedKey)
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("[Fatal error]\nRequest failed:", err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				b, _ := ioutil.ReadAll(resp.Body)
				fmt.Printf("[Server error]\n%d: %s\n", resp.StatusCode, string(b))
				return
			}

			// Parse the response
			var r siteAddResponseModel
			if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
				fmt.Println("[Fatal error]\nInvalid JSON response:", err)
				return
			}

			fmt.Println(siteAddResponseModelFormat(&r))
		},
	}
	siteCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&domain, "domain", "d", "", "Primary domain name")
	c.MarkFlagRequired("domain")
	c.Flags().StringArrayVarP(&aliases, "alias", "a", []string{}, "Alias domain (can be used multiple times)")
	c.Flags().StringVarP(&tlsCertificate, "certificate", "c", "", "Name of the TLS certificate")
	c.Flags().BoolVar(&noClientCaching, "no-client-caching", false, "Disable setting the Cache-Control for static files")

	// Add shared flags
	addSharedFlags(c)
}
