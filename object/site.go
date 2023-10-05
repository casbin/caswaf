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

import (
	"fmt"

	"github.com/casbin/caswaf/run"
	"github.com/casbin/caswaf/util"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"xorm.io/core"
)

type Node struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Diff    string `json:"diff"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Site struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	UpdatedTime string `xorm:"varchar(100)" json:"updatedTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`

	Domain   string  `xorm:"varchar(100)" json:"domain"`
	Host     string  `xorm:"varchar(100)" json:"host"`
	SslMode  string  `xorm:"varchar(100)" json:"sslMode"`
	SslCert  string  `xorm:"varchar(100)" json:"sslCert"`
	PublicIp string  `xorm:"varchar(100)" json:"publicIp"`
	Node     string  `xorm:"varchar(100)" json:"node"`
	IsSelf   bool    `json:"isSelf"`
	Nodes    []*Node `json:"nodes"`

	CasdoorApplication string `xorm:"varchar(100)" json:"casdoorApplication"`

	SslCertObj     *Cert                   `xorm:"-" json:"sslCertObj"`
	ApplicationObj *casdoorsdk.Application `xorm:"-" json:"applicationObj"`
}

func GetGlobalSites() []*Site {
	sites := []*Site{}
	err := ormer.Engine.Asc("owner").Desc("created_time").Find(&sites)
	if err != nil {
		panic(err)
	}

	return sites
}

func GetSites(owner string) []*Site {
	sites := []*Site{}
	err := ormer.Engine.Desc("created_time").Find(&sites, &Site{Owner: owner})
	if err != nil {
		panic(err)
	}

	for _, site := range sites {
		if site.SslCert == "" && site.Domain != "" && (site.SslMode == "HTTPS and HTTP" || site.SslMode == "HTTPS Only") {
			site.SslCert = getBaseDomain(site.Domain)
			UpdateSite(site.GetId(), site)
		}
	}

	return sites
}

func getSite(owner string, name string) *Site {
	site := Site{Owner: owner, Name: name}
	existed, err := ormer.Engine.Get(&site)
	if err != nil {
		panic(err)
	}

	if existed {
		return &site
	} else {
		return nil
	}
}

func GetSite(id string) *Site {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getSite(owner, name)
}

func GetMaskedSite(site *Site, node string) *Site {
	if site == nil {
		return nil
	}

	if site.PublicIp == "(empty)" {
		site.PublicIp = ""
	}

	site.IsSelf = false
	if site.Node == node {
		site.IsSelf = true
	}

	return site
}

func GetMaskedSites(sites []*Site, node string) []*Site {
	for _, site := range sites {
		site = GetMaskedSite(site, node)
	}
	return sites
}

func UpdateSite(id string, site *Site) bool {
	owner, name := util.GetOwnerAndNameFromId(id)
	if getSite(owner, name) == nil {
		return false
	}

	site.UpdatedTime = util.GetCurrentTime()

	_, err := ormer.Engine.ID(core.PK{owner, name}).AllCols().Update(site)
	if err != nil {
		panic(err)
	}

	refreshSiteMap()
	site.checkNodes()

	//return affected != 0
	return true
}

func UpdateSiteNoRefresh(id string, site *Site) bool {
	owner, name := util.GetOwnerAndNameFromId(id)
	if getSite(owner, name) == nil {
		return false
	}

	_, err := ormer.Engine.ID(core.PK{owner, name}).AllCols().Update(site)
	if err != nil {
		panic(err)
	}

	//return affected != 0
	return true
}

func AddSite(site *Site) bool {
	affected, err := ormer.Engine.Insert(site)
	if err != nil {
		panic(err)
	}

	if affected != 0 {
		refreshSiteMap()
	}

	return affected != 0
}

func DeleteSite(site *Site) bool {
	affected, err := ormer.Engine.ID(core.PK{site.Owner, site.Name}).Delete(&Site{})
	if err != nil {
		panic(err)
	}

	if affected != 0 {
		refreshSiteMap()
	}

	return affected != 0
}

func (site *Site) GetId() string {
	return fmt.Sprintf("%s/%s", site.Owner, site.Name)
}

func (site *Site) checkNodes() {
	hostname := util.GetHostname()
	for i, node := range site.Nodes {
		if node.Name != hostname {
			continue
		}

		if site.Host == "" {
			continue
		}

		ok, msg := pingUrl(site.Host)
		status := "Running"
		if !ok {
			status = "Stopped"
		}

		run.CreateRepo(site.Name)

		version := getSiteVersion(site.Name)

		path := run.GetRepoPath(site.Name)
		diff := run.GitDiff(path)

		if node.Status != status || node.Message != msg || node.Version != version || node.Diff != diff {
			site.Nodes[i].Status = status
			site.Nodes[i].Message = msg
			site.Nodes[i].Version = version
			site.Nodes[i].Diff = diff
			UpdateSite(site.GetId(), site)
		}
	}
}
