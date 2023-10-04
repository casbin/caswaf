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
	"testing"

	"github.com/casbin/caswaf/proxy"
	"github.com/casbin/caswaf/util"
)

func TestGetCertExpireTime(t *testing.T) {
	InitConfig()

	cert := getCert("admin", "casbin.com")
	println(getCertExpireTime(cert.Certificate))
}

func TestRenewAllCerts(t *testing.T) {
	InitConfig()
	proxy.InitHttpClient()

	certs := GetCerts("admin")
	for i, cert := range certs {
		res := RenewCert(cert)
		fmt.Printf("[%d/%d] Renewed cert: [%s] to [%s], res = %v\n", i+1, len(certs), cert.Name, cert.ExpireTime, res)
	}
}

func TestApplyAllCerts(t *testing.T) {
	InitConfig()

	baseDir := "F:/github_repos/nginx/conf/ssl"
	certs := GetCerts("admin")
	for _, cert := range certs {
		if cert.Certificate == "" || cert.PrivateKey == "" {
			continue
		}

		util.WriteStringToPath(cert.Certificate, fmt.Sprintf("%s/%s.pem", baseDir, cert.Name))
		util.WriteStringToPath(cert.PrivateKey, fmt.Sprintf("%s/%s.key", baseDir, cert.Name))
	}
}
