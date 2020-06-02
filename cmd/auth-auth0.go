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

func init() {
	c := &cobra.Command{
		Use:   "auth0",
		Short: "Authenticate using Auth0",
		Long: `Launches a web browser to authenticate with the Auth0 application connected to the node, then stores the authentication token. This command manages the entire authentication workflow for the user, and it requires a desktop environment running on the client's machine.

The Auth0 application is defined in the node's configuration. Users must be part of the Auth0 directory and have permissions to use the app.

Once you have authenticated with Auth0, the client obtains an OAuth token which it uses to authorize API calls with the node. Tokens have a limited lifespan, which is configurable by the admin (stkcli supports automatically refreshing tokens when possible).
`,
		DisableAutoGenTag: true,

		// Get the command for this authentication method
		Run: openIDAuthCommand("auth0"),
	}

	authCmd.AddCommand(c)

	// Add shared flags
	addSharedFlags(c)
}
