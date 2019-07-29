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
	"time"
)

// GET /status (status)
type statusResponseModelSync struct {
	Running  bool       `json:"running"`
	LastSync *time.Time `json:"lastSync"`
}
type statusResponseModelHealth struct {
	Domain       string     `json:"domain"`
	App          *string    `json:"app"`
	StatusCode   *int       `json:"status"`
	ResponseSize *int       `json:"size"`
	Error        *string    `json:"error"`
	Time         *time.Time `json:"time"`
}
type statusResponseModel struct {
	Sync   statusResponseModelSync     `json:"sync"`
	Health []statusResponseModelHealth `json:"health"`
}

// GET /info (auth)
type infoResponseModel struct {
	AuthMethod string `json:"authMethod"`
	Hostname   string `json:"hostname"`
	Version    string `json:"version"`
}

// POST /site (site add)
type siteAddRequestModel struct {
	Domain         string   `json:"domain"`
	Aliases        []string `json:"aliases"`
	TLSCertificate string   `json:"tlsCertificate"`
	ClientCaching  bool     `json:"clientCaching"`
}

// GET /site/<domain> (site get)
type siteGetResponseModelApp struct {
	// App details
	Name    string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
}
type siteGetResponseModel struct {
	Domain         string                   `json:"domain"`
	Aliases        []string                 `json:"aliases"`
	TLSCertificate string                   `json:"tlsCertificate"`
	ClientCaching  bool                     `json:"clientCaching"`
	Error          *string                  `json:"error"`
	App            *siteGetResponseModelApp `json:"app"`
}

// GET /site (site list)
type siteListResponseModel []siteGetResponseModel

// POST /site/<domain>/deploy (deploy app)
type deployRequestModel struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
type deployResponseModel struct {
	DeploymentID string     `json:"id"`
	SiteID       string     `json:"site"`
	AppName      string     `json:"app"`
	AppVersion   string     `json:"version"`
	Status       string     `json:"status"`
	Error        *string    `json:"deploymentError"`
	Time         *time.Time `json:"time"`
}
