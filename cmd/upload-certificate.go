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
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/spf13/cobra"
	pkcs12 "software.sslmate.com/src/go-pkcs12"

	"github.com/ItalyPaleAle/smpcli/utils"
)

func init() {
	var (
		name        string
		certificate string
		certKey     string
	)

	var keyVaultName string
	var keyVaultURL string

	// This function requests the name of the Azure Key Vault from the node
	var getKeyVaultName = func() error {
		baseURL, client := getURLClient()
		// Get the shared key
		sharedKey, found, err := nodeStore.GetSharedKey(optAddress)
		if err != nil {
			return fmt.Errorf("Error while reading node store: %s", err.Error())
		}
		if !found {
			return fmt.Errorf("No authentication data for the domain %s; please make sure you've executed the 'auth' command.\n", optAddress)
		}

		// Invoke the /keyvaultname endpoint to get the name and URL of the key vault
		req, err := http.NewRequest("GET", baseURL+"/keyvaultname", nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", sharedKey)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errors.New("Invalid response status code")
		}

		// Parse the response
		var r map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return err
		}

		var ok bool
		keyVaultName, ok = r["name"]
		if !ok || keyVaultName == "" {
			return errors.New("Invalid response: empty name")
		}
		keyVaultURL, ok = r["url"]
		if !ok || keyVaultURL == "" {
			return errors.New("Invalid response: empty url")
		}

		return nil
	}

	// This function gets a client authenticated with Azure Key Vault
	var getKeyVault = func() *keyvault.BaseClient {
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
		// Convert certificate to base64
		pfxB64 := base64.StdEncoding.EncodeToString(pfx)

		// Store the certificate
		ctx := context.Background()
		result, err := akvClient.ImportCertificate(ctx, keyVaultURL, name, keyvault.CertificateImportParameters{
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
		prv, err := loadPrivateKey(certKey)
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
			if !checkFile(certKey) {
				return
			}

			// Convert the certificate and key to PCKS12
			pfx := createPFX()
			if pfx == nil {
				return
			}

			// Get the details of the Azure Key Vault
			if err := getKeyVaultName(); err != nil {
				fmt.Println("[Error]:", err)
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
		},
	}
	uploadCmd.AddCommand(c)

	// Flags
	c.Flags().StringVarP(&name, "name", "c", "", "Certificate name (required)")
	c.MarkFlagRequired("name")
	c.Flags().StringVarP(&certificate, "certificate", "f", "", "Certificate file (required)")
	c.MarkFlagRequired("certificate")
	c.Flags().StringVarP(&certKey, "certificate-key", "p", "", "Private key (required)")
	c.MarkFlagRequired("certificate-key")

	// Add shared flags
	addSharedFlags(c)
}
