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

	"github.com/spf13/cobra"
)

func init() {
	// Config command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Usage: config get/set [key] [value]")
		},
	}
	rootCmd.AddCommand(configCmd)

	// Config get
	configCmd.AddCommand(&cobra.Command{
		Use:   "set",
		Short: "Set configuration",
		Long:  ``,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)
		},
	})
}
