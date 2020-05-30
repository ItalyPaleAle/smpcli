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

package cmd

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/statiko-dev/stkcli/utils"
)

func init() {
	var (
		app        string
		path       string
		signingKey string
	)

	c := &cobra.Command{
		Use:   "upload",
		Short: "Upload an app or bundle",
		Long: `Uploads an app or app bundle to the node, to be stored in the node's app repository.

This command accepts four parameters:

- ` + "`" + `--path` + "`" + ` or ` + "`" + `-f` + "`" + ` is the path to a file or folder to upload
- ` + "`" + `--app` + "`" + ` or ` + "`" + `-a` + "`" + ` is the name of the name of the bundle, which can be used to identify the app when you want to deploy it in a node (do not include an extension)
- ` + "`" + `--signing-key` + "`" + ` is the path to a private RSA key used for codesigning

Paths can be folders containing your app's files; stkcli will automatically create a tar.bz2 archive for you. Alternatively, you can point the ` + "`" + `--path` + "`" + ` parameter to an existing archive (various formats are supported, including zip, tar.gz, tar.bz2, and more), and it will uploaded as-is.

App names must be unique. You cannot re-upload an app using the same file name.
`,
		DisableAutoGenTag: true,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()
			auth := nodeStore.GetAuthToken(optAddress)

			// Check if the path exists
			exists, err := utils.PathExists(path)
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Error while reading filesystem", err)
				return
			}
			if !exists {
				utils.ExitWithError(utils.ErrorUser, "File or folder not found", err)
				return
			}

			// App name and bundle name
			app = strings.ToLower(app)
			bundleName := ""

			// Check if the path is a folder
			folder, err := utils.FolderExists(path)
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Filesystem error", err)
				return
			}
			if folder {
				// Bundle is the app's name and ".tar.bz2" extension
				bundleName = app + ".tar.bz2"
			} else {
				// It's a file, so check the file type
				pathLc := strings.ToLower(path)
				switch {
				case strings.HasSuffix(pathLc, ".zip"),
					strings.HasSuffix(pathLc, ".tar"),
					strings.HasSuffix(pathLc, ".tgz"),
					strings.HasSuffix(pathLc, ".tsz"),
					strings.HasSuffix(pathLc, ".txz"),
					strings.HasSuffix(pathLc, ".rar"):
					bundleName = app + pathLc[(len(pathLc)-4):]
				case strings.HasSuffix(pathLc, ".tar.bz2"),
					strings.HasSuffix(pathLc, ".tar.lz4"):
					bundleName = app + pathLc[(len(pathLc)-8):]
				case strings.HasSuffix(pathLc, ".tar.gz"),
					strings.HasSuffix(pathLc, ".tar.sz"),
					strings.HasSuffix(pathLc, ".tar.xz"):
					bundleName = app + pathLc[(len(pathLc)-7):]
				case strings.HasSuffix(pathLc, ".tbz2"),
					strings.HasSuffix(pathLc, ".tlz4"):
					bundleName = app + pathLc[(len(pathLc)-5):]
				default:
					utils.ExitWithError(utils.ErrorUser, "Invalid file type", nil)
					return
				}
			}

			// If it's a folder, create an archive; upload bundles as-is
			var file io.ReadCloser
			if folder {
				// Create a tar.bz2 archive
				var w *io.PipeWriter
				file, w = io.Pipe()
				go func() {
					if err := utils.TarBZ2(path, w); err != nil {
						utils.ExitWithError(utils.ErrorApp, "Error while creating a tar.bz2 archive", err)
						panic(1)
					}
					w.Close()
				}()
			} else {
				// Get a buffer reader
				file, err = os.Open(path)
				if err != nil {
					utils.ExitWithError(utils.ErrorApp, "Error while reading file", err)
					return
				}
			}

			// The stream is split between two readers: one for the hashing, one for writing the stream to disk
			h := sha256.New()
			tee := io.TeeReader(file, h)

			// Upload the app's file
			// This also makes the stream proceed so the hash is calculated
			// Start by creating the body as multipart/form-data
			pr, pw := io.Pipe()
			mpw := multipart.NewWriter(pw)
			go func() {
				partw, err := mpw.CreateFormFile("file", bundleName)
				if err != nil {
					utils.ExitWithError(utils.ErrorApp, "Error while preparing request", err)
					return
				}
				_, err = io.Copy(partw, tee)
				if err != nil {
					utils.ExitWithError(utils.ErrorApp, "Error while preparing request", err)
					return
				}
				err = mpw.Close()
				if err != nil {
					utils.ExitWithError(utils.ErrorApp, "Error while preparing request", err)
					return
				}
				err = pw.Close()
				if err != nil {
					utils.ExitWithError(utils.ErrorApp, "Error while preparing request", err)
					return
				}
				err = file.Close()
				if err != nil {
					utils.ExitWithError(utils.ErrorApp, "Error while preparing request", err)
					return
				}
			}()

			// Invoke the /app endpoint
			err = utils.RequestJSON(utils.RequestOpts{
				Authorization:   auth,
				Body:            pr,
				BodyContentType: mpw.FormDataContentType(),
				Client:          client,
				Method:          utils.RequestPOST,
				StatusCode:      http.StatusNoContent,
				URL:             baseURL + "/app",
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}

			fmt.Println("Uploaded app's bundle")

			// Calculate the SHA256 hash
			hashed := h.Sum(nil)
			metadata := appMetadataRequestModel{
				Hash: base64.StdEncoding.EncodeToString(hashed),
			}

			// If we have a key, calculate the digital signature
			if signingKey != "" {
				// Load key
				privateKey := loadRSAPrivateKey(signingKey)
				if privateKey == nil {
					utils.ExitWithError(utils.ErrorApp, "Could not load RSA private key", nil)
					return
				}

				// Calculate the signature
				rng := rand.Reader
				signatureBytes, err := rsa.SignPKCS1v15(rng, privateKey, crypto.SHA256, hashed)
				if err != nil {
					utils.ExitWithError(utils.ErrorApp, "Error while creating signature", err)
					return
				}

				// Convert the signature to base64
				metadata.Signature = base64.StdEncoding.EncodeToString(signatureBytes)
			}

			// Body
			buf := new(bytes.Buffer)
			err = json.NewEncoder(buf).Encode(metadata)
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Error while encoding to JSON", err)
				return
			}

			// Invoke the /app/:name endpoint and save the metadata
			err = utils.RequestJSON(utils.RequestOpts{
				Authorization:   auth,
				Body:            buf,
				BodyContentType: "application/json",
				Client:          client,
				Method:          utils.RequestPOST,
				StatusCode:      http.StatusNoContent,
				URL:             baseURL + "/app/" + bundleName,
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}

			fmt.Println("Stored bundle's metadata")
			fmt.Println("Done:", bundleName)
		},
	}
	appCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&app, "app", "a", "", "app bundle name, with no extension (required)")
	c.MarkFlagRequired("app")
	c.Flags().StringVarP(&path, "path", "f", "", "path to local file or folder to bundle (required)")
	c.MarkFlagRequired("path")
	c.Flags().StringVarP(&signingKey, "signing-key", "s", "", "path to a RSA private key for code signing")

	// Add shared flags
	addSharedFlags(c)
}
