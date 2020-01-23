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
	"strconv"
	"strings"
	"time"
)

// Fromat siteGetResponseModel
func siteGetResponseModelFormat(m *siteGetResponseModel) (result string) {
	aliases := strings.Join(m.Aliases, ", ")

	err := "\033[2m<nil>\033[0m"
	if m.Error != nil {
		err = *m.Error
	}

	app := "\033[2m<nil>\033[0m"
	if m.App != nil {
		app = fmt.Sprintf("%s-%s", m.App.Name, m.App.Version)
	}

	result = fmt.Sprintf(`Domain:        %s
Aliases:       %s
TLSCert:       %s
Error:         %s
App:           %s`, m.Domain, aliases, m.TLSCertificate, err, app)
	return
}

// Format siteListResponseModel
func siteListResponseModelFormat(m siteListResponseModel) (result string) {
	result = ""
	l := len(m)
	for i := 0; i < l; i++ {
		result += siteGetResponseModelFormat(&m[i])
		if i < l-1 {
			result += "\n\n"
		}
	}
	return
}

// Format statusResponseModel
func statusResponseModelFormat(m *statusResponseModel) (result string) {
	// Info (Nginx and sync status)
	nginxRunning := "no"
	if m.Nginx.Running {
		nginxRunning = "yes"
	}
	syncRunning := "no"
	if m.Sync.Running {
		syncRunning = "yes"
	}
	syncError := "\033[2m<nil>\033[0m"
	if m.Sync.SyncError != "" {
		syncError = m.Sync.SyncError
	}

	result = fmt.Sprintf("Info\n----\n"+`
Nginx is running: %s
Sync is running:  %s
Last sync time:   %s
Sync error:       %s

`, nginxRunning, syncRunning, m.Sync.LastSync.Format(time.RFC3339), syncError)

	// Health
	result += "Health\n------\n\n"

	l := len(m.Health)
	for i := 0; i < l; i++ {
		el := m.Health[i]

		if el.App != nil {
			err := "\033[2m<nil>\033[0m"
			if el.Error != nil {
				err = *el.Error
			}

			statusCode := "\033[2m<nil>\033[0m"
			if el.StatusCode != nil {
				statusCode = strconv.Itoa(*el.StatusCode)
			}

			responseSize := "\033[2m<nil>\033[0m"
			if el.ResponseSize != nil {
				responseSize = strconv.Itoa(*el.ResponseSize)
			}

			result += fmt.Sprintf(`Domain:       %s
App:          %s
StatusCode:   %s
ResponseSize: %s
Error:        %s
Time:         %s

`, el.Domain, *el.App, statusCode, responseSize, err, el.Time.Format(time.RFC3339))
		} else {
			// If there's no app deployed there's less data
			result += fmt.Sprintf(`Domain:       %s
App:          %s

`, el.Domain, "\033[2m<nil>\033[0m")
		}
	}

	return
}

// Fromat deployResponseModel
func deployResponseModelFormat(m *deployResponseModel) (result string) {
	err := "\033[2m<nil>\033[0m"
	if m.Error != nil {
		err = *m.Error
	}

	t := "\033[2m<nil>\033[0m"
	if m.Time != nil {
		t = m.Time.Format(time.RFC3339)
	}

	result = fmt.Sprintf(`DeploymentID: %s
SiteID:       %s
AppName:      %s
AppVersion:   %s
Status:       %s
Error:        %s
Time:         %s`, m.DeploymentID, m.SiteID, m.AppName, m.AppVersion, m.Status, err, t)
	return
}
