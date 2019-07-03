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
	"io/ioutil"
	"net/http"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// adoptCmd represents the adopt command
var adoptCmd = &cobra.Command{
	Use:   "adopt",
	Short: "Adopts/resets a node",
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

		// Ask user to confirm a potentially destructive action
		fmt.Println("Warning: this is a potentially destructive action, that will reset the node to a clean state. Please type the address (hostname or IP) of the node again to confirm, without quotes.")
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("Type '%s'", optAddress),
		}
		result, err := prompt.Run()
		if err != nil || result != optAddress {
			fmt.Println("Aborted")
			return
		}

		// Invoke the /adopt endpoint
		req, err := http.NewRequest("POST", baseURL+"/adopt", nil)
		req.Header.Set("Authorization", sharedKey)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("[Fatal error]\nRequest failed:", err)
			return
		}
		defer resp.Body.Close()

		bytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bytes))
	},
}

func init() {
	rootCmd.AddCommand(adoptCmd)
}
