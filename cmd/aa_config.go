// +build !docsgen

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
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/ItalyPaleAle/stkcli/utils"
)

// This is the first init that is executed
func init() {
	// Load config
	if err := loadConfig(); err != nil {
		panic(err)
	}
}

// Load configuration
func loadConfig() error {
	// Get the home directory
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	// Ensure the config folder exists
	folder := filepath.FromSlash(home + "/.stkcli")
	if err := utils.EnsureFolder(folder); err != nil {
		return err
	}

	// Load the config file
	file := folder + "/config.yaml"
	viper.SetConfigType("yaml")
	viper.SetConfigFile(file)

	// Set defaults
	viper.SetDefault("node", "")
	viper.SetDefault("port", 2265)
	viper.SetDefault("insecure", false)
	viper.SetDefault("http", false)

	// Read in the config file if it exists
	exists, err := utils.FileExists(file)
	if err != nil {
		return err
	}
	if exists {
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
	}

	return nil
}
