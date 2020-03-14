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
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dsnet/compress/bzip2"
)

// TarBZ2 creates a tar.bz2 archive from a folder
// Source: https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func TarBZ2(src string, writers ...io.Writer) error {
	// Ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	mw := io.MultiWriter(writers...)

	bzw, err := bzip2.NewWriter(mw, &bzip2.WriterConfig{Level: 9})
	defer bzw.Close()
	if err != nil {
		return err
	}

	tw := tar.NewWriter(bzw)
	defer tw.Close()

	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		// return on any error
		if err != nil {
			return err
		}

		// Create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// Update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))
		fmt.Println("Adding", "/"+header.Name)

		// Write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// Return on non-regular files (thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update)
		if !fi.Mode().IsRegular() {
			return nil
		}

		// Open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// Copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		// Manually close here after each file operation; defering would cause each file close to wait until all operations have completed.
		f.Close()

		return nil
	})
}
