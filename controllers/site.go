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

func (c *ApiController) GetGlobalSites() {
	if c.RequireSignedIn() {
		return
	}

	c.Data["json"] = object.GetGlobalSites()
	c.ServeJSON()
}

func (c *ApiController) GetSites() {
	if c.RequireSignedIn() {
		return
	}

	owner := c.Input().Get("owner")

	c.Data["json"] = object.GetSites(owner)
	c.ServeJSON()
}

func (c *ApiController) GetSite() {
	if c.RequireSignedIn() {
		return
	}

	id := c.Input().Get("id")

	c.Data["json"] = object.GetSite(id)
	c.ServeJSON()
}

func (c *ApiController) UpdateSite() {
	if c.RequireSignedIn() {
		return
	}

	id := c.Input().Get("id")

	var site object.Site
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &site)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = object.UpdateSite(id, &site)
	c.ServeJSON()
}

func (c *ApiController) AddSite() {
	if c.RequireSignedIn() {
		return
	}

	var site object.Site
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &site)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = object.AddSite(&site)
	c.ServeJSON()
}

func (c *ApiController) DeleteSite() {
	if c.RequireSignedIn() {
		return
	}

	var site object.Site
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &site)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = object.DeleteSite(&site)
	c.ServeJSON()
}
