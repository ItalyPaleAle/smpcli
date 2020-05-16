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
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/spf13/cobra"

	"github.com/statiko-dev/stkcli/utils"
)

func init() {
	var (
		app         string
		path        string
		noSignature bool
	)

	// This function signs the checksum of the app's payload with the key stored in Azure Key Vault
	var signHash = func(ctx context.Context, hash []byte) (signature string, err error) {
		// Get the URL of the Key Vault, requesting it from the node
		keyVaultURL, codesignKeyName, codesignKeyVersion, err := getKeyVaultInfo()
		if err != nil {
			return
		}

		// Convert the hash to base64
		hashB64 := base64.URLEncoding.EncodeToString(hash)

		// Get the Azure Key Vault client
		akvClient := getKeyVault()
		if akvClient == nil {
			return
		}

		// Request Azure Key Vault to sign the message
		res, err := akvClient.Sign(ctx, keyVaultURL, codesignKeyName, codesignKeyVersion, keyvault.KeySignParameters{
			Algorithm: "RS256",
			Value:     &hashB64,
		})
		if err != nil {
			return
		}

		// Check the response
		if res.Result == nil || *res.Result == "" {
			err = errors.New("Empty response")
		}

		// The response is encded with Base64 with URL-encoding; we need to switch to the standard encoding
		signature = strings.ReplaceAll(*res.Result, "-", "+")
		signature = strings.ReplaceAll(signature, "_", "/")

		// Ensure that we have the proper padding
		if len(signature)%4 == 2 {
			signature += "=="
		} else if len(signature)%4 == 3 {
			signature += "="
		}

		return
	}

	// This function uploads the app's bundle to Azure Blob Storage
	// Returns true in case of success, and false if there's an error
	var uploadArchive = func(file io.Reader, sasURLs *uploadAuthResponseModel) bool {
		// URL to upload the archive to
		uApp, err := url.Parse(sasURLs.ArchiveURL)
		if err != nil {
			utils.ExitWithError(utils.ErrorApp, "Error while building app URL", err)
			return false
		}

		// Uploader client
		credential := azblob.NewAnonymousCredential()
		pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{
			Retry: azblob.RetryOptions{
				MaxTries: 3,
			},
		})
		ctx := context.Background()

		// The stream is split between two readers: one for the hashing, one for writing the stream to disk
		h := sha256.New()
		tee := io.TeeReader(file, h)

		// Access conditions for blob uploads: disallow the operation if the blob already exists
		// See: https://docs.microsoft.com/en-us/rest/api/storageservices/specifying-conditional-headers-for-blob-service-operations#Subheading1
		accessConditions := azblob.BlobAccessConditions{
			ModifiedAccessConditions: azblob.ModifiedAccessConditions{
				IfNoneMatch: "*",
			},
		}

		// Upload the app's file
		// This also makes the stream proceed so the hash is calculated
		blockBlobURL := azblob.NewBlockBlobURL(*uApp, pipeline)
		_, err = azblob.UploadStreamToBlockBlob(ctx, tee, blockBlobURL, azblob.UploadStreamToBlockBlobOptions{
			BufferSize:       3 * 1024 * 1024,
			MaxBuffers:       2,
			AccessConditions: accessConditions,
		})
		if err != nil {
			if stgErr, ok := err.(azblob.StorageError); !ok {
				utils.ExitWithError(utils.ErrorApp, "Network error while uploading the archive", err)
			} else {
				utils.ExitWithError(utils.ErrorApp, "Azure Storage error failed while uploading the archive:\n"+stgErr.Response().Status, nil)
			}
			return false
		}
		fmt.Println("Uploaded app's bundle")

		// Calculate the SHA256 hash
		hashed := h.Sum(nil)

		// Calculate the digital signature
		if !noSignature {
			signature, err := signHash(ctx, hashed)
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Cannot calculate signature", err)
				return false
			}

			// Add the signature as metadata
			metadata := azblob.Metadata{}
			metadata["signature"] = signature
			_, err = blockBlobURL.SetMetadata(ctx, metadata, azblob.BlobAccessConditions{})
			if err != nil {
				utils.ExitWithError(utils.ErrorApp, "Cannot update blob's metadata", err)
				return false
			}
		} else {
			fmt.Fprintln(os.Stderr, "\033[33mWARN: Skipping cryptographically signing the app's bundle. Nodes will not be able to verify the integrity and the origin of the code.\033[0m")
		}

		return true
	}

	c := &cobra.Command{
		Use:   "app",
		Short: "Upload an app or bundle",
		Long: `Uploads an app to the Azure Storage Account associated with the node's instance

IMPORTANT: In order to use this command, you must have the Azure CLI installed and you must be authenticated to the Azure subscription where the Key Vault resides (with ` + "`" + `az login` + "`" + `). Additionally, your Azure account must have the following permissions in the Key Vault's data plane: keys (create, update, import, sign), certificate (create, update, import).

This command accepts four parameters:

- ` + "`" + `--path` + "`" + ` or ` + "`" + `-f` + "`" + ` is the path to a file or folder to upload
- ` + "`" + `--app` + "`" + ` or ` + "`" + `-a` + "`" + ` is the name of the name of the bundle, which can be used to identify the app when you want to deploy it in a node (do not include an extension)
- ` + "`" + `--no-signature` + "`" + ` is a boolean that when present will skip calculating the checksum of the app's bundle and signing it with the codesign key

Paths can be folders containing your app's files; stkcli will automatically create a tar.bz2 archive for you. Alternatively, you can point the ` + "`" + `--path` + "`" + ` parameter to an existing archive (various formats are supported, including zip, tar.gz, tar.bz2, and more), and it will uploaded as-is.

App names must be unique. You cannot re-upload an app using the same file name.

When using ` + "`" + `--no-signature` + "`" + `, stkcli will not calculate the checksum of the app's bundle, and it will not cryptographically sign it with the codesigning key. Statiko nodes might be configured to not accept unsigned app bundles for security reasons. However, when uploading unsigned bundles, you do not need to be signed into an Azure account in the local system.
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

			// Request body for getting the SAS token for Azure Storage from the node
			reqBody := &uploadAuthRequestModel{
				Name: bundleName,
			}
			buf := new(bytes.Buffer)
			err = json.NewEncoder(buf).Encode(reqBody)
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Error while encoding to JSON", err)
				return
			}

			// Invoke the /uploadauth endpoint and request the SAS token from the node
			var sasURLs uploadAuthResponseModel
			err = utils.RequestJSON(utils.RequestOpts{
				Authorization:   auth,
				Body:            buf,
				BodyContentType: "application/json",
				Client:          client,
				Method:          utils.RequestPOST,
				Target:          &sasURLs,
				URL:             baseURL + "/uploadauth",
			})
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Request failed", err)
				return
			}

			// If it's a folder, create an archive; upload bundles as-is
			if folder {
				// Create a tar.bz2 archive
				r, w := io.Pipe()
				go func() {
					if err := utils.TarBZ2(path, w); err != nil {
						utils.ExitWithError(utils.ErrorApp, "Error while creating a tar.bz2 archive", err)
						panic(1)
					}
					w.Close()
				}()

				// Upload the archive
				result := uploadArchive(r, &sasURLs)
				if !result {
					// The command has already printed the error
					return
				}
			} else {
				// Get a buffer reader
				file, err := os.Open(path)
				if err != nil {
					utils.ExitWithError(utils.ErrorApp, "Error while reading file", err)
					return
				}

				// Upload the archive
				result := uploadArchive(file, &sasURLs)
				if !result {
					// The command has already printed the error
					return
				}
			}
		},
	}
	uploadCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&app, "app", "a", "", "app bundle name, with no extension (required)")
	c.MarkFlagRequired("app")
	c.Flags().StringVarP(&path, "path", "f", "", "path to local file or folder to bundle (required)")
	c.MarkFlagRequired("path")
	c.Flags().BoolVarP(&noSignature, "no-signature", "", false, "do not cryptographically sign the app's bundle")

	// Add shared flags
	addSharedFlags(c)
}
