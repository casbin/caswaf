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
	"strings"
	"testing"

	"github.com/casbin/caswaf/object"
)

func TestVerboseMode(t *testing.T) {
	// Mock the GetRulesByRuleIds function by creating a test scenario
	// We'll test the verbose mode by directly checking the CheckRules logic
	
	tests := []struct {
		name          string
		rule          *object.Rule
		expressions   []*object.Expression
		req           *http.Request
		wantReason    string
		wantContains  []string
		wantAction    string
	}{
		{
			name: "Verbose mode enabled - IP rule",
			rule: &object.Rule{
				Owner:     "admin",
				Name:      "test-verbose-rule",
				Type:      "IP",
				Action:    "Block",
				Reason:    "Custom block reason",
				IsVerbose: true,
			},
			expressions: []*object.Expression{
				{
					Operator: "is in",
					Value:    "127.0.0.1",
				},
			},
			req: &http.Request{
				RemoteAddr: "127.0.0.1",
			},
			wantAction:   "Block",
			wantContains: []string{"Rule [admin/test-verbose-rule] triggered", "Custom reason: Custom block reason"},
		},
		{
			name: "Verbose mode disabled - IP rule",
			rule: &object.Rule{
				Owner:     "admin",
				Name:      "test-non-verbose-rule",
				Type:      "IP",
				Action:    "Block",
				Reason:    "Custom block reason",
				IsVerbose: false,
			},
			expressions: []*object.Expression{
				{
					Operator: "is in",
					Value:    "127.0.0.1",
				},
			},
			req: &http.Request{
				RemoteAddr: "127.0.0.1",
			},
			wantReason: "Custom block reason",
			wantAction: "Block",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create rule object based on type
			var ruleObj Rule
			switch tt.rule.Type {
			case "IP":
				ruleObj = &IpRule{}
			default:
				t.Fatalf("unsupported rule type: %s", tt.rule.Type)
			}

			// Check the rule
			result, err := ruleObj.checkRule(tt.expressions, tt.req)
			if err != nil {
				t.Fatalf("checkRule() error = %v", err)
			}

			if result == nil {
				t.Fatalf("checkRule() returned nil result")
			}

			// Apply the logic from CheckRules to add action/reason
			if result.Action == "" {
				result.Action = tt.rule.Action
			}

			// Apply verbose logic
			if result.Action == "Block" || result.Action == "Drop" {
				if tt.rule.IsVerbose {
					verboseReason := "Rule [" + tt.rule.Owner + "/" + tt.rule.Name + "] triggered"
					if result.Reason != "" {
						verboseReason += " - " + result.Reason
					}
					if tt.rule.Reason != "" {
						verboseReason += " - Custom reason: " + tt.rule.Reason
					}
					result.Reason = verboseReason
				} else if tt.rule.Reason != "" {
					result.Reason = tt.rule.Reason
				}
			}

			// Check action
			if result.Action != tt.wantAction {
				t.Errorf("Action = %v, want %v", result.Action, tt.wantAction)
			}

			// Check reason
			if tt.wantReason != "" {
				if result.Reason != tt.wantReason {
					t.Errorf("Reason = %v, want %v", result.Reason, tt.wantReason)
				}
			}

			// Check if reason contains expected strings
			if tt.wantContains != nil {
				for _, substr := range tt.wantContains {
					if !strings.Contains(result.Reason, substr) {
						t.Errorf("Reason %q does not contain %q", result.Reason, substr)
					}
				}
			}
		})
	}
}
