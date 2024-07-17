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

package object

import (
	"github.com/casbin/caswaf/util"
	"github.com/xorm-io/core"
)

type Expression struct {
	Name     string `json:"name"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Rule struct {
	Owner       string       `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string       `xorm:"varchar(100) notnull pk" json:"name"`
	Type        string       `xorm:"varchar(100) notnull" json:"type"`
	Expressions []Expression `xorm:"mediumtext" json:"expressions"`
	CreatedTime string       `xorm:"varchar(100) notnull" json:"createdTime"`
	UpdatedTime string       `xorm:"varchar(100) notnull" json:"updatedTime"`
}

func GetRules() ([]*Rule, error) {
	rules := []*Rule{}
	err := ormer.Engine.Asc("owner").Desc("created_time").Find(&rules)
	return rules, err
}

func getRule(owner string, name string) (*Rule, error) {
	rule := Rule{Owner: owner, Name: name}
	existed, err := ormer.Engine.Get(&rule)
	if err != nil {
		return nil, err
	}
	if existed {
		return &rule, nil
	} else {
		return nil, nil
	}
}

func GetRule(id string) (*Rule, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getRule(owner, name)
}

func UpdateRule(id string, rule *Rule) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if s, err := getRule(owner, name); err != nil {
		return false, err
	} else if s == nil {
		return false, nil
	}
	rule.UpdatedTime = util.GetCurrentTime()
	_, err := ormer.Engine.ID(core.PK{owner, name}).AllCols().Update(rule)
	if err != nil {
		return false, err
	}
	return true, nil
}

func AddRule(rule *Rule) (bool, error) {
	if _, err := ormer.Engine.Insert(rule); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func DeleteRule(rule *Rule) (bool, error) {
	affected, err := ormer.Engine.ID(core.PK{rule.Owner, rule.Name}).Delete(&Rule{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func GetWAFRules() string {
	// Get all rules of type "waf".
	rules := []*Rule{}
	err := ormer.Engine.Where("type = ?", "waf").Find(&rules)
	if err != nil {
		return ""
	}

	res := ""
	// get all expressions from rules
	for _, rule := range rules {
		for _, expression := range rule.Expressions {
			res += expression.Value + "\n"
		}
	}
	return res
}
