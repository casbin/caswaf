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

package service

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/beego/beego"
	"github.com/casbin/caswaf/object"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

func joinPath(a string, b string) string {
	if strings.HasSuffix(a, "/") && strings.HasPrefix(b, "/") {
		b = b[1:]
	} else if !strings.HasSuffix(a, "/") && !strings.HasPrefix(b, "/") {
		b = "/" + b
	}
	res := a + b
	return res
}

func responseError(w http.ResponseWriter, format string, a ...interface{}) {
	w.WriteHeader(http.StatusInternalServerError)

	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)
	_, err := fmt.Fprintf(w, msg)
	if err != nil {
		panic(err)
	}
}

func getCertificateForDomain(domain string) (*tls.Certificate, error) {
	site := object.GetSiteByDomain(domain)
	if site == nil {
		return nil, fmt.Errorf("getCertificateForDomain() error: site not found for domain: [%s]", domain)
	}
	if site.SslCertObj == nil {
		return nil, fmt.Errorf("getCertificateForDomain() error: cert: [%s] not found for site: [%s]", site.SslCert, site.Name)
	}

	tlsCert, certErr := tls.X509KeyPair([]byte(site.SslCertObj.Certificate), []byte(site.SslCertObj.PrivateKey))

	return &tlsCert, certErr
}

func getCasdoorClientFromSite(site *object.Site) (*casdoorsdk.Client, error) {
	if site.ApplicationObj == nil {
		return nil, fmt.Errorf("site.ApplicationObj is empty")
	}

	casdoorEndpoint := beego.AppConfig.String("casdoorEndpoint")
	if casdoorEndpoint == "http://localhost:8000" {
		casdoorEndpoint = "http://localhost:7001"
	}

	clientId := site.ApplicationObj.ClientId
	clientSecret := site.ApplicationObj.ClientSecret

	certificate := ""
	if site.ApplicationObj.CertObj != nil {
		certificate = site.ApplicationObj.CertObj.Certificate
	}

	res := casdoorsdk.NewClient(casdoorEndpoint, clientId, clientSecret, certificate, site.ApplicationObj.Organization, site.CasdoorApplication)
	return res, nil
}
