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
	"os"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/spf13/cobra"
	pkcs12 "software.sslmate.com/src/go-pkcs12"

	"github.com/ItalyPaleAle/stkcli/utils"
)

func init() {
	var (
		name        string
		certificate string
		certKey     string
	)

	// This function uploads the PFX certificate to Azure Key Vault
	var uploadCertificate = func(pfx []byte, akvClient *keyvault.BaseClient) bool {
		// Get the URL of the Key Vault, requesting it from the node
		keyVaultURL, _, _, err := getKeyVaultInfo()
		if err != nil {
			utils.ExitWithError(utils.ErrorApp, "Error while requesting name of Key Vault", err)
			return false
		}

		// Convert certificate to base64
		pfxB64 := base64.StdEncoding.EncodeToString(pfx)

		// Store the certificate
		ctx := context.Background()
		result, err := akvClient.ImportCertificate(ctx, keyVaultURL, name, keyvault.CertificateImportParameters{
			Base64EncodedCertificate: &pfxB64,
			Password:                 nil,
		})
		if err != nil {
			utils.ExitWithError(utils.ErrorApp, "Error while storing certificate in Azure Key Vault", err)
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
		var prv *rsa.PrivateKey
		if block == nil {
			return nil, errors.New("cannot decode PEM block containing private key: empty block")
		}
		if block.Type == "RSA PRIVATE KEY" {
			prv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		} else if block.Type == "PRIVATE KEY" {
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			var ok bool
			prv, ok = key.(*rsa.PrivateKey)
			if !ok {
				return nil, errors.New("cannot decode PEM block containing private key: invalid key algorithm")
			}
		} else {
			return nil, errors.New("cannot decode PEM block containing private key: invalid block type")
		}

		return prv, nil
	}

	// This function creates a PCKS12-encoded file (PFX) with the certificates
	var createPFX = func() []byte {
		// Load the certificate and key
		crt, err := loadCertificate(certificate)
		if err != nil {
			utils.ExitWithError(utils.ErrorApp, "Cannot load certificate", err)
			return nil
		}
		prv, err := loadPrivateKey(certKey)
		if err != nil {
			utils.ExitWithError(utils.ErrorApp, "Cannot load private key", err)
			return nil
		}

		// Crete the PCKS12 bag
		pcksData, err := pkcs12.Encode(rand.Reader, prv, crt, nil, "")
		if err != nil {
			utils.ExitWithError(utils.ErrorApp, "Cannot create PKCS12 bag", err)
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
				utils.ExitWithError(utils.ErrorApp, "File not found: "+path, nil)
				return false
			}
			utils.ExitWithError(utils.ErrorApp, "Error while reading filesystem", err)
			return false
		}
		if !isFile {
			utils.ExitWithError(utils.ErrorApp, "File not found: "+path, nil)
			return false
		}
		return true
	}

	c := &cobra.Command{
		Use:   "certificate",
		Short: "Upload a TLS certificate",
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
