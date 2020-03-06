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
	"github.com/spf13/cobra"
)

// stateCmd represents the state command
var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Get or restore state",
	Long: `The state namespace contains commands used to manage the state of the node: dump to a file, or restore from a local file.

The node's state is a JSON document containing the list of apps and sites configured. You can dump it to a local file and restore it on the same node or a different one.

The state file can contains secrets (e.g self-signed TLS certificates) which are encrypted with the symmetric key from the node's configuration file.
`,
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(stateCmd)
}
