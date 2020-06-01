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

// dhParamsCmd represents the app command
var dhParamsCmd = &cobra.Command{
	Use:   "dhparams",
	Short: "Set DH parameters for the cluster",
	Long: `The dhparams namespace contains commands to set new Diffie-Hellman parameters for the cluster and check the status of the ones currently in use.
`,
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(dhParamsCmd)
}
