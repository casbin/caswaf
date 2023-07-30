// Copyright 2023 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/casbin/caswaf/object"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

func getSigninUrl(casdoorClient *casdoorsdk.Client, callbackUrl string, originalPath string) string {
	scope := "read"
	return fmt.Sprintf("%s/login/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s",
		casdoorClient.Endpoint, casdoorClient.ClientId, url.QueryEscape(callbackUrl), scope, url.QueryEscape(originalPath))
}

func redirectToCasdoor(casdoorClient *casdoorsdk.Client, w http.ResponseWriter, r *http.Request) {
	scheme := r.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}

	callbackUrl := fmt.Sprintf("%s://%s/callback", scheme, r.Host)
	originalPath := r.RequestURI
	signinUrl := getSigninUrl(casdoorClient, callbackUrl, originalPath)
	http.Redirect(w, r, signinUrl, http.StatusFound)
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	site := object.GetSiteByDomain(r.Host)
	if site == nil {
		fmt.Fprintf(w, "CasWAF error: site not found for host: %s", r.Host)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" {
		fmt.Fprint(w, "CasWAF error: the code should not be empty")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if state == "" {
		fmt.Fprint(w, "CasWAF error: the state should not be empty")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	casdoorClient := casdoorsdk.NewClient(site.CasdoorEndpoint, site.CasdoorClientId, site.CasdoorClientSecret, site.CasdoorCertificate, site.CasdoorOrganization, site.CasdoorApplication)
	token, err := casdoorClient.GetOAuthToken(code, state)
	if err != nil {
		fmt.Fprintf(w, "CasWAF error: casdoorClient.GetOAuthToken() error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "casdoor_access_token",
		Value: token.AccessToken,
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	originalPath := state
	http.Redirect(w, r, originalPath, http.StatusFound)
}
