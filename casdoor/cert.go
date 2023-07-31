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

package casdoor

import "github.com/casbin/caswaf/util"

type Cert struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	DisplayName     string `xorm:"varchar(100)" json:"displayName"`
	Scope           string `xorm:"varchar(100)" json:"scope"`
	Type            string `xorm:"varchar(100)" json:"type"`
	CryptoAlgorithm string `xorm:"varchar(100)" json:"cryptoAlgorithm"`
	BitSize         int    `json:"bitSize"`
	ExpireInYears   int    `json:"expireInYears"`

	Certificate            string `xorm:"mediumtext" json:"certificate"`
	PrivateKey             string `xorm:"mediumtext" json:"privateKey"`
	AuthorityPublicKey     string `xorm:"mediumtext" json:"authorityPublicKey"`
	AuthorityRootPublicKey string `xorm:"mediumtext" json:"authorityRootPublicKey"`
}

func GetCerts(owner string) ([]*Cert, error) {
	certs := []*Cert{}
	err := adapter.Engine.Desc("created_time").Find(&certs, &Cert{Owner: owner})
	if err != nil {
		return certs, err
	}

	return certs, nil
}

func GetCert(id string) (*Cert, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getCert(owner, name)
}

func getCert(owner string, name string) (*Cert, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	cert := Cert{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&cert)
	if err != nil {
		return &cert, err
	}

	if existed {
		return &cert, nil
	} else {
		return nil, nil
	}
}
