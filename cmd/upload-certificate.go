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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pkcs12 "software.sslmate.com/src/go-pkcs12"

	"smpcli/utils"
)

func init() {
	var (
		name        string
		certificate string
		key         string
		dhparams    string
	)

	// This function gets a client authenticated with Azure Key Vault
	var getKeyVault = func() *keyvault.BaseClient {
		vaultName := viper.GetString("AzureKeyVault")
		if len(vaultName) < 1 {
			fmt.Println("[Error]\nConfiguration variable `AzureKeyVault` must be set before uploading a certificate.")
			return nil
		}

		// Create a new client
		akvClient := keyvault.New()

		// Authorize from the Azure CLI
		authorizer, err := auth.NewAuthorizerFromCLI()
		if err != nil {
			fmt.Println("[Fatal error]\nError while authorizing the Azure Key Vault client:", err)
			return nil
		}
		akvClient.Authorizer = authorizer

		return &akvClient
	}

	// This function uploads the PFX certificate to Azure Key Vault
	var uploadCertificate = func(pfx []byte, akvClient *keyvault.BaseClient) bool {
		// Base URL for the vault
		vaultName := viper.GetString("AzureKeyVault")
		vaultBaseURL := fmt.Sprintf("https://%s.%s", vaultName, azure.PublicCloud.KeyVaultDNSSuffix)

		// Convert certificate to base64
		pfxB64 := base64.StdEncoding.EncodeToString(pfx)

		// Store the certificate
		ctx := context.Background()
		result, err := akvClient.ImportCertificate(ctx, vaultBaseURL, name, keyvault.CertificateImportParameters{
			Base64EncodedCertificate: &pfxB64,
			Password:                 nil,
		})
		if err != nil {
			fmt.Println("[Fatal error]\nError while storing certificate in Azure Key Vault:", err)
			return false
		}
		fmt.Printf("Stored %s\n", *result.ID)

		return true
	}

	// This function loads the certificate
	var loadCertificate = func(file string) (*x509.Certificate, error) {
		// Load certificate from disk
		dataPEM, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		// Get the certificate
		block, _ := pem.Decode(dataPEM)
		if block == nil || block.Type != "CERTIFICATE" {
			return nil, errors.New("Cannot decode PEM block containing certificate")
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		return cert, nil
	}

	// This function loads the private key
	var loadPrivateKey = func(file string) (*rsa.PrivateKey, error) {
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

	// This function creates a PCKS12-encoded file (PFX) with the certificates
	var createPFX = func() []byte {
		// Load the certificate and key
		crt, err := loadCertificate(certificate)
		if err != nil {
			fmt.Println("[Fatal error]\nCannot load certificate:", err)
			return nil
		}
		prv, err := loadPrivateKey(key)
		if err != nil {
			fmt.Println("[Fatal error]\nCannot load private key:", err)
			return nil
		}

		// Crete the PCKS12 bag
		pcksData, err := pkcs12.Encode(rand.Reader, prv, crt, nil, "")
		if err != nil {
			fmt.Println("[Fatal error]\nCannot create PKCS12 bag:", err)
			return nil
		}

		return pcksData
	}

	// This function uploads the dhparams file to Azure Storage
	var uploadDhparams = func() bool {
		// Get config
		storageAccount := viper.GetString("AzureStorageAccount")
		storageKey := viper.GetString("AzureStorageKey")
		storageContainer := viper.GetString("AzureStorageContainer")
		if len(storageAccount) < 1 || len(storageKey) < 1 || len(storageContainer) < 1 {
			fmt.Println("[Error]\nConfiguration variables `AzureStorageAccount`, `AzureStorageKey` and `AzureStorageContainer` must be set before uploading a certificate.")
			return false
		}

		// Stream to dhparams file
		file, err := os.Open(dhparams)
		if err != nil {
			fmt.Println("[Fatal error]\nError while opening dhparams file:", err)
			return false
		}

		// URL to upload to
		dst := fmt.Sprintf("https://%s.blob.core.windows.net/%s/dhparams/%s.pem", storageAccount, storageContainer, name)
		u, err := url.Parse(dst)
		if err != nil {
			fmt.Println("[Fatal error]\nError while building dhparams URL:", err)
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

		// Upload the app's file
		blockBlobURL := azblob.NewBlockBlobURL(*u, pipeline)
		_, err = azblob.UploadStreamToBlockBlob(ctx, file, blockBlobURL, azblob.UploadStreamToBlockBlobOptions{
			BufferSize: 3 * 1024 * 1024,
			MaxBuffers: 2,
		})
		if err != nil {
			fmt.Println("[Fatal error]\nError while uploading file:", err)
			return false
		}
		fmt.Printf("Uploaded %s\n", dst)

		return true
	}

	// This function returns true if the file exists and it's a regular file
	var checkFile = func(path string) bool {
		// Check if the path exists
		isFile, err := utils.IsRegularFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("[Error]\nFile not found:", path)
				return false
			}
			fmt.Println("[Fatal error]\nError while reading filesystem:", err)
			return false
		}
		if !isFile {
			fmt.Println("[Error]\nFile not found:", path)
			return false
		}
		return true
	}

	c := &cobra.Command{
		Use:   "certificate",
		Short: "Upload a certificate",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if all files exist
			if !checkFile(certificate) {
				return
			}
			if !checkFile(key) {
				return
			}
			if !checkFile(dhparams) {
				return
			}

			// Convert the certificate and key to PCKS12
			pfx := createPFX()
			if pfx == nil {
				return
			}

			// Get the Azure Key Vault client
			akvClient := getKeyVault()
			if akvClient == nil {
				return
			}

			// Upload the certificate to Azure Key Vault
			result := uploadCertificate(pfx, akvClient)
			if !result {
				return
			}

			// Upload the dhparams file to Azure Storage
			result = uploadDhparams()
			if !result {
				return
			}
		},
	}
	uploadCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&name, "name", "n", "", "Certificate name (required)")
	c.MarkFlagRequired("name")
	c.Flags().StringVarP(&certificate, "certificate", "c", "", "Certificate file (required)")
	c.MarkFlagRequired("certificate")
	c.Flags().StringVarP(&key, "key", "k", "", "Private key (required)")
	c.MarkFlagRequired("key")
	c.Flags().StringVarP(&dhparams, "dhparams", "d", "", "DH Parameters file (required)")
	c.MarkFlagRequired("dhparams")
}
