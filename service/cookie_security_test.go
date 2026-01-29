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
	"testing"

	"github.com/casbin/caswaf/object"
)

func TestAddSecureFlagsToCookies(t *testing.T) {
	tests := []struct {
		name               string
		inputCookies       []string
		enableSecure       bool
		enableHttpOnly     bool
		enableSameSite     bool
		expectedContains   []string
		unexpectedContains []string
	}{
		{
			name:             "Add Secure flag only",
			inputCookies:     []string{"sessionid=abc123; Path=/"},
			enableSecure:     true,
			enableHttpOnly:   false,
			enableSameSite:   false,
			expectedContains: []string{"Secure"},
		},
		{
			name:             "Add HttpOnly flag only",
			inputCookies:     []string{"sessionid=abc123; Path=/"},
			enableSecure:     false,
			enableHttpOnly:   true,
			enableSameSite:   false,
			expectedContains: []string{"HttpOnly"},
		},
		{
			name:             "Add SameSite flag only",
			inputCookies:     []string{"sessionid=abc123; Path=/"},
			enableSecure:     false,
			enableHttpOnly:   false,
			enableSameSite:   true,
			expectedContains: []string{"SameSite=Lax"},
		},
		{
			name:             "Add all flags",
			inputCookies:     []string{"sessionid=abc123; Path=/"},
			enableSecure:     true,
			enableHttpOnly:   true,
			enableSameSite:   true,
			expectedContains: []string{"Secure", "HttpOnly", "SameSite=Lax"},
		},
		{
			name:               "Don't add existing Secure flag",
			inputCookies:       []string{"sessionid=abc123; Path=/; Secure"},
			enableSecure:       true,
			enableHttpOnly:     false,
			enableSameSite:     false,
			expectedContains:   []string{"Secure"},
			unexpectedContains: []string{"Secure; Secure"},
		},
		{
			name:               "Don't add existing HttpOnly flag",
			inputCookies:       []string{"sessionid=abc123; Path=/; HttpOnly"},
			enableSecure:       false,
			enableHttpOnly:     true,
			enableSameSite:     false,
			expectedContains:   []string{"HttpOnly"},
			unexpectedContains: []string{"HttpOnly; HttpOnly"},
		},
		{
			name:               "Don't add existing SameSite flag",
			inputCookies:       []string{"sessionid=abc123; Path=/; SameSite=Strict"},
			enableSecure:       false,
			enableHttpOnly:     false,
			enableSameSite:     true,
			expectedContains:   []string{"SameSite=Strict"},
			unexpectedContains: []string{"SameSite=Lax"},
		},
		{
			name:             "Handle multiple cookies",
			inputCookies:     []string{"sessionid=abc123; Path=/", "token=xyz789; Path=/"},
			enableSecure:     true,
			enableHttpOnly:   true,
			enableSameSite:   false,
			expectedContains: []string{"Secure", "HttpOnly"},
		},
		{
			name:             "Handle PHPSESSID cookie",
			inputCookies:     []string{"PHPSESSID=abc123; Path=/"},
			enableSecure:     true,
			enableHttpOnly:   true,
			enableSameSite:   true,
			expectedContains: []string{"Secure", "HttpOnly", "SameSite=Lax"},
		},
		{
			name:             "Handle loginToken cookie",
			inputCookies:     []string{"loginToken=xyz789; Path=/"},
			enableSecure:     true,
			enableHttpOnly:   true,
			enableSameSite:   true,
			expectedContains: []string{"Secure", "HttpOnly", "SameSite=Lax"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock response
			resp := &http.Response{
				Header: http.Header{},
			}
			resp.Header["Set-Cookie"] = tt.inputCookies

			// Create mock site with configuration
			site := &object.Site{
				EnableCookieSecure:   tt.enableSecure,
				EnableCookieHttpOnly: tt.enableHttpOnly,
				EnableCookieSameSite: tt.enableSameSite,
			}

			// Call the function
			err := addSecureFlagsToCookies(resp, site)
			if err != nil {
				t.Errorf("addSecureFlagsToCookies() returned error: %v", err)
				return
			}

			// Verify results
			modifiedCookies := resp.Header["Set-Cookie"]
			if len(modifiedCookies) != len(tt.inputCookies) {
				t.Errorf("Expected %d cookies, got %d", len(tt.inputCookies), len(modifiedCookies))
				return
			}

			// Check expected strings are present
			for _, cookie := range modifiedCookies {
				for _, expected := range tt.expectedContains {
					if !strings.Contains(cookie, expected) {
						t.Errorf("Expected cookie to contain '%s', but got: %s", expected, cookie)
					}
				}
				for _, unexpected := range tt.unexpectedContains {
					if strings.Contains(cookie, unexpected) {
						t.Errorf("Did not expect cookie to contain '%s', but got: %s", unexpected, cookie)
					}
				}
			}
		})
	}
}

func TestAddSecureFlagsToCookies_NilInputs(t *testing.T) {
	tests := []struct {
		name string
		resp *http.Response
		site *object.Site
	}{
		{
			name: "Nil response",
			resp: nil,
			site: &object.Site{EnableCookieSecure: true},
		},
		{
			name: "Nil site",
			resp: &http.Response{Header: http.Header{}},
			site: nil,
		},
		{
			name: "Both nil",
			resp: nil,
			site: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := addSecureFlagsToCookies(tt.resp, tt.site)
			if err != nil {
				t.Errorf("addSecureFlagsToCookies() should handle nil inputs gracefully, got error: %v", err)
			}
		})
	}
}

func TestAddSecureFlagsToCookies_NoCookies(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}
	site := &object.Site{
		EnableCookieSecure:   true,
		EnableCookieHttpOnly: true,
		EnableCookieSameSite: true,
	}

	err := addSecureFlagsToCookies(resp, site)
	if err != nil {
		t.Errorf("addSecureFlagsToCookies() should handle response with no cookies gracefully, got error: %v", err)
	}

	if len(resp.Header["Set-Cookie"]) != 0 {
		t.Errorf("Expected no cookies, got %d", len(resp.Header["Set-Cookie"]))
	}
}
