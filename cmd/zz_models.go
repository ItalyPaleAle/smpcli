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
type infoResponseModelOpenID struct {
	AuthorizeURL string `json:"authorizeUrl"`
	TokenURL     string `json:"tokenUrl"`
	ClientID     string `json:"clientId"`
}
type infoResponseModel struct {
	AuthMethods []string                 `json:"authMethods"`
	Auth0       *infoResponseModelOpenID `json:"auth0"`
	AzureAD     *infoResponseModelOpenID `json:"azureAD"`
	Hostname    string                   `json:"hostname"`
	Version     string                   `json:"version"`
}

// POST /site (site add)
type siteAddRequestModel struct {
	Domain    string                `json:"domain,omitempty"`
	Aliases   []string              `json:"aliases,omitempty"`
	Temporary bool                  `json:"temporary,omitempty"`
	TLS       *siteTLSConfiguration `json:"tls,omitempty"`
}

// PATCH /site/<domain> (site set)
type siteSetRequestModel struct {
	Aliases []string              `json:"aliases,omitempty"`
	TLS     *siteTLSConfiguration `json:"tls,omitempty"`
}

// GET /site/<domain> (site get)
type siteGetResponseModelApp struct {
	// App details
	Name string `json:"name"`
}
type siteGetResponseModel struct {
	Domain    string                   `json:"domain"`
	Temporary bool                     `json:"temporary"`
	Aliases   []string                 `json:"aliases"`
	TLS       *siteTLSConfiguration    `json:"tls"`
	App       *siteGetResponseModelApp `json:"app"`
}

// GET /site (site list)
type siteListResponseModel []siteGetResponseModel

// POST /site/<domain>/deploy (deploy app)
type deployRequestModel struct {
	Name string `json:"name"`
}

// GET /app (app list)
type appListResponseModel []struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
}

// POST /app/<app> (set app bundle's metadata)
type appMetadataRequestModel struct {
	Signature string `json:"signature"`
	Hash      string `json:"hash"`
}

// GET /certificate (certificate list)
type certificateListResponseModel []string

// POST /certificate (site add)
type certificateAddRequestModel struct {
	Name        string `json:"name"`
	Certificate string `json:"cert" `
	Key         string `json:"key"`
	Force       bool   `json:"force"`
}

// GET /dhparams (DH params show)
type dhParamsGetResponseModel struct {
	Type       string     `json:"type"`
	Date       *time.Time `json:"date"`
	Generating bool       `json:"generating"`
}

// POST /dhparams (DH params set)
type dhParamsSetRequestModel struct {
	DHParams string `json:"dhparams"`
}

// Common
type siteTLSConfiguration struct {
	Type        string `json:"type"`
	Certificate string `json:"cert,omitempty"`
	Version     string `json:"ver,omitempty"`
}

const (
	TLSCertificateImported      = "imported"
	TLSCertificateAzureKeyVault = "akv"
	TLSCertificateSelfSigned    = "selfsigned"
	TLSCertificateACME          = "acme"
)
