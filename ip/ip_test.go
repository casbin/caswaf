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

package ip

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsAbroadIp(t *testing.T) {
	// Change to the project root directory for the test
	wd, _ := os.Getwd()
	if filepath.Base(wd) == "ip" {
		os.Chdir("..")
		defer os.Chdir(wd)
	}

	// Initialize the IP database
	InitIpDb()

	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		// Intranet IPs should return false
		{"Private 10.x.x.x", "10.0.0.1", false},
		{"Private 192.168.x.x", "192.168.1.1", false},
		{"Private 172.16.x.x", "172.16.0.1", false},
		{"Loopback", "127.0.0.1", false},
		{"Link-local", "169.254.1.1", false},

		// Note: Testing public IPs requires the IP database to be properly loaded
		// The actual behavior depends on the IP geolocation database content
		// These tests verify that intranet IPs are handled correctly
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAbroadIp(tt.ip)
			if result != tt.expected {
				t.Errorf("IsAbroadIp(%s) = %v, expected %v", tt.ip, result, tt.expected)
			}
		})
	}
}
