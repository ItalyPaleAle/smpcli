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
	"github.com/statiko-dev/stkcli/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	c := &cobra.Command{
		Use:   "remove",
		Short: "Remove all default connection options",
		Long: `Removes all default connection options that were set with ` + "`" + `defaults set` + "`" + `, and goes back to the system defaults.

After invoking this command, you're required to specify the ` + "`" + `--node` + "`" + ` (or ` + "`" + `-n` + "`" + `) flag for all stkcli commands that interact with a node.
`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			// Reset all values to the default
			viper.Set("node", "")
			viper.Set("port", 2265)
			viper.Set("insecure", false)
			viper.Set("http", false)

			// Save
			err := viper.WriteConfig()
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Error while saving configuration", err)
				return
			}
		},
	}
	defaultsCmd.AddCommand(c)
}
