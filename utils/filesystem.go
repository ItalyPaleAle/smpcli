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
	"os"
)

// PathExists returns true if the path exists on disk
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// EnsureFolder creates a folder if it doesn't exist already
func EnsureFolder(path string) error {
	exists, err := PathExists(path)
	if err != nil {
		return err
	} else if !exists {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
