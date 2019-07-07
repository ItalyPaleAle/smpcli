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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
		domain  string
		app     string
		version string
		path    string
	)

	// This funtion sends the request to the node to deploy the app
	// Returns true in case of success, and false if there's an error
	var sendRequest = func(sharedKey string) bool {
		baseURL, client := getURLClient()

		// Request body
		reqBody := &deployRequestModel{
			App:     app,
			Version: version,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(reqBody)

		// Invoke the /site endpoint and add the site
		req, err := http.NewRequest("POST", baseURL+"/site/"+domain+"/deploy", buf)
		if err != nil {
			fmt.Println("[Fatal error]\nCould not build the request:", err)
			return false
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", sharedKey)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("[Fatal error]\nRequest failed:", err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			b, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("[Server error]\n%d: %s\n", resp.StatusCode, string(b))
			return false
		}

		// Parse the response
		var r deployResponseModel
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			fmt.Println("[Fatal error]\nInvalid JSON response:", err)
			return false
		}

		fmt.Println(deployResponseModelFormat(&r))
		return true
	}

	// This function uploads the tar.bz2 archive to Azure Blob Storage
	// Returns true in case of success, and false if there's an error
	var uploadArchive = func(path string) bool {
		// Get variables
		storageAccount := viper.GetString("AzureStorageAccount")
		storageKey := viper.GetString("AzureStorageKey")
		storageContainer := viper.GetString("AzureStorageContainer")
		if len(storageAccount) < 1 || len(storageKey) < 1 || len(storageContainer) < 1 {
			fmt.Printf("[Error]\nConfiguration variables `AzureStorageAccount`, `AzureStorageKey` and `AzureStorageContainer` must be set before uploading a file.")
			return false
		}

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
		blockBlobURL := azblob.NewBlockBlobURL(*uApp, azblob.NewPipeline(credential, azblob.PipelineOptions{
			Retry: azblob.RetryOptions{
				MaxTries: 3,
			},
		}))
		ctx := context.Background()

		// Get a buffer reader
		file, err := os.Open(path)
		if err != nil {
			fmt.Println("[Fatal error]\nError while reading file:", err)
			return false
		}

		// Upload the app's file
		_, err = azblob.UploadStreamToBlockBlob(ctx, file, blockBlobURL, azblob.UploadStreamToBlockBlobOptions{
			BufferSize: 3 * 1024 * 1024,
			MaxBuffers: 2,
		})
		if err != nil {
			fmt.Println("[Fatal error]\nError while uploading file:", err)
			return false
		}
		fmt.Printf("Uploaded %s\n", dstApp)

		return true
	}

	c := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an app",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
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

			// Check if we need to upload a file or folder first
			if len(path) > 0 {
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
					// Upload the archive
					result := uploadArchive(path)
					if !result {
						return
					}
				}
			}

			// Returns true if succeeded
			_ = sendRequest(sharedKey)
		},
	}
	rootCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&domain, "domain", "d", "", "Primary domain name")
	c.MarkFlagRequired("domain")
	c.Flags().StringVarP(&app, "app", "a", "", "App's bundle name")
	c.MarkFlagRequired("app")
	c.Flags().StringVarP(&version, "version", "v", "", "App's bundle version")
	c.MarkFlagRequired("version")
	c.Flags().StringVarP(&path, "path", "f", "", "Path to local file or folder to bundle (optional)")

	// Add shared flags
	addSharedFlags(c)
}
