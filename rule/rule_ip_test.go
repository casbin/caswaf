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

package rule

import (
	"net/http"
	"testing"

	"github.com/casbin/caswaf/ip"
	"github.com/casbin/caswaf/object"
)

func TestIpRule_checkRule_IsAbroad(t *testing.T) {
	// Initialize IP database for testing
	err := ip.Init("../ip/17monipdb.dat")
	if err != nil {
		t.Fatalf("Failed to initialize IP database: %v", err)
	}

	tests := []struct {
		name        string
		expressions []*object.Expression
		remoteAddr  string
		wantMatch   bool
		wantErr     bool
	}{
		{
			name: "Test abroad IP (US IP)",
			expressions: []*object.Expression{
				{
					Operator: "is abroad",
					Value:    "",
				},
			},
			remoteAddr: "8.8.8.8",
			wantMatch:  true,
			wantErr:    false,
		},
		{
			name: "Test China IP (should not match abroad)",
			expressions: []*object.Expression{
				{
					Operator: "is abroad",
					Value:    "",
				},
			},
			remoteAddr: "61.135.169.121", // Baidu IP in China
			wantMatch:  false,
			wantErr:    false,
		},
		{
			name: "Test another China IP (should not match abroad)",
			expressions: []*object.Expression{
				{
					Operator: "is abroad",
					Value:    "",
				},
			},
			remoteAddr: "220.181.38.148", // Another Baidu IP in China
			wantMatch:  false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &IpRule{}
			req := &http.Request{
				RemoteAddr: tt.remoteAddr + ":8080",
			}
			result, err := r.checkRule(tt.expressions, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := result != nil
			if got != tt.wantMatch {
				t.Errorf("checkRule() matched = %v, want %v", got, tt.wantMatch)
			}
		})
	}
}

func TestIpRule_checkRule_ExistingOperators(t *testing.T) {
	tests := []struct {
		name        string
		expressions []*object.Expression
		remoteAddr  string
		wantMatch   bool
		wantErr     bool
	}{
		{
			name: "Test is in with single IP",
			expressions: []*object.Expression{
				{
					Operator: "is in",
					Value:    "192.168.1.1",
				},
			},
			remoteAddr: "192.168.1.1",
			wantMatch:  true,
			wantErr:    false,
		},
		{
			name: "Test is in with CIDR",
			expressions: []*object.Expression{
				{
					Operator: "is in",
					Value:    "192.168.1.0/24",
				},
			},
			remoteAddr: "192.168.1.100",
			wantMatch:  true,
			wantErr:    false,
		},
		{
			name: "Test is not in",
			expressions: []*object.Expression{
				{
					Operator: "is not in",
					Value:    "10.0.0.1",
				},
			},
			remoteAddr: "192.168.1.1",
			wantMatch:  true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &IpRule{}
			req := &http.Request{
				RemoteAddr: tt.remoteAddr + ":8080",
			}
			result, err := r.checkRule(tt.expressions, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := result != nil
			if got != tt.wantMatch {
				t.Errorf("checkRule() matched = %v, want %v", got, tt.wantMatch)
			}
		})
	}
}
