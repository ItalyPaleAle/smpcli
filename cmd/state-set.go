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
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"smpcli/utils"
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

			// Read the file
			exists, err := utils.PathExists(stateFile)
			if err != nil {
				fmt.Println("[Fatal error]\nError while checking file:", err)
				return
			}
			if !exists {
				fmt.Println("[Error]\nFile does not exist.")
				return
			}
			state, err := ioutil.ReadFile(stateFile)
			if err != nil {
				fmt.Println("[Fatal error]\nError while reading file:", err)
				return
			}
			if state == nil || len(state) == 0 {
				fmt.Println("[Error]\nFile is empty.")
				return
			}
			stateBuf := bytes.NewBuffer(state)

			// Ask user to confirm a potentially destructive action
			fmt.Println("Warning: this is a potentially destructive action, that will replace the state of a node. Please type the address (hostname or IP) of the node again to confirm, without quotes.")
			prompt := promptui.Prompt{
				Label: fmt.Sprintf("Type '%s'", optAddress),
			}
			result, err := prompt.Run()
			if err != nil || result != optAddress {
				fmt.Println("Aborted")
				return
			}

			// Invoke the /state endpoint
			req, err := http.NewRequest("POST", baseURL+"/state", stateBuf)
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
			if resp.StatusCode != http.StatusNoContent {
				b, _ := ioutil.ReadAll(resp.Body)
				fmt.Printf("[Server error]\n%d: %s\n", resp.StatusCode, string(b))
				return
			}
		},
	}

	stateCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&stateFile, "file", "f", "", "File containing the desired state")
	c.MarkFlagRequired("file")

	// Add shared flags
	addSharedFlags(c)
}
