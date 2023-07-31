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

	for _, cert := range certs {
		if cert.Certificate != "" && cert.ExpireTime == "" {
			cert.ExpireTime = getCertExpireTime(cert.Certificate)
			UpdateCert(cert.GetId(), cert)
		}
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

func (cert *Cert) GetId() string {
	return fmt.Sprintf("%s/%s", cert.Owner, cert.Name)
}
