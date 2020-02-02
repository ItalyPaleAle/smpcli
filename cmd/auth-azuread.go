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
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ItalyPaleAle/smpcli/utils"
)

func init() {
	c := &cobra.Command{
		Use:   "azuread",
		Short: "Authenticate using an Azure AD account",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			baseURL, client := getURLClient()

			// Invoke the /info endpoint to see what's the authentication method
			var rInfo infoResponseModel
			err := utils.RequestJSON(utils.RequestOpts{
				Client: client,
				Target: &rInfo,
				URL:    baseURL + "/info",
			})
			if err != nil {
				fmt.Println("[Fatal error]\nRequest failed:", err)
				return
			}

			// Ensure the node supports pre-shared key authentication
			if !utils.SliceContainsString(rInfo.AuthMethods, "azureAD") || rInfo.AzureAD == nil {
				fmt.Println("[Fatal error]\nThis node does not support authenticating with an Azure AD account")
				return
			}

			// Redirect users to the authentication URL
			authorizeURL := fmt.Sprintf("%s?client_id=%s&response_type=code&redirect_uri=%s&response_mode=query&domain_hint=organizations&scope=openid+offline_access", rInfo.AzureAD.AuthorizeURL, rInfo.AzureAD.ClientID, url.QueryEscape("http://localhost:3993"))
			utils.LaunchBrowser(authorizeURL)

			// Start a web server to listen to authorization codes
			authCode := ""
			ctx, ctxCancel := context.WithCancel(context.Background())
			defer ctxCancel()
			mux := http.NewServeMux()
			server := &http.Server{
				Addr:           "127.0.0.1:3993",
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
				Handler:        mux,
			}
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				// Ensure we have the code in the response
				query := r.URL.Query()
				if query != nil && query.Get("code") != "" {
					authCode = query.Get("code")
					fmt.Fprintf(w, "Authenticated with Azure AD. You can close this window.")
					ctxCancel()
				} else {
					fmt.Fprintf(w, "Error: response did not contain an authorization code")
				}
			})
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			select {
			// Shutdown the server when the context is canceled
			case <-ctx.Done():
				server.Shutdown(ctx)
			}

			// Exchange the authorization code for a token
			body := url.Values{}
			// No client_secret because this is a client-side app
			body.Set("client_id", rInfo.AzureAD.ClientID)
			body.Set("code", authCode)
			body.Set("grant_type", "authorization_code")
			body.Set("redirect_uri", "http://localhost:3993")
			body.Set("scope", "openid offline_access")

			// Request
			var rToken struct {
				ExpiresIn    int    `json:"expires_in"`
				IDToken      string `json:"id_token"`
				RefreshToken string `json:"refresh_token"`
			}
			err = utils.RequestJSON(utils.RequestOpts{
				Body:            strings.NewReader(body.Encode()),
				BodyContentType: "application/x-www-form-urlencoded",
				Method:          utils.RequestPOST,
				Target:          &rToken,
				URL:             rInfo.AzureAD.TokenURL,
			})
			if err != nil {
				fmt.Println("[Fatal error]\nRequest failed:", err)
				return
			}

			if rToken.IDToken == "" || rToken.RefreshToken == "" {
				fmt.Println("[Fatal error]\nResponse did not contain an id_token or a refresh_token")
				return
			}

			// Test the auth token by requesting the node's state, invoking the /state endpoint
			// We're not requesting anything from the response
			var rState struct{}
			err = utils.RequestJSON(utils.RequestOpts{
				Authorization: rToken.IDToken,
				Client:        client,
				Target:        &rState,
				URL:           baseURL + "/state",
			})
			if err != nil {
				// Check if the error is a 401
				if strings.HasPrefix(err.Error(), "invalid response status code: 401") {
					fmt.Println("[Error]\nInvalid pre-shared key")
				} else {
					fmt.Println("[Server error]\n", err)
				}
				return
			}

			// Store the key in the node store
			/*if err := nodeStore.StoreSharedKey(optAddress, sharedKey); err != nil {
				fmt.Println("[Fatal error]\nError while storing the pre-shared key:", err)
				return
			}*/

			fmt.Println("Success! You're authenticated")
		},
	}

	authCmd.AddCommand(c)

	// Add shared flags
	addSharedFlags(c)
}
