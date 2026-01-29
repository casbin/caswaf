// Copyright 2024 The casbin Authors. All Rights Reserved.
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
	"testing"

	"github.com/casbin/caswaf/object"
)

func TestSetHSTSHeader(t *testing.T) {
	tests := []struct {
		name                  string
		site                  *object.Site
		isHTTPS               bool
		expectedHeader        string
		expectHeaderSet       bool
	}{
		{
			name: "HSTS enabled with includeSubDomains over HTTPS",
			site: &object.Site{
				EnableHSTS:            true,
				HSTSMaxAge:            31536000,
				HSTSIncludeSubDomains: true,
			},
			isHTTPS:         true,
			expectedHeader:  "max-age=31536000; includeSubDomains",
			expectHeaderSet: true,
		},
		{
			name: "HSTS enabled without includeSubDomains over HTTPS",
			site: &object.Site{
				EnableHSTS:            true,
				HSTSMaxAge:            31536000,
				HSTSIncludeSubDomains: false,
			},
			isHTTPS:         true,
			expectedHeader:  "max-age=31536000",
			expectHeaderSet: true,
		},
		{
			name: "HSTS enabled over HTTP - should not set header",
			site: &object.Site{
				EnableHSTS:            true,
				HSTSMaxAge:            31536000,
				HSTSIncludeSubDomains: true,
			},
			isHTTPS:         false,
			expectedHeader:  "",
			expectHeaderSet: false,
		},
		{
			name: "HSTS disabled over HTTPS",
			site: &object.Site{
				EnableHSTS:            false,
				HSTSMaxAge:            31536000,
				HSTSIncludeSubDomains: true,
			},
			isHTTPS:         true,
			expectedHeader:  "",
			expectHeaderSet: false,
		},
		{
			name: "HSTS enabled with zero max-age - should not set header",
			site: &object.Site{
				EnableHSTS:            true,
				HSTSMaxAge:            0,
				HSTSIncludeSubDomains: true,
			},
			isHTTPS:         true,
			expectedHeader:  "",
			expectHeaderSet: false,
		},
		{
			name: "HSTS enabled with custom max-age over HTTPS",
			site: &object.Site{
				EnableHSTS:            true,
				HSTSMaxAge:            86400,
				HSTSIncludeSubDomains: false,
			},
			isHTTPS:         true,
			expectedHeader:  "max-age=86400",
			expectHeaderSet: true,
		},
		{
			name:            "Nil site - should not set header",
			site:            nil,
			isHTTPS:         true,
			expectedHeader:  "",
			expectHeaderSet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			setHSTSHeader(tt.site, header, tt.isHTTPS)

			hstsHeader := header.Get("Strict-Transport-Security")
			if tt.expectHeaderSet {
				if hstsHeader != tt.expectedHeader {
					t.Errorf("Expected HSTS header to be %q, got %q", tt.expectedHeader, hstsHeader)
				}
			} else {
				if hstsHeader != "" {
					t.Errorf("Expected HSTS header to not be set, got %q", hstsHeader)
				}
			}
		})
	}
}
