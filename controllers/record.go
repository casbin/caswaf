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
