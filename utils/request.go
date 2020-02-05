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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Default HTTP client (do not use this for requests to the nodes)
var defaultHTTPClient *http.Client

const (
	RequestDELETE = "DELETE"
	RequestGET    = "GET"
	RequestPATCH  = "PATCH"
	RequestPOST   = "POST"
	RequestPUT    = "PUT"
)

// RequestOpts contains the parameters for the RequestJSON function
type RequestOpts struct {
	Authorization   string
	Body            io.Reader
	BodyContentType string
	Client          *http.Client
	Method          string
	StatusCode      int
	Target          interface{} // Only used by RequestJSON
	URL             string
}

// RequestJSON fetches a JSON document from the web
func RequestJSON(opts RequestOpts) (err error) {
	// Make the request
	response, err := RequestRaw(opts)
	if err != nil {
		return err
	}
	defer response.Close()

	if opts.Target != nil {
		// Decode the JSON into the target
		err = json.NewDecoder(response).Decode(opts.Target)
		if err != nil {
			return err
		}
	}
	return nil
}

// RequestRaw fetches a document from the web and returns the stream as is
func RequestRaw(opts RequestOpts) (response io.ReadCloser, err error) {
	// Check options and default values
	if opts.URL == "" {
		return nil, errors.New("empty URL")
	}
	if opts.Client == nil {
		if defaultHTTPClient == nil {
			defaultHTTPClient = &http.Client{
				Timeout: 30 * time.Second,
			}
		}
		opts.Client = defaultHTTPClient
	}
	if opts.Method == "" {
		opts.Method = RequestGET
	}
	if opts.Method == RequestGET && opts.Body != nil {
		return nil, errors.New("cannot have a request body for GET requests")
	}
	if opts.Body != nil && opts.BodyContentType == "" {
		return nil, errors.New("must specify a content type for the body when there's a request body")
	}

	// Build the request
	req, err := http.NewRequest(opts.Method, opts.URL, opts.Body)
	if err != nil {
		return
	}
	// Set the body's Content-Type if we have a body
	if opts.Body != nil {
		req.Header.Set("Content-Type", opts.BodyContentType)
	}
	// Authorization, if any
	if opts.Authorization != "" {
		req.Header.Set("Authorization", opts.Authorization)
	}

	// Send the request
	resp, err := opts.Client.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		return nil, err
	}

	// If we're expecting a specific status code, check for that, otherwise fallback to check that we're below 400
	if (opts.StatusCode > 0 && resp.StatusCode != opts.StatusCode) || (opts.StatusCode <= 0 && resp.StatusCode >= 399) {
		b, _ := ioutil.ReadAll(resp.Body)
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		return nil, fmt.Errorf("invalid response status code: %d; content: %s", resp.StatusCode, string(b))
	}

	return resp.Body, nil
}
