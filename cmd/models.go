package cmd

import (
	"time"
)

// GET /status
type statusResponseModelApp struct {
	ID         string     `json:"id"` // UUID that doesn't need parsing
	Domain     string     `json:"domain"`
	AppName    *string    `json:"appName"`
	AppVersion *string    `json:"appVersion"`
	Updated    *time.Time `json:"updated"`
}
type statusResponseModel struct {
	Apps   []statusResponseModelApp `json:"apps"`
	Health []interface{}            `json:"health"`
}

// GET /info
type infoResponseModel struct {
	AuthMethod string `json:"authMethod"`
	Hostname   string `json:"hostname"`
	Version    string `json:"version"`
}

// POST /adopt
type adoptResponseModel struct {
	Message string `json:"message"`
}

// POST /site
type siteAddRequestModel struct {
	Domain         string   `json:"domain"`
	Aliases        []string `json:"aliases"`
	TLSCertificate string   `json:"tlsCertificate"`
	ClientCaching  bool     `json:"clientCaching"`
}
type siteAddResponseModel struct {
	ID             string   `json:"id"` // UUID that doesn't need parsing
	Domain         string   `json:"domain"`
	Aliases        []string `json:"aliases"`
	TLSCertificate string   `json:"tlsCertificate"`
	ClientCaching  bool     `json:"clientCaching"`
}
