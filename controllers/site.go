package controllers

import (
	"encoding/json"

	"github.com/casbin/caswaf/object"
)

func (c *ApiController) GetGlobalSites() {
	c.Data["json"] = object.GetGlobalSites()
	c.ServeJSON()
}

func (c *ApiController) GetSites() {
	owner := c.Input().Get("owner")

	c.Data["json"] = object.GetSites(owner)
	c.ServeJSON()
}

func (c *ApiController) GetSite() {
	id := c.Input().Get("id")

	c.Data["json"] = object.GetSite(id)
	c.ServeJSON()
}

func (c *ApiController) UpdateSite() {
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
	var site object.Site
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &site)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = object.AddSite(&site)
	c.ServeJSON()
}

func (c *ApiController) DeleteSite() {
	var site object.Site
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &site)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = object.DeleteSite(&site)
	c.ServeJSON()
}
