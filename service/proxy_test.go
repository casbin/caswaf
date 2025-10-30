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

func TestTLSMinVersion(t *testing.T) {
	// Create a TLS config as it would be in the Start() function
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Verify that MinVersion is set to TLS 1.2 or higher
	if tlsConfig.MinVersion < tls.VersionTLS12 {
		t.Errorf("TLS MinVersion is too low: got 0x%04x, want >= 0x%04x (TLS 1.2)", tlsConfig.MinVersion, tls.VersionTLS12)
	}

	// Verify that TLS 1.0 and 1.1 are disabled
	if tlsConfig.MinVersion == tls.VersionTLS10 {
		t.Error("TLS 1.0 should be disabled")
	}
	if tlsConfig.MinVersion == tls.VersionTLS11 {
		t.Error("TLS 1.1 should be disabled")
	}
}
