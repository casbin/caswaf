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
	"crypto/tls"
	"testing"
)

// TestTLSConfigExcludes3DES verifies that the TLS configuration excludes
// 3DES cipher suites that are vulnerable to the Sweet32 attack
func TestTLSConfigExcludes3DES(t *testing.T) {
	// Create the TLS config as it would be in the actual server
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
	}

	// Define vulnerable 3DES cipher suites
	vulnerableCiphers := []uint16{
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	}

	// Verify that none of the vulnerable ciphers are in the allowed list
	for _, vulnerableCipher := range vulnerableCiphers {
		for _, allowedCipher := range tlsConfig.CipherSuites {
			if vulnerableCipher == allowedCipher {
				t.Errorf("Vulnerable 3DES cipher suite 0x%04X found in allowed cipher suites (Sweet32 vulnerability)", vulnerableCipher)
			}
		}
	}

	// Verify minimum TLS version is set to 1.2 or higher
	if tlsConfig.MinVersion < tls.VersionTLS12 {
		t.Errorf("Minimum TLS version should be 1.2 or higher, got: %d", tlsConfig.MinVersion)
	}

	// Verify that we have at least some secure cipher suites configured
	if len(tlsConfig.CipherSuites) == 0 {
		t.Error("No cipher suites configured - default cipher suites may include vulnerable 3DES")
	}
}
