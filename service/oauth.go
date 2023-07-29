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

func redirectToCasdoor(casdoorClient *casdoorsdk.Client, w http.ResponseWriter, r *http.Request) {
	callbackUrl := fmt.Sprintf("%s://%s/callback", r.URL.Scheme, r.Host)
	signinUrl := casdoorClient.GetSigninUrl(callbackUrl)
	w.Header().Set("Location", signinUrl)

	w.WriteHeader(http.StatusFound)
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

	referrerUrl, err := url.Parse(r.Referer())
	if err != nil {
		panic(err)
	}

	targetUrl := fmt.Sprintf("%s://%s", referrerUrl.Scheme, joinPath(site.Domain, referrerUrl.Path))
	w.Header().Set("Location", targetUrl)
	w.WriteHeader(http.StatusFound)
}
