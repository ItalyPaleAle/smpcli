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
	"path"
	"path/filepath"
	"strings"

	"github.com/dsnet/compress/bzip2"
)

// TarBZ2 creates a tar.bz2 archive from a folder
// Adapted from: https://gist.github.com/sdomino/e6bc0c98f87843bc26bb
func TarBZ2(src string, writers ...io.Writer) error {
	// Clean the source folder
	src = path.Clean(src)

	// Ensure the src actually exists before trying to tar it, and that it's a directory
	exists, err := FolderExists(src)
	if err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("path does not exist or is not a folder: %s", src)
	}

	// Create a bzip2 stream compressor
	mw := io.MultiWriter(writers...)
	bzw, err := bzip2.NewWriter(mw, &bzip2.WriterConfig{Level: 9})
	defer bzw.Close()
	if err != nil {
		return err
	}

	// Tar writer
	tw := tar.NewWriter(bzw)
	defer tw.Close()

	// Walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Exclude the root folder
		if file == src || file == src+string(os.PathSeparator) {
			return nil
		}

		// Create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// Update the name to correctly reflect the desired destination when un-taring
		header.Name = strings.TrimPrefix(file, src+string(os.PathSeparator))
		fmt.Println("Adding", header.Name)

		// Write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If this is a directory, go to the next segment
		if fi.Mode().IsDir() {
			return nil
		}

		// Open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		// Copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})
}
