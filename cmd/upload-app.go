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

	"github.com/ItalyPaleAle/smpcli/utils"
)

func init() {
	var (
		app         string
		version     string
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

	// This function uploads the tar.bz2 archive to Azure Blob Storage
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
		Long:  ``,
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

			// Request body for getting the SAS token for Azure Storage from the node
			reqBody := &uploadAuthRequestModel{
				Name:    app,
				Version: version,
			}
			buf := new(bytes.Buffer)
			err = json.NewEncoder(buf).Encode(reqBody)
			if err != nil {
				utils.ExitWithError(utils.ErrorNode, "Error while encoding to JSON", err)
				return
			}

			// Invoke the /uploadauth endpoing and request the SAS token from the node
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

			// Check if the path is already a tar.bz2 archive
			pathLc := strings.ToLower(path)
			if strings.HasSuffix(pathLc, ".tar.bz2") {
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
			} else {
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
			}
		},
	}
	uploadCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&app, "app", "a", "", "App's bundle name (required)")
	c.MarkFlagRequired("app")
	c.Flags().StringVarP(&version, "version", "v", "", "App's bundle version (required)")
	c.MarkFlagRequired("version")
	c.Flags().StringVarP(&path, "path", "f", "", "Path to local file or folder to bundle")
	c.Flags().BoolVarP(&noSignature, "no-signature", "", false, "do not cryptographically sign the app's bundle")

	// Add shared flags
	addSharedFlags(c)
}
