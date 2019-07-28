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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
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

			// Invoke the /state endpoint and get the state
			req, err := http.NewRequest("GET", baseURL+"/state", nil)
			if err != nil {
				fmt.Println("[Fatal error]\nCould not build the request:", err)
				return
			}
			req.Header.Set("Authorization", sharedKey)
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("[Fatal error]\nRequest failed:", err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				b, _ := ioutil.ReadAll(resp.Body)
				fmt.Printf("[Server error]\n%d: %s\n", resp.StatusCode, string(b))
				return
			}

			// If we have a file, write the response to disk
			if len(outFile) != 0 {
				out, err := os.Create("filename.ext")
				if err != nil {
					fmt.Println("[Fatal error]\nCannot create file:", err)
					return
				}
				defer out.Close()
				io.Copy(out, resp.Body)
			} else {
				// Write to stdout
				io.Copy(os.Stdout, resp.Body)
			}
		},
	}
	stateCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&outFile, "out", "o", "", "Output file where to store state (if not set, print to stdout)")

	// Add shared flags
	addSharedFlags(c)
}
