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

package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// Format of the nodes.json document
type nodeProperties struct {
	SharedKey    string `json:"sharedKey,omitempty"`
	IDToken      string `json:"idToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ClientID     string `json:"clientId,omitempty"`
	TokenURL     string `json:"tokenUrl,omitempty"`
}
type nodeDocument map[string]*nodeProperties

// NodeStore class for managing the node store
type NodeStore struct {
	path string
}

// Init the object
func (s *NodeStore) Init() error {
	// Get the home directory
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	// Ensure the folder exists
	storeFolder := filepath.FromSlash(home + "/.stkcli")
	if err := EnsureFolder(storeFolder); err != nil {
		return err
	}

	// File name
	s.path = filepath.FromSlash(storeFolder + "/nodes.json")

	return nil
}

// GetAuthToken returns the value for the Authorization header
// It will throw an error and terminate the app if there's no token or if the auth token has expired and can't be refreshed
func (s *NodeStore) GetAuthToken(address string) string {
	// If we have the NODE_KEY environmental variable, use that as fallback
	env := os.Getenv("NODE_KEY")

	// First, check if we have the data in the store
	document, err := s.read()
	if err != nil {
		ExitWithError(ErrorApp, "Could not read store file", err)
		return ""
	}

	// Check if we have something
	obj, foundObj := document[address]
	if !foundObj || (obj.SharedKey == "" && obj.IDToken == "" && obj.RefreshToken == "") {
		if env != "" {
			return env
		} else {
			ExitWithError(ErrorUser, "No authentication data for the node "+address+"; please make sure you've executed the 'auth' command.", nil)
			return ""
		}
	}

	// If we have a pre-shared key, we can proceed right away
	if obj.SharedKey != "" {
		return obj.SharedKey
	} else {
		// If we have an ID Token, check if it's still valid
		if CheckJWTValid(obj.IDToken) {
			return obj.IDToken
		}

		// Token has expired, so try refreshing it
		body := url.Values{}
		// No client_secret because this is a client-side app
		body.Set("client_id", obj.ClientID)
		body.Set("grant_type", "refresh_token")
		body.Set("refresh_token", obj.RefreshToken)
		body.Set("scope", "openid offline_access")

		// Request a new token
		var resp struct {
			ExpiresIn    int    `json:"expires_in"`
			IDToken      string `json:"id_token"`
			RefreshToken string `json:"refresh_token"`
		}
		err = RequestJSON(RequestOpts{
			Body:            strings.NewReader(body.Encode()),
			BodyContentType: "application/x-www-form-urlencoded",
			Method:          RequestPOST,
			Target:          &resp,
			URL:             obj.TokenURL,
		})
		if err != nil || resp.IDToken == "" || resp.RefreshToken == "" {
			ExitWithError(ErrorUser, "Your session for the node "+address+" has expired. Please authenticate again with the 'auth' command.", nil)
			return ""
		}

		// Store the updated tokens
		err = s.StoreAuthToken(address, resp.IDToken, resp.RefreshToken, obj.ClientID, obj.TokenURL)
		if err != nil {
			ExitWithError(ErrorApp, "Error while trying to save the new tokens", err)
			return ""
		}

		// Return the token
		return resp.IDToken
	}

	// We should never get here
	ExitWithError(ErrorApp, "Reaching unexpected code", nil)
	return ""
}

// StoreSharedKey adds the shared key to the store
func (s *NodeStore) StoreSharedKey(address string, sharedKey string) error {
	// Read the current file
	document, err := s.read()
	if err != nil {
		return err
	}

	// Add the item
	document[address] = &nodeProperties{
		SharedKey: sharedKey,
	}

	// Store the updated object
	if err := s.save(document); err != nil {
		return err
	}

	return nil
}

// StoreAuthToken adds the ID Token and the Refresh Token to the store
func (s *NodeStore) StoreAuthToken(address string, idToken string, refreshToken string, clientID string, tokenURL string) error {
	// Read the current file
	document, err := s.read()
	if err != nil {
		return err
	}

	// Add the item
	document[address] = &nodeProperties{
		IDToken:      idToken,
		RefreshToken: refreshToken,
		ClientID:     clientID,
		TokenURL:     tokenURL,
	}

	// Store the updated object
	if err := s.save(document); err != nil {
		return err
	}

	return nil
}

func (s *NodeStore) read() (nodeDocument, error) {
	// If file doesn't exist, return an empty document
	exists, err := PathExists(s.path)
	if err != nil {
		return nil, err
	}
	if !exists {
		data := make(nodeDocument)
		return data, nil
	}

	// Read the JSON
	bytes, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	var data nodeDocument
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *NodeStore) save(data nodeDocument) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(s.path, bytes, 0600); err != nil {
		return err
	}

	return nil
}
