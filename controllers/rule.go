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

package controllers

import (
	"encoding/json"

	"github.com/casbin/caswaf/object"
	"github.com/casbin/caswaf/service"
	"github.com/casbin/caswaf/util"
	"github.com/hsluoyz/modsecurity-go/seclang/parser"
)

func (c *ApiController) GetRules() {
	if c.RequireSignedIn() {
		return
	}

	rules, err := object.GetRules()
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(rules)
}

func (c *ApiController) GetRule() {
	if c.RequireSignedIn() {
		return
	}

	id := c.Input().Get("id")
	rule, err := object.GetRule(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(rule)
}

func (c *ApiController) AddRule() {
	if c.RequireSignedIn() {
		return
	}

	currentTime := util.GetCurrentTime()
	rule := object.Rule{
		CreatedTime: currentTime,
		UpdatedTime: currentTime,
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &rule)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	err = checkWAFRule(makeWAFRules(rule.Expressions))
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.Data["json"] = wrapActionResponse(object.AddRule(&rule))
	go service.UpdateWAF()
	c.ServeJSON()
}

func (c *ApiController) UpdateRule() {
	if c.RequireSignedIn() {
		return
	}

	var rule object.Rule
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &rule)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	err = checkWAFRule(makeWAFRules(rule.Expressions))
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	id := c.Input().Get("id")
	c.Data["json"] = wrapActionResponse(object.UpdateRule(id, &rule))
	go service.UpdateWAF()
	c.ServeJSON()
}

func (c *ApiController) DeleteRule() {
	if c.RequireSignedIn() {
		return
	}

	var rule object.Rule
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &rule)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(object.DeleteRule(&rule))
	go service.UpdateWAF()
	c.ServeJSON()
}

func makeWAFRules(expressions []object.Expression) []string {
	rules := make([]string, len(expressions))
	for i, expression := range expressions {
		rules[i] = expression.Value
	}
	return rules
}

func checkWAFRule(rules []string) error {
	for _, rule := range rules {
		scanner := parser.NewSecLangScannerFromString(rule)
		_, err := scanner.AllDirective()
		if err != nil {
			return err
		}
	}
	return nil
}
