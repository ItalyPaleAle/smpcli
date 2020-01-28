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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ItalyPaleAle/smpcli/utils"
)

func init() {
	var (
		app     string
		version string
		path    string
	)

	// This function loads the private key
	var loadSigningKey = func(file string) (*rsa.PrivateKey, error) {
		// Load key from disk
		dataPEM, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		// Parse the key
		block, _ := pem.Decode(dataPEM)
		if block == nil || block.Type != "RSA PRIVATE KEY" {
			return nil, errors.New("Cannot decode PEM block containing private key")
		}
		prv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		return prv, nil
	}

	// This function uploads the tar.bz2 archive to Azure Blob Storage
	// Returns true in case of success, and false if there's an error
	var uploadArchive = func(file io.Reader, sasURLs *uploadAuthResponseModel) bool {
		// Check if the key exists
		signingKeyFile := viper.GetString("SigningKey")
		if len(signingKeyFile) < 1 {
			fmt.Printf("[Error]\nNo private signing key set in the `SigningKey` configuration variable")
			return false
		}
		exists, err := utils.PathExists(signingKeyFile)
		if err != nil {
			fmt.Println("[Fatal error]\nError while reading filesystem:", err)
			return false
		}
		if !exists {
			fmt.Println("[Error]\nPrivate signing key not found:", signingKeyFile)
			return false
		}

		// Load the signing key
		signingKey, err := loadSigningKey(signingKeyFile)

		// URL to upload the archive to
		uApp, err := url.Parse(sasURLs.ArchiveURL)
		if err != nil {
			fmt.Println("[Fatal error]\nError while building app URL:", err)
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
				fmt.Println("[Fatal error]\nNetwork error while uploading the archive:\n", err)
			} else {
				fmt.Println("[Fatal error]\nAzure Storage error failed while uploading the archive:\n", stgErr.Response().Status)
			}
			return false
		}
		fmt.Println("Uploaded app's bundle")

		// Calculate the SHA256 hash
		hashed := h.Sum(nil)

		// Calculate the digital signature
		rng := rand.Reader
		signatureRaw, err := rsa.SignPKCS1v15(rng, signingKey, crypto.SHA256, hashed[:])
		if err != nil {
			fmt.Println("[Fatal error]\nCannot calculate signature:\n", err)
			return false
		}

		// Convert the signature to base64
		signature := base64.StdEncoding.EncodeToString(signatureRaw)

		// Add the signature as metadata
		metadata := azblob.Metadata{}
		metadata["signature"] = signature
		_, err = blockBlobURL.SetMetadata(ctx, metadata, azblob.BlobAccessConditions{})
		if err != nil {
			fmt.Println("[Fatal error]\nCannot update blob's metadata:\n", err)
			return false
		}

		return true
	}

	c := &cobra.Command{
		Use:   "app",
		Short: "Upload an app or bundle",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()

			// Get the shared key
			sharedKey, found, err := nodeStore.GetSharedKey(optAddress)
			if err != nil {
				fmt.Println("[Fatal error]\nError while reading node store:", err)
				return
			}
			if !found {
				fmt.Printf("[Error]\nNo authentication data for the domain %s; please make sure you've executed the 'auth' command.\n", optAddress)
				return
			}

			// Check if the path exists
			exists, err := utils.PathExists(path)
			if err != nil {
				fmt.Println("[Fatal error]\nError while reading filesystem:", err)
				return
			}
			if !exists {
				fmt.Println("[Error]\nFile/folder not found:", path)
				return
			}

			// Request body for getting the SAS token for Azure Storage from the node
			reqBody := &uploadAuthRequestModel{
				Name:    app,
				Version: version,
			}
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(reqBody)

			// Request the SAS token from the node
			req, err := http.NewRequest("POST", baseURL+"/uploadauth", buf)
			if err != nil {
				fmt.Println("[Fatal error]\nCould not build the request:", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", sharedKey)
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("[Fatal error]\nRequest failed:", err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				b, _ := ioutil.ReadAll(resp.Body)
				fmt.Printf("[Server error]\n%d: %s\n", resp.StatusCode, string(b))
				return
			}

			// Parse the response
			var sasURLs uploadAuthResponseModel
			if err := json.NewDecoder(resp.Body).Decode(&sasURLs); err != nil {
				fmt.Println("[Fatal error]\nInvalid JSON response:", err)
				return
			}

			// Check if the path is already a tar.bz2 archive
			pathLc := strings.ToLower(path)
			if strings.HasSuffix(pathLc, ".tar.bz2") {
				// Get a buffer reader
				file, err := os.Open(path)
				if err != nil {
					fmt.Println("[Fatal error]\nError while reading file:", err)
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
						fmt.Println("[Fatal error]\nError while creating tar.bz2 archive:", err)
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

	// Add shared flags
	addSharedFlags(c)
}
