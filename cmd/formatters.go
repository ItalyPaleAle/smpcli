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
	"strings"
)

func siteAddResponseModelFormat(m *siteAddResponseModel) (result string) {
	clientCaching := "no"
	if m.ClientCaching {
		clientCaching = "yes"
	}
	aliases := strings.Join(m.Aliases, ", ")
	result = fmt.Sprintf(`ID:            %s
Domain:        %s
Aliases:       %s
TLSCert:       %s
ClientCaching: %s`, m.ID, m.Domain, aliases, m.TLSCertificate, clientCaching)
	return
}

func siteListResponseModelFormat(m siteListResponseModel) (result string) {
	result = ""
	l := len(m)
	for i := 0; i < l; i++ {
		result += siteAddResponseModelFormat(&m[i])
		if i < l-1 {
			result += "\n\n"
		}
	}
	return
}
