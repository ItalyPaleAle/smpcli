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

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [address]",
	Short: "Gets the status of a node",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		baseURL, client := getURLClient()

		// Invoke the /status endpoint to see what's the authentication method
		resp, err := client.Get(baseURL + "/status")
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
	rootCmd.AddCommand(statusCmd)
}
