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

package controllers

import (
	"encoding/json"

	"github.com/casbin/caswaf/object"
)

func (c *ApiController) GetRecords() {
	if c.RequireSignedIn() {
		return
	}

	owner := c.Input().Get("owner")
	if owner == "admin" {
		owner = ""
	}

	sites, err := object.GetRecords(owner)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	// object.GetMaskedSites(sites, util.GetHostname())
	c.ResponseOk(sites)
}

func (c *ApiController) DeleteRecord() {
	if c.RequireSignedIn() {
		return
	}

	var record object.Record
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &record)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(object.DeleteRecord(&record))
	c.ServeJSON()
}

func (c *ApiController) UpdateRecord() {
	if c.RequireSignedIn() {
		return
	}

	owner := c.Input().Get("owner")
	id := c.Input().Get("id")

	var record object.Record
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &record)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(object.UpdateRecord(owner, id, &record))
	c.ServeJSON()
}

func (c *ApiController) GetRecord() {
	if c.RequireSignedIn() {
		return
	}

	owner := c.Input().Get("owner")
	id := c.Input().Get("id")
	record, err := object.GetRecord(owner, id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(record)
}

func (c *ApiController) AddRecord() {
	if c.RequireSignedIn() {
		return
	}

	var record object.Record
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &record)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(object.AddRecord(&record))
	c.ServeJSON()
}
