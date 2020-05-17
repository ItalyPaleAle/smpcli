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
type statusResponseModelStore struct {
	Healthy bool `json:"healthy"`
}
type statusResponseModel struct {
	NodeName string                      `json:"name"`
	Nginx    statusResponseModelNginx    `json:"nginx"`
	Sync     statusResponseModelSync     `json:"sync"`
	Store    statusResponseModelStore    `json:"store"`
	Health   []statusResponseModelHealth `json:"health"`
}

// GET /clusterstatus (cluster status)
type clusterStatusResponseModel map[string]*statusResponseModel

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
	Domain  string                `json:"domain"`
	Aliases []string              `json:"aliases"`
	TLS     *siteTLSConfiguration `json:"tls"`
}

// PATCH /site/<domain> (site set)
type siteSetRequestModel struct {
	Aliases []string              `json:"aliases"`
	TLS     *siteTLSConfiguration `json:"tls"`
}

// GET /site/<domain> (site get)
type siteGetResponseModelApp struct {
	// App details
	Name string `json:"name"`
}
type siteGetResponseModel struct {
	Domain  string                   `json:"domain"`
	Aliases []string                 `json:"aliases"`
	TLS     *siteTLSConfiguration    `json:"tls"`
	App     *siteGetResponseModelApp `json:"app"`
}

// GET /site (site list)
type siteListResponseModel []siteGetResponseModel

// POST /site/<domain>/deploy (deploy app)
type deployRequestModel struct {
	Name string `json:"name"`
}

// POST /uploadauth (request Azure Storage SAS token)
type uploadAuthRequestModel struct {
	Name string `json:"name"`
}
type uploadAuthResponseModel struct {
	ArchiveURL   string `json:"archiveUrl"`
	SignatureURL string `json:"signatureUrl"`
}

// Common
type siteTLSConfiguration struct {
	Type        string `json:"type"`
	Certificate string `json:"cert"`
	Version     string `json:"ver"`
}

const (
	TLSCertificateImported      = "imported"
	TLSCertificateAzureKeyVault = "akv"
	TLSCertificateSelfSigned    = "selfsigned"
	TLSCertificateACME          = "acme"
)
