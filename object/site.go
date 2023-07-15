package object

import (
	"github.com/casbin/caswaf/util"
	"xorm.io/core"
)

type Site struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`

	Domain  string `xorm:"varchar(100)" json:"domain"`
	Host    string `xorm:"varchar(100)" json:"host"`
	SslMode string `xorm:"varchar(100)" json:"sslMode"`
	SslCert string `xorm:"varchar(100)" json:"sslCert"`
}

func GetGlobalSites() []*Site {
	sites := []*Site{}
	err := adapter.engine.Asc("owner").Desc("created_time").Find(&sites)
	if err != nil {
		panic(err)
	}

	return sites
}

func GetSites(owner string) []*Site {
	sites := []*Site{}
	err := adapter.engine.Desc("created_time").Find(&sites, &Site{Owner: owner})
	if err != nil {
		panic(err)
	}

	return sites
}

func getSite(owner string, name string) *Site {
	site := Site{Owner: owner, Name: name}
	existed, err := adapter.engine.Get(&site)
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

func UpdateSite(id string, site *Site) bool {
	owner, name := util.GetOwnerAndNameFromId(id)
	if getSite(owner, name) == nil {
		return false
	}

	_, err := adapter.engine.ID(core.PK{owner, name}).AllCols().Update(site)
	if err != nil {
		panic(err)
	}

	//return affected != 0
	return true
}

func AddSite(site *Site) bool {
	affected, err := adapter.engine.Insert(site)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteSite(site *Site) bool {
	affected, err := adapter.engine.ID(core.PK{site.Owner, site.Name}).Delete(&Site{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}
