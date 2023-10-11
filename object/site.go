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
	"strings"

	"github.com/casbin/caswaf/run"
	"github.com/casbin/caswaf/util"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"xorm.io/core"
)

type Node struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Diff     string `json:"diff"`
	Pid      int    `json:"pid"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	Provider string `json:"provider"`
}

type Site struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	UpdatedTime string `xorm:"varchar(100)" json:"updatedTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`

	Tag          string   `xorm:"varchar(100)" json:"tag"`
	Domain       string   `xorm:"varchar(100)" json:"domain"`
	OtherDomains []string `xorm:"varchar(500)" json:"otherDomains"`
	NeedRedirect bool     `json:"needRedirect"`
	Challenges   []string `xorm:"mediumtext" json:"challenges"`
	Host         string   `xorm:"varchar(100)" json:"host"`
	Port         int      `json:"port"`
	SslMode      string   `xorm:"varchar(100)" json:"sslMode"`
	SslCert      string   `xorm:"varchar(100)" json:"sslCert"`
	PublicIp     string   `xorm:"varchar(100)" json:"publicIp"`
	Node         string   `xorm:"varchar(100)" json:"node"`
	IsSelf       bool     `json:"isSelf"`
	Nodes        []*Node  `xorm:"mediumtext" json:"nodes"`

	CasdoorApplication string `xorm:"varchar(100)" json:"casdoorApplication"`

	SslCertObj     *Cert                   `xorm:"-" json:"sslCertObj"`
	ApplicationObj *casdoorsdk.Application `xorm:"-" json:"applicationObj"`
}

func GetGlobalSites() ([]*Site, error) {
	sites := []*Site{}
	err := ormer.Engine.Asc("owner").Desc("created_time").Find(&sites)
	if err != nil {
		return nil, err
	}

	return sites, nil
}

func GetSites(owner string) ([]*Site, error) {
	sites := []*Site{}
	err := ormer.Engine.Asc("tag").Asc("name").Desc("created_time").Find(&sites, &Site{Owner: owner})
	if err != nil {
		return nil, err
	}

	for _, site := range sites {
		if site.SslCert == "" && site.Domain != "" && (site.SslMode == "HTTPS and HTTP" || site.SslMode == "HTTPS Only") {
			site.SslCert, err = getBaseDomain(site.Domain)
			if err != nil {
				return nil, err
			}

			_, err = UpdateSite(site.GetId(), site)
			if err != nil {
				return nil, err
			}
		}
	}

	return sites, nil
}

func getSite(owner string, name string) (*Site, error) {
	site := Site{Owner: owner, Name: name}
	existed, err := ormer.Engine.Get(&site)
	if err != nil {
		return nil, err
	}

	if existed {
		return &site, nil
	}
	return nil, nil
}

func GetSite(id string) (*Site, error) {
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

func UpdateSite(id string, site *Site) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if s, err := getSite(owner, name); err != nil {
		return false, err
	} else if s == nil {
		return false, nil
	}

	site.UpdatedTime = util.GetCurrentTime()

	_, err := ormer.Engine.ID(core.PK{owner, name}).AllCols().Update(site)
	if err != nil {
		return false, err
	}

	err = refreshSiteMap()
	if err != nil {
		return false, err
	}

	err = site.checkNodes()
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateSiteNoRefresh(id string, site *Site) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if s, err := getSite(owner, name); err != nil {
		return false, err
	} else if s == nil {
		return false, nil
	}

	_, err := ormer.Engine.ID(core.PK{owner, name}).AllCols().Update(site)
	if err != nil {
		return false, err
	}

	return true, nil
}

func AddSite(site *Site) (bool, error) {
	affected, err := ormer.Engine.Insert(site)
	if err != nil {
		return false, err
	}

	if affected != 0 {
		err = refreshSiteMap()
		if err != nil {
			return false, err
		}
	}

	return affected != 0, nil
}

func DeleteSite(site *Site) (bool, error) {
	affected, err := ormer.Engine.ID(core.PK{site.Owner, site.Name}).Delete(&Site{})
	if err != nil {
		return false, err
	}

	if affected != 0 {
		err = refreshSiteMap()
		if err != nil {
			return false, err
		}
	}

	return affected != 0, nil
}

func (site *Site) GetId() string {
	return fmt.Sprintf("%s/%s", site.Owner, site.Name)
}

func (site *Site) GetChallengeMap() map[string]string {
	m := map[string]string{}
	for _, challenge := range site.Challenges {
		tokens := strings.Split(challenge, ":")
		m[tokens[0]] = tokens[1]
	}
	return m
}

func (site *Site) GetHost() string {
	if site.Host != "" {
		return site.Host
	}

	if site.Port == 0 {
		return ""
	}

	res := fmt.Sprintf("http://localhost:%d", site.Port)
	return res
}

func addErrorToMsg(msg string, function string, err error) string {
	if msg == "" {
		return fmt.Sprintf("%s(): %s", function, err.Error())
	} else {
		return fmt.Sprintf("%s || %s(): %s", msg, function, err.Error())
	}
}

func (site *Site) checkNodes() error {
	hostname := util.GetHostname()
	for i, node := range site.Nodes {
		if node.Name != hostname {
			continue
		}

		if site.GetHost() == "" {
			continue
		}

		ok, msg := pingUrl(site.GetHost())
		status := "Running"
		if !ok {
			msg = ""
			if node.Pid > 0 {
				var err error
				ok, err = run.IsProcessActive(node.Pid)
				if err != nil {
					msg = addErrorToMsg(msg, "IsProcessActive", err)
				}
			}
		}
		if !ok {
			status = "Stopped"
		}

		diff := ""
		if i != 0 {
			diff = site.Nodes[0].Diff
		}

		pid, err := run.CreateRepo(site.Name, !ok, diff, node.Provider)
		if err != nil {
			msg = addErrorToMsg(msg, "CreateRepo", err)
		}

		if pid == 0 {
			pid = node.Pid
		}

		version, err := getSiteVersion(site.Name)
		if err != nil {
			msg = addErrorToMsg(msg, "getSiteVersion", err)
		}

		path := run.GetRepoPath(site.Name)
		newDiff, err := run.GitDiff(path)
		if err != nil {
			msg = addErrorToMsg(msg, "GitDiff", err)
		}

		if node.Status != status || node.Message != msg || node.Version != version || node.Diff != newDiff || node.Pid != pid {
			site.Nodes[i].Version = version
			site.Nodes[i].Diff = newDiff
			site.Nodes[i].Pid = pid
			site.Nodes[i].Status = status
			site.Nodes[i].Message = msg
			_, err = UpdateSite(site.GetId(), site)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
