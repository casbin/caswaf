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
	"errors"
	"net"
	"strings"

	"github.com/casbin/caswaf/object"
	"github.com/casbin/caswaf/util"
	"github.com/hsluoyz/modsecurity-go/seclang/parser"
)

func (c *ApiController) GetRules() {
	if c.RequireSignedIn() {
		return
	}
	owner := c.Input().Get("owner")
	if owner == "admin" {
		owner = ""
	}

	rules, err := object.GetRules(owner)
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
	err = checkExpressions(rule.Expressions, rule.Type)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.Data["json"] = wrapActionResponse(object.AddRule(&rule))
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

	err = checkExpressions(rule.Expressions, rule.Type)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	id := c.Input().Get("id")
	c.Data["json"] = wrapActionResponse(object.UpdateRule(id, &rule))
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
	c.ServeJSON()
}

func checkExpressions(expressions []*object.Expression, ruleType string) error {
	values := make([]string, len(expressions))
	for i, expression := range expressions {
		values[i] = expression.Value
	}
	switch ruleType {
	case "WAF":
		return checkWafRule(values)
	case "IP":
		return checkIpRule(values)
	case "IpRate":
		return checkIpRateRule(expressions)
	}
	return nil
}

func checkWafRule(rules []string) error {
	for _, rule := range rules {
		scanner := parser.NewSecLangScannerFromString(rule)
		_, err := scanner.AllDirective()
		if err != nil {
			return err
		}
	}
	return nil
}

func checkIpRule(ipLists []string) error {
	for _, ipList := range ipLists {
		for _, ip := range strings.Split(ipList, ",") {
			_, _, err := net.ParseCIDR(ip)
			if net.ParseIP(ip) == nil && err != nil {
				return errors.New("Invalid IP address: " + ip)
			}
		}
	}
	return nil
}

func checkIpRateRule(expressions []*object.Expression) error {
	if len(expressions) != 1 {
		return errors.New("IpRate rule should have only one expression")
	}
	expression := expressions[0]
	_, err := util.ParseIntWithError(expression.Operator)
	if err != nil {
		return err
	}
	_, err = util.ParseIntWithError(expression.Value)
	if err != nil {
		return err
	}
	return nil
}
