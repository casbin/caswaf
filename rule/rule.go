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
	"fmt"
	"net/http"

	"github.com/casbin/caswaf/object"
)

type Rule interface {
	checkRule(expressions []*object.Expression, req *http.Request) (bool, string, string, error)
}

type ActionResult struct {
	Type       string
	StatusCode int
}

func CheckRules(ruleIds []string, r *http.Request) (*ActionResult, string, error) {
	rules, err := object.GetRulesByRuleIds(ruleIds)
	if err != nil {
		return nil, "", err
	}
	for i, rule := range rules {
		var ruleObj Rule
		switch rule.Type {
		case "User-Agent":
			ruleObj = &UaRule{}
		case "IP":
			ruleObj = &IpRule{}
		case "WAF":
			ruleObj = &WafRule{}
		case "IP Rate Limiting":
			ruleObj = &IpRateRule{
				ruleName: rule.GetId(),
			}
		case "Compound":
			ruleObj = &CompoundRule{}
		default:
			return nil, "", fmt.Errorf("unknown rule type: %s for rule: %s", rule.Type, rule.GetId())
		}

		isHit, action, reason, err := ruleObj.checkRule(rule.Expressions, r)
		if err != nil {
			return nil, "", err
		}
		
		// Use rule's action if no action specified by the rule check
		if action == "" {
			action = rule.Action
		}
		
		// Determine status code
		statusCode := rule.StatusCode
		if statusCode == 0 {
			// Set default status codes if not specified
			switch action {
			case "Block":
				statusCode = 403
			case "Drop":
				statusCode = 400
			case "Allow":
				statusCode = 200
			case "CAPTCHA":
				statusCode = 302
			default:
				return nil, "", fmt.Errorf("unknown rule action: %s for rule: %s", action, rule.GetId())
			}
		}
		
		actionResult := &ActionResult{
			Type:       action,
			StatusCode: statusCode,
		}
		
		if isHit {
			if action == "Block" || action == "Drop" {
				if rule.Reason != "" {
					reason = rule.Reason
				} else {
					reason = fmt.Sprintf("hit rule %s: %s", ruleIds[i], reason)
				}
				return actionResult, reason, nil
			} else if action == "Allow" {
				return actionResult, reason, nil
			} else if action == "CAPTCHA" {
				return actionResult, reason, nil
			} else {
				return nil, "", fmt.Errorf("unknown rule action: %s for rule: %s", action, rule.GetId())
			}
		}
	}
	
	// Default action if no rule matched
	return &ActionResult{
		Type:       "Allow",
		StatusCode: 200,
	}, "", nil
}
