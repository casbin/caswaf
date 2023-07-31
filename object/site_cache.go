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

import "github.com/casbin/caswaf/casdoor"

var siteMap = map[string]*Site{}

func InitSiteMap() {
	refreshSiteMap()
}

func getCertMap() map[string]*casdoor.Cert {
	certs, err := casdoor.GetCerts("admin")
	if err != nil {
		panic(err)
	}

	res := map[string]*casdoor.Cert{}
	for _, cert := range certs {
		res[cert.Name] = cert
	}
	return res
}

func getApplicationMap() map[string]*casdoor.Application {
	certMap := getCertMap()

	applications, err := casdoor.GetApplications("admin")
	if err != nil {
		panic(err)
	}

	res := map[string]*casdoor.Application{}
	for _, application := range applications {
		if application.Cert != "" {
			if cert, ok := certMap[application.Cert]; ok {
				application.CertObj = cert
			}
		}

		res[application.Name] = application
	}
	return res
}

func refreshSiteMap() {
	siteMap = map[string]*Site{}

	applicationMap := getApplicationMap()
	sites := GetGlobalSites()
	for _, site := range sites {
		if _, ok := siteMap[site.Domain]; !ok {
			if site.SslCert != "" {
				site.SslCertObj = getCert("admin", site.SslCert)
			}

			if site.CasdoorApplication != "" && site.ApplicationObj == nil {
				if v, ok2 := applicationMap[site.CasdoorApplication]; ok2 {
					site.ApplicationObj = v
				}
			}

			if site.Domain != "" && site.PublicIp == "" {
				go func(site *Site) {
					site.PublicIp = resolveDomainToIp(site.Domain)
					UpdateSiteNoRefresh(site.GetId(), site)
				}(site)
			}

			siteMap[site.Domain] = site
		}
	}
}

func GetSiteByDomain(domain string) *Site {
	if site, ok := siteMap[domain]; ok {
		return site
	} else {
		return nil
	}
}
