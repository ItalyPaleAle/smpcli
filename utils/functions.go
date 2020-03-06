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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// SliceContainsString returns true if the slice of strings contains a certain string
func SliceContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// LaunchBrowser opens a web browser at a specified URL
func LaunchBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	case "linux":
		exec.Command("xdg-open", url).Start()
	default:
		fmt.Printf("Please open a web browser to this URL to authenticate:\n%s\n", url)
		return
	}
	fmt.Printf("If your browser didn't automatically open, please visit this URL to authenticate:\n%s\n", url)
}

// CheckJWTValid returns true if the JWT token is well-formed and not expired
// This function does not validate the JWT token beyond checking if it's still valid
func CheckJWTValid(jwt string) bool {
	// Split the token in its 3 parts
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return false
	}

	// Decode the header
	headerData, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	var header map[string]string
	err = json.Unmarshal(headerData, &header)
	if err != nil {
		return false
	}
	if header["typ"] != "JWT" {
		return false
	}

	// Decode the claims
	claimsData, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	var claims map[string]string
	err = json.Unmarshal(claimsData, &claims)
	if err != nil {
		return false
	}

	// Check the expiration
	now := time.Now().Unix()
	if claims["exp"] == "" {
		return false
	}
	exp, err := strconv.ParseInt(claims["exp"], 10, 64)
	if err != nil || exp < now {
		return false
	}

	return true
}
