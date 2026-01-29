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
	"net/http"
	"strings"

	"github.com/casbin/caswaf/object"
)

// addSecureFlagsToCookies adds Secure, HttpOnly, and SameSite flags to Set-Cookie headers
// based on the site's configuration
func addSecureFlagsToCookies(resp *http.Response, site *object.Site) error {
	if resp == nil || site == nil {
		return nil
	}

	// Get all Set-Cookie headers
	setCookieHeaders := resp.Header["Set-Cookie"]
	if len(setCookieHeaders) == 0 {
		return nil
	}

	// Process each cookie
	modifiedCookies := make([]string, 0, len(setCookieHeaders))
	for _, cookie := range setCookieHeaders {
		modifiedCookie := cookie

		// Add Secure flag if enabled and not already present
		if site.EnableCookieSecure && !strings.Contains(strings.ToLower(cookie), "secure") {
			modifiedCookie = modifiedCookie + "; Secure"
		}

		// Add HttpOnly flag if enabled and not already present
		if site.EnableCookieHttpOnly && !strings.Contains(strings.ToLower(cookie), "httponly") {
			modifiedCookie = modifiedCookie + "; HttpOnly"
		}

		// Add SameSite flag if enabled and not already present
		if site.EnableCookieSameSite && !strings.Contains(strings.ToLower(cookie), "samesite") {
			modifiedCookie = modifiedCookie + "; SameSite=Lax"
		}

		modifiedCookies = append(modifiedCookies, modifiedCookie)
	}

	// Replace the Set-Cookie headers with modified ones
	resp.Header["Set-Cookie"] = modifiedCookies

	return nil
}
