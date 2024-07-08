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

type Rule struct {
	Id          int64  `xorm:"id pk autoincr" json:"id"`
	Rule        string `xorm:"varchar(512) notnull" json:"rule"`
	IsActive    bool   `xorm:"bool" json:"isActive"`
	CreatedTime string `xorm:"varchar(100) notnull" json:"createdTime"`
	UpdatedTime string `xorm:"varchar(100) notnull" json:"updatedTime"`
}

func GetRules() ([]*Rule, error) {
	rules := []*Rule{}
	err := ormer.Engine.Asc("id").Find(&rules)
	return rules, err
}

func GetRule(id int64) (*Rule, error) {
	rule := Rule{Id: id}
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

func UpdateRule(id int64, rule *Rule) (bool, error) {
	if affected, err := ormer.Engine.ID(id).AllCols().Update(rule); err != nil {
		return false, err
	} else {
		return affected != 0, nil
	}
}

func AddRule(rule *Rule) (bool, error) {
	if _, err := ormer.Engine.Insert(rule); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func DeleteRule(id int64) (bool, error) {
	rule := Rule{Id: id}
	if affected, err := ormer.Engine.Delete(&rule); err != nil {
		return false, err
	} else {
		return affected != 0, nil
	}
}

func GetWAFRules() string {
	objects, err := GetRules()
	if err != nil {
		return ""
	}

	var res string

	for _, rule := range objects {
		if rule.IsActive {
			res += rule.Rule + "\n"
		}
	}
	return res
}
