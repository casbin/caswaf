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

	"github.com/casbin/caswaf/conf"
	"github.com/casbin/caswaf/object"
	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/types"
	"github.com/hsluoyz/modsecurity-go/seclang/parser"
)

type WafRule struct{}

func (r *WafRule) checkRule(expressions []*object.Expression, req *http.Request) (bool, string, string, error) {
	var ruleStr string
	for _, expression := range expressions {
		ruleStr += expression.Value
	}
	waf, err := coraza.NewWAF(
		coraza.NewWAFConfig().
			WithErrorCallback(logError).
			WithDirectives(conf.WafConf).
			WithDirectives(ruleStr),
	)
	if err != nil {
		return false, "", "", fmt.Errorf("create WAF failed")
	}
	tx := waf.NewTransaction()
	processRequest(tx, req)
	matchedRules := tx.MatchedRules()
	for _, matchedRule := range matchedRules {
		rule := matchedRule.Rule()
		directive, err := parser.NewSecLangScannerFromString(rule.Raw()).AllDirective()
		if err != nil {
			return false, "", "", err
		}
		for _, d := range directive {
			ruleDirective := d.(*parser.RuleDirective)
			for _, action := range ruleDirective.Actions.Action {
				switch action.Tk {
				case parser.TkActionBlock, parser.TkActionDeny:
					return true, "Block", fmt.Sprintf("blocked by WAF rule: %d", rule.ID()), nil
				case parser.TkActionAllow:
					return true, "Allow", "", nil
				case parser.TkActionDrop:
					return true, "Drop", fmt.Sprintf("dropped by WAF rule: %d", rule.ID()), nil
				default:
					// skip other actions
					continue
				}
			}
		}
	}
	return false, "", "", nil
}

func processRequest(tx types.Transaction, req *http.Request) {
	// Process URI and method
	tx.ProcessURI(req.URL.String(), req.Method, req.Proto)

	// Process request headers
	for key, values := range req.Header {
		for _, value := range values {
			tx.AddRequestHeader(key, value)
		}
	}
	tx.ProcessRequestHeaders()

	// Process request body (if any)
	if req.Body != nil {
		_, err := tx.ProcessRequestBody()
		if err != nil {
			return
		}
	}
}

func logError(error types.MatchedRule) {
	msg := error.ErrorLog()
	fmt.Printf("[WAFlogError][%s] %s\n", error.Rule().Severity(), msg)
}
