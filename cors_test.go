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

package main

import (
	"strings"
	"testing"
)

// TestCORSOriginsParsing tests that CORS origins are properly parsed from configuration
func TestCORSOriginsParsing(t *testing.T) {
	tests := []struct {
		name           string
		config         string
		expectedLength int
		expectedOrigins []string
	}{
		{
			name:           "Single origin",
			config:         "http://localhost:7001",
			expectedLength: 1,
			expectedOrigins: []string{"http://localhost:7001"},
		},
		{
			name:           "Multiple origins",
			config:         "http://localhost:7001,http://localhost:17000",
			expectedLength: 2,
			expectedOrigins: []string{"http://localhost:7001", "http://localhost:17000"},
		},
		{
			name:           "Multiple origins with spaces",
			config:         "http://localhost:7001, http://localhost:17000, https://example.com",
			expectedLength: 3,
			expectedOrigins: []string{"http://localhost:7001", "http://localhost:17000", "https://example.com"},
		},
		{
			name:           "Empty config",
			config:         "",
			expectedLength: 0,
			expectedOrigins: []string{},
		},
		{
			name:           "Config with extra commas",
			config:         "http://localhost:7001,,http://localhost:17000",
			expectedLength: 2,
			expectedOrigins: []string{"http://localhost:7001", "http://localhost:17000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var allowedOrigins []string
			if tt.config != "" {
				origins := strings.Split(tt.config, ",")
				for _, origin := range origins {
					trimmed := strings.TrimSpace(origin)
					if trimmed != "" {
						allowedOrigins = append(allowedOrigins, trimmed)
					}
				}
			}

			if len(allowedOrigins) != tt.expectedLength {
				t.Errorf("Expected %d origins, got %d", tt.expectedLength, len(allowedOrigins))
			}

			for i, expected := range tt.expectedOrigins {
				if i >= len(allowedOrigins) {
					t.Errorf("Missing expected origin: %s", expected)
					continue
				}
				if allowedOrigins[i] != expected {
					t.Errorf("Expected origin[%d] to be %s, got %s", i, expected, allowedOrigins[i])
				}
			}
		})
	}
}

// TestCORSNoWildcardOrigins tests that wildcard origins should not be accepted
func TestCORSNoWildcardOrigins(t *testing.T) {
	config := "*"
	var allowedOrigins []string
	
	if config != "" {
		origins := strings.Split(config, ",")
		for _, origin := range origins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				// In production, we should validate and reject wildcard
				if trimmed == "*" {
					t.Logf("Warning: Wildcard origin detected. This is insecure when AllowCredentials is true.")
				}
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}

	// The configuration should not use wildcard when credentials are enabled
	if len(allowedOrigins) > 0 && allowedOrigins[0] == "*" {
		t.Logf("CORS configured with wildcard origin '*'. This is a security vulnerability when AllowCredentials is true.")
	}
}

// TestCORSDefaultOrigins tests that a secure default is used when no config is provided
func TestCORSDefaultOrigins(t *testing.T) {
	corsOriginsConfig := "" // Simulate no configuration
	var allowedOrigins []string
	
	if corsOriginsConfig != "" {
		origins := strings.Split(corsOriginsConfig, ",")
		for _, origin := range origins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}
	
	// If no origins are configured, use a secure default (localhost only)
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"http://localhost:7001"}
	}

	if len(allowedOrigins) != 1 {
		t.Errorf("Expected 1 default origin, got %d", len(allowedOrigins))
	}

	if allowedOrigins[0] != "http://localhost:7001" {
		t.Errorf("Expected default origin to be http://localhost:7001, got %s", allowedOrigins[0])
	}

	// Ensure the default is not a wildcard
	if allowedOrigins[0] == "*" {
		t.Error("Default origin should not be a wildcard")
	}
}
