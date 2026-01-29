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
// based on the site's configuration.
// Note: This feature should only be enabled when the reverse proxy is accessed exclusively via HTTPS,
// as the Secure flag prevents cookies from being sent over HTTP connections.
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
		if site.EnableCookieSecure && !hasSecureFlag(cookie) {
			modifiedCookie = modifiedCookie + "; Secure"
		}

		// Add HttpOnly flag if enabled and not already present
		if site.EnableCookieHttpOnly && !hasHttpOnlyFlag(cookie) {
			modifiedCookie = modifiedCookie + "; HttpOnly"
		}

		// Add SameSite flag if enabled and not already present
		if site.EnableCookieSameSite && !hasSameSiteFlag(cookie) {
			modifiedCookie = modifiedCookie + "; SameSite=Lax"
		}

		modifiedCookies = append(modifiedCookies, modifiedCookie)
	}

	// Replace the Set-Cookie headers with modified ones
	resp.Header["Set-Cookie"] = modifiedCookies

	// Return value is always nil for now, but error signature is kept for future extensibility
	return nil
}

// hasSecureFlag checks if a Set-Cookie header already has the Secure flag
func hasSecureFlag(cookie string) bool {
	return hasCookieAttribute(cookie, "secure")
}

// hasHttpOnlyFlag checks if a Set-Cookie header already has the HttpOnly flag
func hasHttpOnlyFlag(cookie string) bool {
	return hasCookieAttribute(cookie, "httponly")
}

// hasSameSiteFlag checks if a Set-Cookie header already has the SameSite flag
func hasSameSiteFlag(cookie string) bool {
	return hasCookieAttribute(cookie, "samesite")
}

// hasCookieAttribute checks if a cookie attribute exists as a standalone attribute (not part of a value)
// Cookie attributes are separated by semicolons, so we split and check each part
func hasCookieAttribute(cookie string, attribute string) bool {
	// Split cookie by semicolons to get individual attributes
	parts := strings.Split(cookie, ";")
	
	attributeLower := strings.ToLower(attribute)
	
	for _, part := range parts {
		// Trim whitespace from the attribute
		trimmedPart := strings.TrimSpace(part)
		
		// Get the attribute name (before = if present)
		attrName := trimmedPart
		if idx := strings.Index(trimmedPart, "="); idx != -1 {
			attrName = trimmedPart[:idx]
		}
		
		// Compare attribute names case-insensitively
		if strings.ToLower(strings.TrimSpace(attrName)) == attributeLower {
			return true
		}
	}
	
	return false
}
