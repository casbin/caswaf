package rule

import (
	"net/http"
	"os"
	"testing"

	"github.com/casbin/caswaf/ip"
	"github.com/casbin/caswaf/object"
)

func init() {
	// Initialize IP database for tests
	// Check if we're in the rule directory and adjust path if needed
	if _, err := os.Stat("../ip/17monipdb.dat"); err == nil {
		os.Chdir("..")
	}
	ip.InitIpDb()
}

func TestIpRule_IsAbroad(t *testing.T) {
	ipRule := &IpRule{}

	tests := []struct {
		name        string
		operator    string
		value       string
		clientIP    string
		shouldMatch bool
		shouldError bool
	}{
		{
			name:        "is abroad with empty value - foreign IP",
			operator:    "is abroad",
			value:       "",
			clientIP:    "8.8.8.8",
			shouldMatch: true,
			shouldError: false,
		},
		{
			name:        "is abroad with some value - foreign IP",
			operator:    "is abroad",
			value:       "1.1.1.1",
			clientIP:    "8.8.8.8",
			shouldMatch: true,
			shouldError: false,
		},
		{
			name:        "is abroad with CIDR value - foreign IP",
			operator:    "is abroad",
			value:       "1.1.1.0/24",
			clientIP:    "8.8.8.8",
			shouldMatch: true,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expressions := []*object.Expression{
				{
					Operator: tt.operator,
					Value:    tt.value,
				},
			}

			req := &http.Request{
				RemoteAddr: tt.clientIP + ":1234",
			}

			result, err := ipRule.checkRule(expressions, req)

			if tt.shouldError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			gotMatch := result != nil
			if gotMatch != tt.shouldMatch {
				t.Errorf("Expected match: %v, got: %v", tt.shouldMatch, gotMatch)
			}
		})
	}
}
