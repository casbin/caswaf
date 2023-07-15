package object

import (
	"github.com/casbin/caswaf/util"
	"xorm.io/core"
)

type Cert struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`

	Type            string `xorm:"varchar(100)" json:"type"`
	CryptoAlgorithm string `xorm:"varchar(100)" json:"cryptoAlgorithm"`
	ExpireTime      string `xorm:"varchar(100)" json:"expireTime"`

	Certificate string `xorm:"mediumtext" json:"certificate"`
	PrivateKey  string `xorm:"mediumtext" json:"privateKey"`
}

func GetGlobalCerts() []*Cert {
	certs := []*Cert{}
	err := adapter.engine.Asc("owner").Desc("created_time").Find(&certs)
	if err != nil {
		panic(err)
	}

	return certs
}

func GetCerts(owner string) []*Cert {
	certs := []*Cert{}
	err := adapter.engine.Desc("created_time").Find(&certs, &Cert{Owner: owner})
	if err != nil {
		panic(err)
	}

	return certs
}

func getCert(owner string, name string) *Cert {
	cert := Cert{Owner: owner, Name: name}
	existed, err := adapter.engine.Get(&cert)
	if err != nil {
		panic(err)
	}

	if existed {
		return &cert
	} else {
		return nil
	}
}

func GetCert(id string) *Cert {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getCert(owner, name)
}

func UpdateCert(id string, cert *Cert) bool {
	owner, name := util.GetOwnerAndNameFromId(id)
	if getCert(owner, name) == nil {
		return false
	}

	if cert.Certificate != "" {
		cert.ExpireTime = getCertExpireTime(cert.Certificate)
	} else {
		cert.ExpireTime = ""
	}

	_, err := adapter.engine.ID(core.PK{owner, name}).AllCols().Update(cert)
	if err != nil {
		panic(err)
	}

	//return affected != 0
	return true
}

func AddCert(cert *Cert) bool {
	affected, err := adapter.engine.Insert(cert)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteCert(cert *Cert) bool {
	affected, err := adapter.engine.ID(core.PK{cert.Owner, cert.Name}).Delete(&Cert{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}
