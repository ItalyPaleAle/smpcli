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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with a node",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		baseURL, client := getURLClient()

		// Invoke the /info endpoint to see what's the authentication method
		resp, err := client.Get(baseURL + "/info")
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
		var r infoResponseModel
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			fmt.Println("[Fatal error]\nInvalid JSON response:", err)
			return
		}

		// Ensure we have the correct data
		if r.AuthMethod != "sharedkey" {
			fmt.Println("[Fatal error]\nThe response from the server is invalid")
			return
		}

		// Prompt the user for the shared key
		prompt := promptui.Prompt{
			Validate: func(input string) error {
				if len(input) < 1 {
					return errors.New("Shared key must not be empty")
				}
				return nil
			},
			Label: "Shared key",
			Mask:  '*',
		}

		sharedKey, err := prompt.Run()
		if err != nil {
			fmt.Println("[Fatal error]\nPrompt failed:", err)
			return
		}

		// Store the key in the node store
		if err := nodeStore.StoreSharedKey(optAddress, sharedKey); err != nil {
			fmt.Println("[Fatal error]\nError while storing the shared key:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
