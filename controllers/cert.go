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
