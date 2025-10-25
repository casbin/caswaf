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

package util

import "testing"

func TestIsIntranetIp(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		// Private IPv4 addresses
		{"Private 10.x.x.x", "10.0.0.1", true},
		{"Private 192.168.x.x", "192.168.1.1", true},
		{"Private 172.16.x.x", "172.16.0.1", true},
		{"Private 172.31.x.x", "172.31.255.255", true},

		// Loopback
		{"Loopback IPv4", "127.0.0.1", true},
		{"Loopback IPv6", "::1", true},

		// Link-local
		{"Link-local IPv4", "169.254.1.1", true},
		{"Link-local IPv6", "fe80::1", true},

		// Public IPs
		{"Public Google DNS", "8.8.8.8", false},
		{"Public Cloudflare DNS", "1.1.1.1", false},
		{"Public IPv4", "123.45.67.89", false},

		// With port
		{"Private with port", "192.168.1.1:8080", true},
		{"Public with port", "8.8.8.8:53", false},

		// Invalid IPs
		{"Invalid IP", "invalid", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsIntranetIp(tt.ip)
			if result != tt.expected {
				t.Errorf("IsIntranetIp(%s) = %v, expected %v", tt.ip, result, tt.expected)
			}
		})
	}
}
