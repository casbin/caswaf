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
)

func (c *ApiController) GetActions() {
	if c.RequireSignedIn() {
		return
	}
	owner := c.Input().Get("owner")
	if owner == "admin" {
		owner = ""
	}

	actions, err := object.GetActions(owner)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(actions)
}

func (c *ApiController) GetAction() {
	if c.RequireSignedIn() {
		return
	}

	id := c.Input().Get("id")
	action, err := object.GetAction(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(action)
}

func (c *ApiController) AddAction() {
	if c.RequireSignedIn() {
		return
	}

	var action object.Action
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &action)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.Data["json"] = wrapActionResponse(object.AddAction(&action))
	c.ServeJSON()
}

func (c *ApiController) UpdateAction() {
	if c.RequireSignedIn() {
		return
	}

	var action object.Action
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &action)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	id := c.Input().Get("id")
	c.Data["json"] = wrapActionResponse(object.UpdateAction(id, &action))
	c.ServeJSON()
}

func (c *ApiController) DeleteAction() {
	if c.RequireSignedIn() {
		return
	}

	var action object.Action
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &action)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(object.DeleteAction(&action))
	c.ServeJSON()
}
