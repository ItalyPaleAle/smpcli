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
	"strconv"

	"github.com/ItalyPaleAle/stkcli/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	c := &cobra.Command{
		Use:   "set",
		Short: "Set default connection options",
		Long: `Set the value for the shared flags that will be used as default in all commands:

- '--node address' or '-n address':
  Sets the address (IP or hostname) of the node to connect to.
  This option is required.
- '--port port' or '-P port':
  If set, will communicate with the node using the port specified.
  System default: 2265
- '--insecure' or '-k' (boolean):
  If set, disables TLS certificate validation when communicating with the node (e.g. to use self-signed certificates).
  System default: false (requires valid TLS certificate)
- '--http' or '-S' (boolean):
  If set, communicates with the node using unencrypted HTTP.
  This option is considered insecure, and should only be used if the node is 'localhost', or if you're connecting to the node over an already-encrypted tunnel (e.g. VPN or SSH port forwarding).
  System default: false (use TLS)

Note that calling the 'defaults set' command overrides the default values for all the four flags above. If those values are not set, the system defaults are used. 
`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// The node address is required
			if optAddress == "" {
				utils.ExitWithError(utils.ErrorUser, "A value for the --node/-n flag is required", nil)
				return
			}

			// Port must be a number
			var port int64 = 2265
			if optPort != "" {
				port, err = strconv.ParseInt(optPort, 10, 32)
				if err != nil || port < 1 {
					utils.ExitWithError(utils.ErrorUser, "Invalid value for the --port/-P flag", nil)
					return
				}
			}

			// Set the values for shared flags
			viper.Set("node", optAddress)
			viper.Set("port", port)
			viper.Set("insecure", optInsecure)
			viper.Set("http", optNoTLS)

			// Save
			err = viper.WriteConfig()
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Error while saving configuration", err)
				return
			}
		},
	}
	defaultsCmd.AddCommand(c)

	// Add shared flags
	addSharedFlags(c)
}
