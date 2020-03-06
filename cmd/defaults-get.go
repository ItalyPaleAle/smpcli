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
	"fmt"

	"github.com/ItalyPaleAle/stkcli/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	c := &cobra.Command{
		Use:   "get",
		Short: "Show all default connection options",
		Long: `Shows all the default flags that are used to connect to a node.
The output of the command resembles the flags that users would pass on the command line to stkcli.
`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			// Check if we have defaults
			node := viper.GetString("node")
			if node == "" {
				utils.ExitWithError(utils.ErrorUser, "There are no user defaults set", nil)
				return
			}

			// Get the port
			port := viper.GetInt("port")
			if port == 2265 {
				// Do not show the default port
				port = 0
			}

			// Show the flags
			fmt.Printf("--node %s ", node)
			if port > 0 {
				fmt.Printf("--port %d ", port)
			}
			if viper.GetBool("insecure") {
				fmt.Print("--insecure ")
			}
			if viper.GetBool("http") {
				fmt.Print("--http ")
			}
			fmt.Print("\n")
		},
	}
	defaultsCmd.AddCommand(c)
}
