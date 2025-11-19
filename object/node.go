// Copyright 2025 The casbin Authors. All Rights Reserved.
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

type NodeItem struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	UpdatedTime string `xorm:"varchar(100)" json:"updatedTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`

	Tag         string `xorm:"varchar(100)" json:"tag"`
	Hostname    string `xorm:"varchar(100)" json:"hostname"`
	IpAddress   string `xorm:"varchar(100)" json:"ipAddress"`
	Description string `xorm:"varchar(500)" json:"description"`
}

func GetGlobalNodes() ([]*NodeItem, error) {
	nodes := []*NodeItem{}
	err := ormer.Engine.Desc("created_time").Find(&nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func GetNodes(owner string) ([]*NodeItem, error) {
	nodes := []*NodeItem{}
	err := ormer.Engine.Desc("created_time").Find(&nodes, &NodeItem{Owner: owner})
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func getNode(owner string, name string) (*NodeItem, error) {
	node := NodeItem{Owner: owner, Name: name}
	existed, err := ormer.Engine.Get(&node)
	if err != nil {
		return nil, err
	}

	if existed {
		return &node, nil
	}
	return nil, nil
}

func GetNode(id string) (*NodeItem, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getNode(owner, name)
}

func UpdateNode(id string, node *NodeItem) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if n, err := getNode(owner, name); err != nil {
		return false, err
	} else if n == nil {
		return false, nil
	}

	node.UpdatedTime = util.GetCurrentTime()

	_, err := ormer.Engine.ID(core.PK{owner, name}).AllCols().Update(node)
	if err != nil {
		return false, err
	}

	return true, nil
}

func AddNode(node *NodeItem) (bool, error) {
	affected, err := ormer.Engine.Insert(node)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func DeleteNode(node *NodeItem) (bool, error) {
	affected, err := ormer.Engine.ID(core.PK{node.Owner, node.Name}).Delete(&NodeItem{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func (node *NodeItem) GetId() string {
	return fmt.Sprintf("%s/%s", node.Owner, node.Name)
}

func GetNodeCount(owner, field, value string) (int64, error) {
	session := GetSession(owner, -1, -1, field, value, "", "")
	return session.Count(&NodeItem{})
}

func GetPaginationNodes(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*NodeItem, error) {
	nodes := []*NodeItem{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Where("owner = ? or owner = ?", "admin", owner).Find(&nodes)
	if err != nil {
		return nodes, err
	}

	return nodes, nil
}
