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

package proxy

import (
	"net/http"
	"testing"

	"github.com/beego/beego"
)

func TestGetCasdoorHttpClient_WithInsecureSkipVerify(t *testing.T) {
	// Set config for insecure skip verify
	beego.AppConfig.Set("casdoorInsecureSkipVerify", "true")

	client := getCasdoorHttpClient()

	if client == nil {
		t.Fatal("Expected non-nil HTTP client")
	}

	// Check if transport is configured with InsecureSkipVerify
	if tr, ok := client.Transport.(*http.Transport); ok {
		if tr.TLSClientConfig == nil {
			t.Fatal("Expected TLSClientConfig to be set")
		}
		if !tr.TLSClientConfig.InsecureSkipVerify {
			t.Fatal("Expected InsecureSkipVerify to be true")
		}
	} else {
		t.Fatal("Expected http.Transport")
	}
}

func TestGetCasdoorHttpClient_WithoutInsecureSkipVerify(t *testing.T) {
	// Set config for secure mode (default)
	beego.AppConfig.Set("casdoorInsecureSkipVerify", "false")

	client := getCasdoorHttpClient()

	if client == nil {
		t.Fatal("Expected non-nil HTTP client")
	}

	// Check if transport uses default settings (no custom TLS config or InsecureSkipVerify is false)
	if tr, ok := client.Transport.(*http.Transport); ok {
		if tr.TLSClientConfig != nil && tr.TLSClientConfig.InsecureSkipVerify {
			t.Fatal("Expected InsecureSkipVerify to be false in secure mode")
		}
	}
}

func TestGetCasdoorHttpClient_DefaultBehavior(t *testing.T) {
	// Remove the config key to test default behavior
	beego.AppConfig.Set("casdoorInsecureSkipVerify", "")

	client := getCasdoorHttpClient()

	if client == nil {
		t.Fatal("Expected non-nil HTTP client")
	}

	// Default should be secure (InsecureSkipVerify = false)
	if tr, ok := client.Transport.(*http.Transport); ok {
		if tr.TLSClientConfig != nil && tr.TLSClientConfig.InsecureSkipVerify {
			t.Fatal("Expected default behavior to be secure (InsecureSkipVerify = false)")
		}
	}
}

func TestInitHttpClient(t *testing.T) {
	beego.AppConfig.Set("casdoorInsecureSkipVerify", "true")
	beego.AppConfig.Set("httpProxy", "")

	InitHttpClient()

	if DefaultHttpClient == nil {
		t.Fatal("Expected DefaultHttpClient to be initialized")
	}

	if ProxyHttpClient == nil {
		t.Fatal("Expected ProxyHttpClient to be initialized")
	}

	if CasdoorHttpClient == nil {
		t.Fatal("Expected CasdoorHttpClient to be initialized")
	}

	// Verify CasdoorHttpClient has insecure skip verify enabled
	if tr, ok := CasdoorHttpClient.Transport.(*http.Transport); ok {
		if tr.TLSClientConfig == nil || !tr.TLSClientConfig.InsecureSkipVerify {
			t.Fatal("Expected CasdoorHttpClient to have InsecureSkipVerify enabled")
		}
	}
}

func TestCasdoorHttpClient_TLSConfig(t *testing.T) {
	testCases := []struct {
		name                   string
		configValue            string
		expectInsecureSkipVerify bool
	}{
		{"Insecure mode enabled", "true", true},
		{"Insecure mode disabled", "false", false},
		{"Default (empty string)", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			beego.AppConfig.Set("casdoorInsecureSkipVerify", tc.configValue)

			client := getCasdoorHttpClient()

			if tr, ok := client.Transport.(*http.Transport); ok {
				if tc.expectInsecureSkipVerify {
					if tr.TLSClientConfig == nil || !tr.TLSClientConfig.InsecureSkipVerify {
						t.Errorf("Expected InsecureSkipVerify to be true for config value: %s", tc.configValue)
					}
				} else {
					if tr.TLSClientConfig != nil && tr.TLSClientConfig.InsecureSkipVerify {
						t.Errorf("Expected InsecureSkipVerify to be false for config value: %s", tc.configValue)
					}
				}
			} else if tc.expectInsecureSkipVerify {
				t.Error("Expected http.Transport for insecure mode")
			}
		})
	}
}
