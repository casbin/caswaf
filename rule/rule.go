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
	checkRule(expressions []*object.Expression, req *http.Request) (bool, string, error)
}

func CheckRules(wafRuleIds []string, r *http.Request) (bool, string, error) {
	rules := object.GetRulesByRuleIds(wafRuleIds)
	for _, rule := range rules {
		var ruleObj Rule
		switch rule.Type {
		case "User-Agent":
			ruleObj = &UaRule{}
		default:
			return false, "", fmt.Errorf("unknown rule type: %s for rule: %s", rule.Type, rule.GetId())
		}

		isHit, reason, err := ruleObj.checkRule(rule.Expressions, r)
		if err != nil {
			return false, "", err
		}

		if isHit {
			if rule.Action == "Block" {
				if rule.Reason != "" {
					reason = rule.Reason
				}

				return false, reason, nil
			} else if rule.Action == "Allow" {
				return true, "", nil
			} else {
				return false, "", fmt.Errorf("unknown rule action: %s for rule: %s", rule.Action, rule.GetId())
			}
		}
	}

	return true, "", nil
}
