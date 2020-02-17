// +build generator

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

// This file is built only when trying to generate the documentation for the CLI
// go run -tags generator .

package cmd

import (
	"github.com/spf13/cobra/doc"
)

func init() {
	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		panic(err)
	}
}
