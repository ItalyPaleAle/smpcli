package utils

import (
	"encoding/json"
	"io/ioutil"
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

	// Read the current file
	document, err := s.read()
	if err != nil {
		return
	}

	// Get the value if it exists
	obj, found := document[address]
	if found {
		sharedKey = obj.SharedKey
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
