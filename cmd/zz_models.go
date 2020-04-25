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
	"time"
)

// GET /status (status)
type statusResponseModelSync struct {
	Running   bool       `json:"running"`
	LastSync  *time.Time `json:"lastSync"`
	SyncError string     `json:"syncError"`
}
type statusResponseModelHealth struct {
	Domain  string     `json:"domain"`
	App     *string    `json:"app"`
	Healthy bool       `json:"healthy"`
	Error   *string    `json:"error"`
	Time    *time.Time `json:"time"`
}
type statusResponseModelNginx struct {
	Running bool `json:"running"`
}
type statusResponseModel struct {
	Nginx  statusResponseModelNginx    `json:"nginx"`
	Sync   statusResponseModelSync     `json:"sync"`
	Health []statusResponseModelHealth `json:"health"`
}

// GET /info (auth)
type infoResponseModelAzureAD struct {
	AuthorizeURL string `json:"authorizeUrl"`
	TokenURL     string `json:"tokenUrl"`
	ClientID     string `json:"clientId"`
}
type infoResponseModel struct {
	AuthMethods []string                  `json:"authMethods"`
	AzureAD     *infoResponseModelAzureAD `json:"azureAD"`
	Hostname    string                    `json:"hostname"`
	Version     string                    `json:"version"`
}

// POST /site (site add)
type siteAddRequestModel struct {
	Domain         string   `json:"domain"`
	Aliases        []string `json:"aliases"`
	TLSCertificate string   `json:"tlsCertificate"`
}

// PATCH /site/<domain> (site set)
type siteSetRequestModel struct {
	Aliases        []string `json:"aliases"`
	TLSCertificate string   `json:"tlsCertificate"`
}

// GET /site/<domain> (site get)
type siteGetResponseModelApp struct {
	// App details
	Name    string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
}
type siteGetResponseModel struct {
	Domain                   string                   `json:"domain"`
	Aliases                  []string                 `json:"aliases"`
	TLSCertificateSelfSigned bool                     `json:"tlsCertificateSelfSigned"`
	TLSCertificate           string                   `json:"tlsCertificate"`
	Error                    *string                  `json:"error"`
	App                      *siteGetResponseModelApp `json:"app"`
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

// POST /uploadauth (request Azure Storage SAS token)
type uploadAuthRequestModel struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
type uploadAuthResponseModel struct {
	ArchiveURL   string `json:"archiveUrl"`
	SignatureURL string `json:"signatureUrl"`
}
