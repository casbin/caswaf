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

package object

import (
	"fmt"

	"github.com/casbin/caswaf/util"
	"github.com/xorm-io/core"
)

type Action struct {
	Owner         string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name          string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime   string `xorm:"varchar(100) notnull" json:"createdTime"`
	Type          string `xorm:"varchar(100) notnull" json:"type"`
	StatusCode    int    `xorm:"int notnull" json:"statusCode"`
	ImmunityTimes int    `xorm:"int notnull" json:"immunityTimes"` // minutes
}

func GetGlobalActions() ([]*Action, error) {
	actions := []*Action{}
	err := ormer.Engine.Asc("owner").Desc("created_time").Find(&actions)
	return actions, err
}

func GetActions(owner string) ([]*Action, error) {
	actions := []*Action{}
	err := ormer.Engine.Desc("updated_time").Find(&actions, &Action{Owner: owner})
	return actions, err
}

func getAction(owner string, name string) (*Action, error) {
	action := Action{Owner: owner, Name: name}
	existed, err := ormer.Engine.Get(&action)
	if err != nil {
		return nil, err
	}
	if existed {
		return &action, nil
	} else {
		return nil, nil
	}
}

func GetAction(id string) (*Action, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getAction(owner, name)
}

func UpdateAction(id string, action *Action) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if s, err := getAction(owner, name); err != nil {
		return false, err
	} else if s == nil {
		return false, nil
	}
	_, err := ormer.Engine.ID(core.PK{owner, name}).AllCols().Update(action)
	if err != nil {
		return false, err
	}
	err = refreshActionMap()
	if err != nil {
		return false, err
	}
	return true, nil
}

func AddAction(action *Action) (bool, error) {
	affected, err := ormer.Engine.Insert(action)
	if err != nil {
		return false, err
	}
	if affected != 0 {
		err = refreshActionMap()
		if err != nil {
			return false, err
		}
	}
	return affected != 0, nil
}

func DeleteAction(action *Action) (bool, error) {
	affected, err := ormer.Engine.ID(core.PK{action.Owner, action.Name}).Delete(&Action{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func (action *Action) GetId() string {
	return fmt.Sprintf("%s/%s", action.Owner, action.Name)
}
