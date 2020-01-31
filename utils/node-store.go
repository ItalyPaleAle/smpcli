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

package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// Format of the nodes.json document
type nodeProperties struct {
	SharedKey string `json:"sharedKey"`
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
	storeFolder := filepath.FromSlash(home + "/.smpcli")
	if err := EnsureFolder(storeFolder); err != nil {
		return err
	}

	// File name
	s.path = filepath.FromSlash(storeFolder + "/nodes.json")

	return nil
}

// GetSharedKey reads the shared key from the store
// The second returned value is a boolean that indicates if the key was found
func (s *NodeStore) GetSharedKey(address string) (sharedKey string, found bool, err error) {
	sharedKey = ""
	found = false
	err = nil

	// If we have the NODE_KEY environmental variable, use that as fallback
	env := os.Getenv("NODE_KEY")
	if env != "" {
		sharedKey = env
		found = true
	}

	// Read the current file
	document, err := s.read()
	if err != nil {
		return
	}

	// Get the value if it exists
	obj, foundObj := document[address]
	if foundObj {
		sharedKey = obj.SharedKey
		found = true
	}
	if len(sharedKey) == 0 {
		found = false
	}

	return
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
