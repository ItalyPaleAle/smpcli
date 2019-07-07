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
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// Config command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Usage:")
			fmt.Println("config get [key]")
			fmt.Println("config set [key] [value]")
		},
	}
	rootCmd.AddCommand(configCmd)

	// Config set
	configCmd.AddCommand(&cobra.Command{
		Use:   "set",
		Short: "Set configuration",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			// Need 2 args
			if len(args) != 2 {
				return errors.New("Usage: config set [key] [value]")
			}

			// Check if the key is valid
			key := args[0]
			if !viper.IsSet(key) {
				return errors.New("Invalid config key: " + key)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Set the config
			key := args[0]
			value := args[1]
			viper.Set(key, value)

			// Save
			viper.WriteConfig()
		},
	})

	// Config get
	configCmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get configuration",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			// Need 1 arg
			if len(args) != 1 {
				return errors.New("Usage: config get [key]")
			}

			// Check if the key is valid
			key := args[0]
			if !viper.IsSet(key) {
				return errors.New("Invalid config key: " + key)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Get the config
			key := args[0]
			value := viper.Get(key)
			fmt.Println(value)
		},
	})
}
