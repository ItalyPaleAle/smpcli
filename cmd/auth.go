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

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with a node",
	Long: `The auth namespace contains the commands to authenticate stkcli with a Statiko node.

The CLI supports two authentication methods:

- ` + "`" + `psk` + "`" + `: pre-shared key
  A key (passphrase) used to authenticate users. The key is stored in the node's configuration file, and is transmitted by clients in the header of API calls. Clients are authenticated if the key they send matches the one in the node's configuration.
  Note that the key is not hashed nor encrypted, so using TLS to connect to nodes is strongly recommended.

- ` + "`" + `azuread` + "`" + `: Azure AD account
- ` + "`" + `auth0` + "`" + `: Auth0
  Clients are authenticated by passing an OAuth token to the node in the header of API calls, as obtained from an Azure AD or Auth0 application. Accounts must be added to the services' directory to be granted permission to use the app.
  This method allows for tighter control over authorized users, and relies on authorization tokens which have a shorter lifespan.

Note that your Statiko nodes might not be configured to support all authentication methods.
If you're the admin of a Statiko node, please refer to the documentation for configuring authentication methods.

Please also note that, in lieu of authorizing stkcli with one of the commands above, you can pass the value for the Authorization header in the REST calls (either the pre-shared key or an OAuth access token) using the ` + "`" + `NODE_KEY` + "`" + ` environmental variable, for each command (e.g. ` + "`" + `NODE_KEY=my-psk stkcli site list` + "`" + `).
`,
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(authCmd)
}
