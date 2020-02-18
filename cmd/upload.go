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
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload apps and certificates",
	Long: `The upload namespace contains commands to conveniently upload app bundles and TLS certificates.

IMPORTANT: In order to use these commands, you must have the Azure CLI installed and you must be authenticated to the Azure subscription where the Key Vault resides (with ` + "`" + `az login` + "`" + `). Additionally, your Azure account must have the following permissions in the Key Vault's data plane: keys (create, update, import, sign), certificate (create, update, import).
`,
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
