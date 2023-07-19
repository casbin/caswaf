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

func (c *ApiController) GetGlobalCerts() {
	c.Data["json"] = object.GetGlobalCerts()
	c.ServeJSON()
}

func (c *ApiController) GetCerts() {
	owner := c.Input().Get("owner")

	c.Data["json"] = object.GetCerts(owner)
	c.ServeJSON()
}

func (c *ApiController) GetCert() {
	id := c.Input().Get("id")

	c.Data["json"] = object.GetCert(id)
	c.ServeJSON()
}

func (c *ApiController) UpdateCert() {
	id := c.Input().Get("id")

	var cert object.Cert
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &cert)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = object.UpdateCert(id, &cert)
	c.ServeJSON()
}

func (c *ApiController) AddCert() {
	var cert object.Cert
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &cert)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = object.AddCert(&cert)
	c.ServeJSON()
}

func (c *ApiController) DeleteCert() {
	var cert object.Cert
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &cert)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = object.DeleteCert(&cert)
	c.ServeJSON()
}
