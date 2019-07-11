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
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"smpcli/utils"
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
	var uploadArchive = func(file io.Reader) bool {
		// Get config
		storageAccount := viper.GetString("AzureStorageAccount")
		storageKey := viper.GetString("AzureStorageKey")
		storageContainer := viper.GetString("AzureStorageContainer")
		signingKeyFile := viper.GetString("SigningKey")
		if len(storageAccount) < 1 || len(storageKey) < 1 || len(storageContainer) < 1 {
			fmt.Printf("[Error]\nConfiguration variables `AzureStorageAccount`, `AzureStorageKey` and `AzureStorageContainer` must be set before uploading a file.")
			return false
		}

		// Check if the key exists
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

		// URL to upload to
		dstApp := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s-%s.tar.bz2", storageAccount, storageContainer, app, version)
		uApp, err := url.Parse(dstApp)
		if err != nil {
			fmt.Println("[Fatal error]\nError while building app URL:", err)
			return false
		}

		// URL where to upload the signature to
		dstSig := dstApp + ".sig"
		uSig, err := url.Parse(dstSig)
		if err != nil {
			fmt.Println("[Fatal error]\nError while building signature URL:", err)
			return false
		}

		// Uploader client
		credential, err := azblob.NewSharedKeyCredential(storageAccount, storageKey)
		if err != nil {
			fmt.Println("[Fatal error]\nError while getting credentials:", err)
			return false
		}
		pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{
			Retry: azblob.RetryOptions{
				MaxTries: 3,
			},
		})
		ctx := context.Background()

		// The stream is split between two readers: one for the hashing, one for writing the stream to disk
		h := sha256.New()
		tee := io.TeeReader(file, h)

		// Upload the app's file
		// This also makes the stream proceed so the hash is calculated
		blockBlobURL := azblob.NewBlockBlobURL(*uApp, pipeline)
		_, err = azblob.UploadStreamToBlockBlob(ctx, tee, blockBlobURL, azblob.UploadStreamToBlockBlobOptions{
			BufferSize: 3 * 1024 * 1024,
			MaxBuffers: 2,
		})
		if err != nil {
			fmt.Println("[Fatal error]\nError while uploading file:", err)
			return false
		}
		fmt.Printf("Uploaded %s\n", dstApp)

		// Calculate the SHA256 hash
		hashed := h.Sum(nil)

		// Calculate the digital signature
		rng := rand.Reader
		signatureRaw, err := rsa.SignPKCS1v15(rng, signingKey, crypto.SHA256, hashed[:])
		if err != nil {
			fmt.Println("[Fatal error]\nCannot calculate signature:", err)
			return false
		}

		// Convert the signature to base64
		signature := base64.StdEncoding.EncodeToString(signatureRaw)

		// Upload the signature
		blockBlobURL = azblob.NewBlockBlobURL(*uSig, pipeline)
		_, err = blockBlobURL.Upload(ctx, strings.NewReader(signature), azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
		if err != nil {
			fmt.Println("[Fatal error]\nError while uploading signature:", err)
			return false
		}
		fmt.Printf("Uploaded %s\n", dstSig)

		return true
	}

	c := &cobra.Command{
		Use:   "app",
		Short: "Upload an app or bundle",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
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
				result := uploadArchive(file)
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
				result := uploadArchive(r)
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
}
