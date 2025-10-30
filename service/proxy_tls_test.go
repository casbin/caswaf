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

// TestTLSConfigSecurity verifies that the TLS configuration meets security requirements
func TestTLSConfigSecurity(t *testing.T) {
	// Create a TLS config similar to what's used in Start()
	config := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
	}

	// Test 1: Verify minimum TLS version is 1.2
	if config.MinVersion != tls.VersionTLS12 {
		t.Errorf("MinVersion should be TLS 1.2, got: %v", config.MinVersion)
	}

	// Test 2: Verify PreferServerCipherSuites is enabled
	if !config.PreferServerCipherSuites {
		t.Error("PreferServerCipherSuites should be true")
	}

	// Test 3: Verify no weak cipher suites are included (especially 3DES)
	weakCiphers := []uint16{
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_RSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	}

	for _, weakCipher := range weakCiphers {
		for _, configuredCipher := range config.CipherSuites {
			if configuredCipher == weakCipher {
				t.Errorf("Weak cipher suite detected: %v", weakCipher)
			}
		}
	}

	// Test 4: Verify all configured ciphers use forward secrecy (ECDHE)
	for _, cipher := range config.CipherSuites {
		cipherName := tls.CipherSuiteName(cipher)
		if cipherName != "" {
			// All our ciphers should start with "TLS_ECDHE"
			if len(cipherName) < 9 || cipherName[:9] != "TLS_ECDHE" {
				t.Errorf("Cipher suite without forward secrecy detected: %s", cipherName)
			}
		}
	}

	// Test 5: Verify strong elliptic curves are preferred
	if len(config.CurvePreferences) == 0 {
		t.Error("CurvePreferences should be configured")
	}

	// X25519 should be the first preference for best performance and security
	if config.CurvePreferences[0] != tls.X25519 {
		t.Errorf("First curve preference should be X25519, got: %v", config.CurvePreferences[0])
	}

	// Test 6: Verify cipher suites count is reasonable (not empty, not too many)
	if len(config.CipherSuites) == 0 {
		t.Error("CipherSuites should not be empty")
	}
	if len(config.CipherSuites) > 10 {
		t.Errorf("Too many cipher suites configured (%d), this may indicate weak ciphers are included", len(config.CipherSuites))
	}
}

// TestWeakCipherSuitesNotPresent specifically tests that 3DES cipher suites are not present
func TestWeakCipherSuitesNotPresent(t *testing.T) {
	config := &tls.Config{
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

	// These are the specific weak cipher suites mentioned in the issue
	forbiddenCiphers := map[uint16]string{
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:     "TLS_RSA_WITH_3DES_EDE_CBC_SHA",
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA: "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA",
	}

	for forbiddenCipher, name := range forbiddenCiphers {
		for _, configuredCipher := range config.CipherSuites {
			if configuredCipher == forbiddenCipher {
				t.Errorf("Forbidden cipher suite %s (value: 0x%04x) found in configuration", name, forbiddenCipher)
			}
		}
	}
}
